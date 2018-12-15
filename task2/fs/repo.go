// Package fs contains a task.Repo implementation which uses the filesystem to store
// task.Task's.
package fs

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/ankeesler/anwork/task2"
)

type repo struct {
	MyTasks    []*task2.Task `json:"tasks"`
	NextTaskID int

	MyEvents    []*task2.Event `json:"events"`
	NextEventID int

	file   string
	loaded bool
}

// New returns a task.Repo that stores task.Task's on the local filesystem.
//
// This task.Repo is NOT thread-safe.
func New(file string) task2.Repo {
	return &repo{file: file}
}

func (r *repo) CreateTask(task *task2.Task) error {
	if err := r.ensureLoaded(); err != nil {
		return err
	}

	task.ID = r.NextTaskID
	r.NextTaskID++

	r.MyTasks = append(r.MyTasks, task)

	return r.commit()
}

func (r *repo) Tasks() ([]*task2.Task, error) {
	if err := r.ensureLoaded(); err != nil {
		return nil, err
	}

	return r.MyTasks, nil
}

func (r *repo) FindTaskByID(id int) (*task2.Task, error) {
	if err := r.ensureLoaded(); err != nil {
		return nil, err
	}

	for _, task := range r.MyTasks {
		if task.ID == id {
			return task, nil
		}
	}

	return nil, nil
}

func (r *repo) FindTaskByName(name string) (*task2.Task, error) {
	if err := r.ensureLoaded(); err != nil {
		return nil, err
	}

	for _, task := range r.MyTasks {
		if task.Name == name {
			return task, nil
		}
	}

	return nil, nil
}

func (r *repo) UpdateTask(task *task2.Task) error {
	if err := r.ensureLoaded(); err != nil {
		return err
	}

	t, err := r.FindTaskByID(task.ID)
	if err != nil {
		panic(err)
	} else if t == nil {
		return &unknownTaskError{name: task.Name, id: task.ID}
	}

	*t = *task

	return r.commit()
}

func (r *repo) DeleteTask(task *task2.Task) error {
	if err := r.ensureLoaded(); err != nil {
		return err
	}

	index := -1
	for i, t := range r.MyTasks {
		if t.ID == task.ID {
			index = i
			break
		}
	}

	if index != -1 {
		r.MyTasks = append(r.MyTasks[:index], r.MyTasks[index+1:]...)
		return r.commit()
	} else {
		return &unknownTaskError{name: task.Name, id: task.ID}
	}
}

func (r *repo) CreateEvent(event *task2.Event) error {
	if err := r.ensureLoaded(); err != nil {
		return err
	}

	event.ID = r.NextEventID
	r.NextEventID++

	r.MyEvents = append(r.MyEvents, event)

	return r.commit()
}

func (r *repo) Events() ([]*task2.Event, error) {
	if err := r.ensureLoaded(); err != nil {
		return nil, err
	}

	return r.MyEvents, nil
}

func (r *repo) FindEventByID(id int) (*task2.Event, error) {
	index := r.findEvent(id)
	if index == -1 {
		return nil, nil
	} else {
		return r.MyEvents[index], nil
	}
}

func (r *repo) DeleteEvent(event *task2.Event) error {
	if err := r.ensureLoaded(); err != nil {
		return err
	}

	index := r.findEvent(event.ID)

	if index != -1 {
		r.MyEvents = append(r.MyEvents[:index], r.MyEvents[index+1:]...)
		return r.commit()
	} else {
		return &unknownEventError{title: event.Title, date: event.Date}
	}
}

func (r *repo) findEvent(id int) int {
	index := -1
	for i, e := range r.MyEvents {
		if e.ID == id {
			index = i
			break
		}
	}
	return index
}

func (r *repo) commit() error {
	data, err := json.Marshal(r)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(r.file, data, 0600)
}

func (r *repo) ensureLoaded() error {
	if r.loaded {
		return nil
	}

	if _, err := os.Stat(r.file); err == nil {
		data, err := ioutil.ReadFile(r.file)
		if err != nil {
			return err
		}

		if err := json.Unmarshal(data, r); err != nil {
			return err
		}
	}

	r.loaded = true

	return nil
}
