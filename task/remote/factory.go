package remote

import "github.com/ankeesler/anwork/task"

type managerFactory struct {
	address string
}

func NewManagerFactory(address string) task.ManagerFactory {
	return &managerFactory{address: address}
}

func (mf *managerFactory) Create() (task.Manager, error) {
	return newManager(mf.address), nil
}

func (mf *managerFactory) Save(manager task.Manager) error {
	return nil
}

func (mf *managerFactory) Reset() error {
	return nil
}
