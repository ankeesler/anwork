package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/ankeesler/anwork/lag"
	taskpkg "github.com/ankeesler/anwork/task"
)

type getTasksHandler struct {
	l    *lag.L
	repo taskpkg.Repo
}

func (h *getTasksHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name != "" {
		task, err := h.repo.FindTaskByName(name)
		if err != nil {
			respondWithError(h.l, w, http.StatusInternalServerError, err)
			return
		}

		tasks := make([]*taskpkg.Task, 0, 1)
		if task != nil {
			tasks = append(tasks, task)
		}
		respond(h.l, w, http.StatusOK, tasks)
		return
	}

	tasks, err := h.repo.Tasks()
	if err != nil {
		respondWithError(h.l, w, http.StatusInternalServerError, err)
		return
	}

	respond(h.l, w, http.StatusOK, tasks)
}

type createTaskHandler struct {
	l    *lag.L
	repo taskpkg.Repo
}

func (h *createTaskHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respondWithError(h.l, w, http.StatusInternalServerError, err)
		return
	}

	var task taskpkg.Task
	if err := json.Unmarshal(data, &task); err != nil {
		respondWithError(h.l, w, http.StatusBadRequest, err)
		return
	}

	h.l.P(lag.I, "creating task %+v", task)
	if err := h.repo.CreateTask(&task); err != nil {
		respondWithError(h.l, w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Add("Location", fmt.Sprintf("/api/v1/tasks/%d", task.ID))
	respond(h.l, w, http.StatusCreated, nil)
}
