package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	api "github.com/ankeesler/anwork/api2"
	"github.com/ankeesler/anwork/task2"
)

type client struct {
	address string
}

// New returns a new API client pointed at an ANWORK API address.
func New(address string) task2.Repo {
	return &client{address: address}
}

func (c *client) CreateTask(task *task2.Task) error {
	return c.do(http.MethodPost, c.tasksURL(), task, nil)
}

func (c *client) Tasks() ([]*task2.Task, error) {
	tasks := make([]*task2.Task, 10)
	if err := c.do(http.MethodGet, c.tasksURL(), nil, &tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (c *client) FindTaskByID(id int) (*task2.Task, error) {
	var task task2.Task

	if err := c.do(http.MethodGet, c.taskURL(id), nil, &task); err != nil {
		return nil, err
	}

	return &task, nil
}

func (c *client) FindTaskByName(name string) (*task2.Task, error) {
	var task task2.Task

	url := fmt.Sprintf("%s?name=%s", c.tasksURL(), name)
	if err := c.do(http.MethodGet, url, nil, &task); err != nil {
		return nil, err
	}

	return &task, nil
}

func (c *client) UpdateTask(task *task2.Task) error {
	return c.do(http.MethodPut, c.taskURL(task.ID), task, nil)
}

func (c *client) DeleteTask(task *task2.Task) error {
	return c.do(http.MethodDelete, c.taskURL(task.ID), nil, nil)
}

func (c *client) CreateEvent(event *task2.Event) error {
	return c.do(http.MethodPost, c.eventsURL(), event, nil)
}

func (c *client) Events() ([]*task2.Event, error) {
	events := make([]*task2.Event, 10)
	if err := c.do(http.MethodGet, c.eventsURL(), nil, &events); err != nil {
		return nil, err
	}

	return events, nil
}

func (c *client) FindEventByID(id int) (*task2.Event, error) {
	var event task2.Event

	if err := c.do(http.MethodGet, c.eventURL(id), nil, &event); err != nil {
		return nil, err
	}

	return &event, nil
}

func (c *client) DeleteEvent(event *task2.Event) error {
	return c.do(http.MethodDelete, c.eventURL(event.ID), nil, nil)
}

func (c *client) tasksURL() string {
	return fmt.Sprintf("http://%s/api/v1/tasks", c.address)
}

func (c *client) taskURL(id int) string {
	return fmt.Sprintf("http://%s/api/v1/tasks/%d", c.address, id)
}

func (c *client) eventsURL() string {
	return fmt.Sprintf("http://%s/api/v1/events", c.address)
}

func (c *client) eventURL(id int) string {
	return fmt.Sprintf("http://%s/api/v1/events/%d", c.address, id)
}

func (c *client) do(method, url string, input interface{}, output interface{}) error {
	body, err := c.encodeBody(input)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return err
	}

	if input != nil {
		req.Header.Add("Content-Type", "application/json")
	}
	if output != nil {
		req.Header.Add("Accept", "application/json")
	}

	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()

	if c.is5xxStatus(rsp) {
		return &badResponseError{code: rsp.Status, message: c.decodeError(rsp.Body)}
	} else if c.is4xxStatus(rsp) {
		return &badResponseError{code: rsp.Status}
	}

	return c.decodeBody(rsp.Body, output)
}

func (c *client) encodeBody(input interface{}) (io.Reader, error) {
	var body io.Reader
	if input != nil {
		data, err := json.Marshal(input)
		if err != nil {
			return nil, err
		}

		body = bytes.NewBuffer(data)
	}
	return body, nil
}

func (c *client) decodeBody(body io.Reader, output interface{}) error {
	if output != nil {
		bytes, err := ioutil.ReadAll(body)
		if err != nil {
			return err
		}

		if err := json.Unmarshal(bytes, output); err != nil {
			return fmt.Errorf("cannot unmarshal response body (%s): '%s'", err.Error(), string(bytes))
		}
	}

	return nil
}

func (c *client) decodeError(body io.Reader) string {
	bodyData, err := ioutil.ReadAll(body)
	if err != nil {
		return fmt.Sprintf("??? (ReadAll: %s)", err.Error())
	}

	errMsg := api.Error{}
	if err := json.Unmarshal(bodyData, &errMsg); err != nil {
		return fmt.Sprintf("??? (Unmarshal: %s)", err.Error())
	}

	return errMsg.Message
}

func (c *client) is4xxStatus(rsp *http.Response) bool {
	return rsp.StatusCode >= 400 && rsp.StatusCode < 500
}

func (c *client) is5xxStatus(rsp *http.Response) bool {
	return rsp.StatusCode >= 500 && rsp.StatusCode < 600
}
