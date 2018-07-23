package remote

import "github.com/ankeesler/anwork/task"

type managerFactory struct {
	client APIClient
}

func NewManagerFactory(client APIClient) task.ManagerFactory {
	return &managerFactory{client: client}
}

func (mf *managerFactory) Create() (task.Manager, error) {
	return newManager(mf.client), nil
}

func (mf *managerFactory) Save(manager task.Manager) error {
	return nil
}
