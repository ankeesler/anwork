// This package contains a task.Manager that acts as an HTTP client for the anwork
// API. All of the state that this type manipulates lives on a remote system.
package remote

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ankeesler/anwork/api"
	"github.com/ankeesler/anwork/task"
)

type manager struct {
	address    string
	httpClient http.Client
}

func newManager(address string) *manager {
	return &manager{address: address}
}

func (m *manager) Create(name string) error {
	url := fmt.Sprintf("%s/api/v1/tasks", m.address)
	payload, err := json.Marshal(api.CreateRequest{Name: name})
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(payload))
	if err != nil {
		panic(err)
	}

	req.Header["content-type"] = []string{"application/json"}
	rsp, err := m.httpClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer rsp.Body.Close()

	if rsp.StatusCode != http.StatusCreated {
		otaErr, err := readErrorResponse(rsp)
		if err != nil {
			return err
		} else {
			return otaErr
		}
	}

	return nil
}

func (m *manager) Delete(name string) error {
	tasks, err := m.getTasks()
	if err != nil {
		return err
	}

	var toDelete *task.Task
	for _, task := range tasks {
		if task.Name == name {
			toDelete = task
			break
		}
	}

	if toDelete == nil {
		return fmt.Errorf("unknown task with name %s", name)
	}

	url := fmt.Sprintf("%s/api/v1/tasks/%d", m.address, toDelete.ID)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		panic(err)
	}

	rsp, err := m.httpClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer rsp.Body.Close()

	if rsp.StatusCode != http.StatusNoContent {
		otaErr, err := readErrorResponse(rsp)
		if err != nil {
			return err
		} else {
			return otaErr
		}
	}

	return nil
}

func (m *manager) Tasks() []*task.Task {
	tasks, err := m.getTasks()
	if err != nil {
		panic(err)
	}

	return tasks
}

func (m *manager) FindByName(name string) *task.Task {
	tasks, err := m.getTasks()
	if err != nil {
		panic(err)
	}

	var foundIt *task.Task
	for _, task := range tasks {
		if task.Name == name {
			foundIt = task
			break
		}
	}

	return foundIt
}

func (m *manager) FindByID(id int) *task.Task {
	url := fmt.Sprintf("%s/api/v1/tasks/%d", m.address, id)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		panic(err)
	}

	req.Header["Accept"] = []string{"application/json"}
	rsp, err := m.httpClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer rsp.Body.Close()

	var task *task.Task
	if rsp.StatusCode == http.StatusOK {
		decoder := json.NewDecoder(rsp.Body)
		if err := decoder.Decode(&task); err != nil {
			panic(err)
		}
	}

	return task
}

func (m *manager) SetPriority(name string, prio int) error {
	return m.updateTask(name, api.UpdateTaskRequest{Priority: prio})
}

func (m *manager) SetState(name string, state task.State) error {
	return m.updateTask(name, api.UpdateTaskRequest{State: state})
}

func (m *manager) Note(name, note string) error {
	t := m.FindByName(name)
	if t == nil {
		return fmt.Errorf("unknown task %s", name)
	}

	url := fmt.Sprintf("%s/api/v1/events", m.address)
	payload, err := json.Marshal(api.AddEventRequest{
		Title:  fmt.Sprintf("Note added to task %s: %s", name, note),
		Date:   time.Now().Unix(),
		Type:   task.EventTypeNote,
		TaskID: t.ID,
	})
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(payload))
	if err != nil {
		panic(err)
	}

	req.Header["Content-Type"] = []string{"application/json"}
	rsp, err := m.httpClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer rsp.Body.Close()

	return nil
}

func (m *manager) Events() []*task.Event {
	url := fmt.Sprintf("%s/api/v1/events", m.address)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		panic(err)
	}

	req.Header["Accept"] = []string{"application/json"}
	rsp, err := m.httpClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer rsp.Body.Close()

	var events []*task.Event
	if rsp.StatusCode == http.StatusOK {
		decoder := json.NewDecoder(rsp.Body)
		if err := decoder.Decode(&events); err != nil {
			panic(err)
		}
	} else {
		panic(fmt.Sprintf("received not OK status code: %s", rsp.Status))
	}

	return events
}

func (m *manager) getTasks() ([]*task.Task, error) {
	url := fmt.Sprintf("%s/api/v1/tasks", m.address)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header["Accept"] = []string{"application/json"}
	rsp, err := m.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	var tasks []*task.Task
	if rsp.StatusCode == http.StatusOK {
		decoder := json.NewDecoder(rsp.Body)
		if err := decoder.Decode(&tasks); err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("GET /api/v1/tasks returned unknown status: %s", rsp.Status)
	}

	return tasks, nil
}

func (m *manager) updateTask(name string, update api.UpdateTaskRequest) error {
	task := m.FindByName(name)
	if task == nil {
		return fmt.Errorf("unknown task %s", name)
	}

	url := fmt.Sprintf("%s/api/v1/tasks/%d", m.address, task.ID)
	payload, err := json.Marshal(update)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(payload))
	if err != nil {
		panic(err)
	}

	req.Header["content-type"] = []string{"application/json"}
	rsp, err := m.httpClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer rsp.Body.Close()

	if rsp.StatusCode != http.StatusNoContent {
		otaErr, err := readErrorResponse(rsp)
		if err != nil {
			return err
		} else {
			return otaErr
		}
	}

	return nil
}

func readErrorResponse(rsp *http.Response) (*api.ErrorResponse, error) {
	var otaErr api.ErrorResponse
	decoder := json.NewDecoder(rsp.Body)
	if err := decoder.Decode(&otaErr); err != nil {
		return nil, err
	}
	return &otaErr, nil
}
