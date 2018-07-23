// This package contains a client for the ANWORK API.
package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/ankeesler/anwork/api"
	"github.com/ankeesler/anwork/task"
)

// This is an HTTP client for the ANWORK API.
type Client struct {
	address    string
	httpClient http.Client
}

type badResponseError struct {
	status  string
	payload []byte
}

func (bre *badResponseError) Error() string {
	if bre.status != "" {
		return fmt.Sprintf("Unexpected response status: %s", bre.status)
	} else {
		return fmt.Sprintf("Unexpected response payload: %s", string(bre.payload))
	}
}

type unknownTaskError struct {
	id int
}

func (ute *unknownTaskError) Error() string {
	return fmt.Sprintf("Unknown task with ID %d", ute.id)
}

// Create a new Client attached to an ANWORK API at some address.
func New(address string) *Client {
	return &Client{address: address}
}

// Create a task.Task.
func (c *Client) CreateTask(name string) error {
	url := fmt.Sprintf("%s/api/v1/tasks", c.address)
	payload, err := json.Marshal(api.CreateRequest{Name: name})
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(payload))
	if err != nil {
		panic(err)
	}

	req.Header["content-type"] = []string{"application/json"}
	rsp, err := c.httpClient.Do(req)
	if err != nil {
		return err
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

// Delete a task.Task.
func (c *Client) DeleteTask(id int) error {
	url := fmt.Sprintf("%s/api/v1/tasks/%d", c.address, id)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}

	rsp, err := c.httpClient.Do(req)
	if err != nil {
		return err
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

// Get all of the tasks stored on the remote.
func (c *Client) GetTasks() ([]*task.Task, error) {
	url := fmt.Sprintf("%s/api/v1/tasks", c.address)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header["Accept"] = []string{"application/json"}
	rsp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	var tasks []*task.Task
	if rsp.StatusCode == http.StatusOK {
		payload, err := ioutil.ReadAll(rsp.Body)
		if err != nil || len(payload) == 0 {
			return nil, fmt.Errorf("unexpected response: %s", string(payload))
		}
		if err := json.Unmarshal(payload, &tasks); err != nil {
			return nil, err
		}
	} else {
		return nil, &badResponseError{status: rsp.Status}
	}

	return tasks, nil
}

// Get the task.Task associated with an ID.
func (c *Client) GetTask(id int) (*task.Task, error) {
	url := fmt.Sprintf("%s/api/v1/tasks/%d", c.address, id)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header["Accept"] = []string{"application/json"}
	rsp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	var task *task.Task
	if rsp.StatusCode == http.StatusOK {
		if err := unmarshal(rsp.Body, &task); err != nil {
			return nil, err
		}
	} else if rsp.StatusCode == http.StatusNotFound {
		return nil, &unknownTaskError{id: id}
	} else {
		return nil, &badResponseError{status: rsp.Status}
	}

	return task, nil
}

// Update a task.Task's priority.
func (c *Client) UpdatePriority(id int, prio int) error {
	return c.updateTask(id, api.UpdateTaskRequest{Priority: prio})
}

// Update a task.Task's task.State.
func (c *Client) UpdateState(id int, state task.State) error {
	return c.updateTask(id, api.UpdateTaskRequest{State: state})
}

// Delete an event.
func (c *Client) DeleteEvent(startTime int) error {
	url := fmt.Sprintf("%s/api/v1/events/%d", c.address, startTime)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}

	rsp, err := c.httpClient.Do(req)
	if err != nil {
		return err
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

// Get all of the events.
func (c *Client) GetEvents() ([]*task.Event, error) {
	url := fmt.Sprintf("%s/api/v1/events", c.address)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header["Accept"] = []string{"application/json"}
	rsp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	var events []*task.Event
	if rsp.StatusCode == http.StatusOK {
		if err := unmarshal(rsp.Body, &events); err != nil {
			return nil, err
		}
	} else {
		return nil, &badResponseError{status: rsp.Status}
	}

	return events, nil
}

func (c *Client) updateTask(id int, update api.UpdateTaskRequest) error {
	url := fmt.Sprintf("%s/api/v1/tasks/%d", c.address, id)
	payload, err := json.Marshal(update)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	req.Header["content-type"] = []string{"application/json"}
	rsp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()

	if rsp.StatusCode == http.StatusNotFound {
		return &unknownTaskError{id: id}
	} else if rsp.StatusCode != http.StatusNoContent {
		otaErr, err := readErrorResponse(rsp)
		if err != nil {
			return err
		} else {
			return otaErr
		}
	}

	return nil
}

func (c *Client) deleteEvent(event *task.Event) error {
	url := fmt.Sprintf("%s/api/v1/events/%d", c.address, event.TaskID)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		panic(err)
	}

	rsp, err := c.httpClient.Do(req)
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
	payload, err := ioutil.ReadAll(rsp.Body)
	if err != nil || len(payload) == 0 {
		return nil, &badResponseError{payload: payload}
	}

	if err := json.Unmarshal(payload, &otaErr); err != nil {
		return nil, err
	}
	return &otaErr, nil
}

func unmarshal(body io.Reader, thing interface{}) error {
	payload, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	} else if payload == nil || len(payload) == 0 {
		return &badResponseError{payload: payload}
	}

	return json.Unmarshal(payload, thing)

}
