package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/ankeesler/anwork/task2"
)

type getTasksHandler struct {
	log  *log.Logger
	repo task2.Repo
}

func (h *getTasksHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name != "" {
		task, err := h.repo.FindTaskByName(name)
		if err != nil {
			respondWithError(h.log, w, http.StatusInternalServerError, err)
			return
		}

		tasks := make([]*task2.Task, 0, 1)
		if task != nil {
			tasks = append(tasks, task)
		}
		respond(h.log, w, http.StatusOK, tasks)
		return
	}

	tasks, err := h.repo.Tasks()
	if err != nil {
		respondWithError(h.log, w, http.StatusInternalServerError, err)
		return
	}

	respond(h.log, w, http.StatusOK, tasks)
}

type createTaskHandler struct {
	log  *log.Logger
	repo task2.Repo
}

func (h *createTaskHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respondWithError(h.log, w, http.StatusInternalServerError, err)
		return
	}

	var task task2.Task
	if err := json.Unmarshal(data, &task); err != nil {
		respondWithError(h.log, w, http.StatusBadRequest, err)
		return
	}

	h.log.Printf("creating task %+v", task)
	if err := h.repo.CreateTask(&task); err != nil {
		respondWithError(h.log, w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Add("Location", fmt.Sprintf("/api/v1/tasks/%d", task.ID))
	respond(h.log, w, http.StatusCreated, nil)
}
