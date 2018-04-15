package local

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	task "github.com/ankeesler/anwork/tasknew"
)

type managerFactory struct {
	outputDir, context string
}

func NewManagerFactory(outputDir, context string) task.ManagerFactory {
	return &managerFactory{
		outputDir: outputDir,
		context:   context,
	}
}

func (mf *managerFactory) Create() (task.Manager, error) {
	if err := mf.validateOutputDir(); err != nil {
		return nil, err
	}

	contextFile := mf.contextFile()
	if _, err := os.Stat(contextFile); os.IsNotExist(err) {
		return &manager{}, nil
	}

	bytes, err := ioutil.ReadFile(contextFile)
	if err != nil {
		return nil, fmt.Errorf("could not read context file (%s): %s", contextFile, err.Error())
	}

	manager := &manager{}
	err = json.Unmarshal(bytes, manager)
	if err != nil {
		return nil, fmt.Errorf("could not read manager from context file (%s): %s",
			contextFile, err.Error())
	}
	return manager, nil
}

func (mf *managerFactory) Save(manager task.Manager) error {
	if err := mf.validateOutputDir(); err != nil {
		return err
	}

	bytes, err := json.Marshal(manager)
	if err != nil {
		return fmt.Errorf("could not marshal manager: %s", err.Error())
	}

	if err := ioutil.WriteFile(mf.contextFile(), bytes, 0644); err != nil {
		return fmt.Errorf("could not write manager to file (%s): %s", mf.contextFile(), err.Error())
	}

	return nil
}

func (mf *managerFactory) validateOutputDir() error {
	if _, err := os.Stat(mf.outputDir); os.IsNotExist(err) {
		return fmt.Errorf("outputDir does not exist: %s", mf.outputDir)
	}
	return nil
}

func (mf *managerFactory) contextFile() string {
	return filepath.Join(mf.outputDir, mf.context)
}
