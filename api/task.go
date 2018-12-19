package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/ankeesler/anwork/lag"
	"github.com/ankeesler/anwork/task"
	"github.com/tedsuo/rata"
)

func findTask(l *lag.L, repo task.Repo, w http.ResponseWriter, r *http.Request) (*task.Task, int) {
	id := rata.Param(r, "id")
	idN, err := strconv.Atoi(id)
	if err != nil {
		respondWithError(l, w, http.StatusBadRequest, err)
		return nil, 0
	}

	task, err := repo.FindTaskByID(idN)
	if err != nil {
		respondWithError(l, w, http.StatusInternalServerError, err)
		return nil, 0
	}

	if task == nil {
		l.P(lag.I, "unknown task with id %d", idN)
		respondWithError(l, w, http.StatusNotFound, fmt.Errorf("unknown task with ID %d", idN))
		return nil, 0
	}

	l.P(lag.I, "found task with id %d: %+v", idN, task)
	return task, idN
}

type getTaskHandler struct {
	l    *lag.L
	repo task.Repo
}

func (h *getTaskHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if task, _ := findTask(h.l, h.repo, w, r); task != nil {
		respond(h.l, w, http.StatusOK, task)
	}
}

type updateTaskHandler struct {
	l    *lag.L
	repo task.Repo
}

func (h *updateTaskHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respond(h.l, w, http.StatusInternalServerError, err)
		return
	}

	var newTask task.Task
	if err := json.Unmarshal(data, &newTask); err != nil {
		respond(h.l, w, http.StatusBadRequest, err)
		return
	}

	if task, idN := findTask(h.l, h.repo, w, r); task != nil {
		newTask.ID = idN
		if err := h.repo.UpdateTask(&newTask); err != nil {
			respondWithError(h.l, w, http.StatusInternalServerError, err)
			return
		}

		respond(h.l, w, http.StatusNoContent, nil)
	}
}

type deleteTaskHandler struct {
	l    *lag.L
	repo task.Repo
}

func (h *deleteTaskHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if task, _ := findTask(h.l, h.repo, w, r); task != nil {
		if err := h.repo.DeleteTask(task); err != nil {
			respondWithError(h.l, w, http.StatusInternalServerError, err)
			return
		}

		respond(h.l, w, http.StatusNoContent, task)
	}
}
