package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/ankeesler/anwork/api"
	"github.com/ankeesler/anwork/task"
)

//go:generate counterfeiter . Cache

// Cache is a very simple interface for a string cache.
type Cache interface {
	// Get returns the string from the cache, if there is one. Iff there is
	// no string in the cache, it will return "", false.
	Get() (string, bool)
	// Set sets the string in the cache.
	Set(string)
}

type client struct {
	tokenCache Cache

	address string
}

// New returns a new API client pointed at an ANWORK API address.
func New(address string, cache Cache) task.Repo {
	return &client{address: address, tokenCache: cache}
}

func (c *client) CreateTask(task *task.Task) error {
	rsp, err := c.doExt(http.MethodPost, c.tasksURL(), task, nil)
	if err != nil {
		return err
	}

	location := rsp.Header.Get("Location")
	if location == "" || !parseID(location, &task.ID) {
		return fmt.Errorf("could not parse ID from Location response header: %s", location)
	}

	return nil
}

func (c *client) Tasks() ([]*task.Task, error) {
	tasks := make([]*task.Task, 10)
	if err := c.do(http.MethodGet, c.tasksURL(), nil, &tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (c *client) FindTaskByID(id int) (*task.Task, error) {
	var task task.Task

	rsp, err := c.doExt(http.MethodGet, c.taskURL(id), nil, &task)
	if rsp != nil && rsp.StatusCode == http.StatusNotFound {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		return &task, nil
	}
}

func (c *client) FindTaskByName(name string) (*task.Task, error) {
	tasks := make([]*task.Task, 0, 1)

	url := fmt.Sprintf("%s?name=%s", c.tasksURL(), name)
	if err := c.do(http.MethodGet, url, nil, &tasks); err != nil {
		return nil, err
	}

	if len(tasks) == 0 {
		return nil, nil
	} else {
		return tasks[0], nil
	}
}

func (c *client) UpdateTask(task *task.Task) error {
	return c.do(http.MethodPut, c.taskURL(task.ID), task, nil)
}

func (c *client) DeleteTask(task *task.Task) error {
	return c.do(http.MethodDelete, c.taskURL(task.ID), nil, nil)
}

func (c *client) CreateEvent(event *task.Event) error {
	rsp, err := c.doExt(http.MethodPost, c.eventsURL(), event, nil)
	if err != nil {
		return err
	}

	location := rsp.Header.Get("Location")
	if location == "" || !parseID(location, &event.ID) {
		return fmt.Errorf("could not parse ID from Location response header: %s", location)
	}

	return nil
}

func (c *client) Events() ([]*task.Event, error) {
	events := make([]*task.Event, 10)
	if err := c.do(http.MethodGet, c.eventsURL(), nil, &events); err != nil {
		return nil, err
	}

	return events, nil
}

func (c *client) FindEventByID(id int) (*task.Event, error) {
	var event task.Event

	rsp, err := c.doExt(http.MethodGet, c.eventURL(id), nil, &event)
	if rsp != nil && rsp.StatusCode == http.StatusNotFound {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		return &event, nil
	}
}

func (c *client) DeleteEvent(event *task.Event) error {
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

func (c *client) authURL() string {
	return fmt.Sprintf("http://%s/api/v1/auth", c.address)
}

func (c *client) do(method, url string, input interface{}, output interface{}) error {
	_, err := c.doExt(method, url, input, output)
	return err
}

func (c *client) doExt(method, url string, input interface{}, output interface{}) (*http.Response, error) {
	body, err := encodeBody(input)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	if input != nil {
		req.Header.Add("Content-Type", "application/json")
	}
	if output != nil {
		req.Header.Add("Accept", "application/json")
	}

	token, err := c.getToken()
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", token)

	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	if is5xxStatus(rsp) {
		return rsp, &badResponseError{code: rsp.Status, message: decodeError(rsp.Body)}
	} else if is4xxStatus(rsp) {
		return rsp, &badResponseError{code: rsp.Status}
	}

	return rsp, decodeBody(rsp.Body, output)
}

func (c *client) getToken() (string, error) {
	if token, ok := c.tokenCache.Get(); ok {
		return token, nil
	}

	token, err := c.reallyGetToken()
	if err != nil {
		return "", err
	}
	c.tokenCache.Set(token)

	return token, nil
}

func (c *client) reallyGetToken() (string, error) {
	req, err := http.NewRequest(http.MethodPost, c.authURL(), nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("Accept", "application/json")

	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer rsp.Body.Close()

	if is5xxStatus(rsp) {
		return "", &badResponseError{code: rsp.Status, message: decodeError(rsp.Body)}
	} else if is4xxStatus(rsp) {
		return "", &badResponseError{code: rsp.Status}
	}

	var auth api.Auth
	if err := decodeBody(rsp.Body, &auth); err != nil {
		return "", err
	}

	decryptedToken := decryptToken(auth.Token)
	return fmt.Sprintf("bearer %s", decryptedToken), nil
}

func decryptToken(token string) string {
	// TODO: implement me!
	return token
}

func encodeBody(input interface{}) (io.Reader, error) {
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

func decodeBody(body io.Reader, output interface{}) error {
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

func decodeError(body io.Reader) string {
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

func is4xxStatus(rsp *http.Response) bool {
	return rsp.StatusCode >= 400 && rsp.StatusCode < 500
}

func is5xxStatus(rsp *http.Response) bool {
	return rsp.StatusCode >= 500 && rsp.StatusCode < 600
}

func parseID(location string, id *int) bool {
	segments := strings.Split(location, "/")
	idS := segments[len(segments)-1]
	idN, err := strconv.Atoi(idS)
	if err != nil {
		return false
	}

	*id = idN
	return true
}
