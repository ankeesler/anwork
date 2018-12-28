package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"code.cloudfoundry.org/lager"
	"github.com/ankeesler/anwork/task"
	"github.com/tedsuo/rata"
)

func findTask(
	logger lager.Logger,
	repo task.Repo,
	w http.ResponseWriter,
	r *http.Request,
) (*task.Task, int) {
	id := rata.Param(r, "id")
	idN, err := strconv.Atoi(id)
	if err != nil {
		respondWithError(logger, w, http.StatusBadRequest, err)
		return nil, 0
	}

	task, err := repo.FindTaskByID(idN)
	if err != nil {
		respondWithError(logger, w, http.StatusInternalServerError, err)
		return nil, 0
	}

	if task == nil {
		logger.Debug("unknown-task", lager.Data{"id": idN})
		respondWithError(logger, w, http.StatusNotFound, fmt.Errorf("unknown task with ID %d", idN))
		return nil, 0
	}

	logger.Debug("found-task", lager.Data{"task": task})
	return task, idN
}

type getTaskHandler struct {
	logger lager.Logger
	repo   task.Repo
}

func (h *getTaskHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if task, _ := findTask(h.logger, h.repo, w, r); task != nil {
		respond(h.logger, w, http.StatusOK, task)
	}
}

type updateTaskHandler struct {
	logger lager.Logger
	repo   task.Repo
}

func (h *updateTaskHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respond(h.logger, w, http.StatusInternalServerError, err)
		return
	}

	var newTask task.Task
	if err := json.Unmarshal(data, &newTask); err != nil {
		respond(h.logger, w, http.StatusBadRequest, err)
		return
	}

	if task, idN := findTask(h.logger, h.repo, w, r); task != nil {
		newTask.ID = idN
		if err := h.repo.UpdateTask(&newTask); err != nil {
			respondWithError(h.logger, w, http.StatusInternalServerError, err)
			return
		}

		respond(h.logger, w, http.StatusNoContent, nil)
	}
}

type deleteTaskHandler struct {
	logger lager.Logger
	repo   task.Repo
}

func (h *deleteTaskHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if task, _ := findTask(h.logger, h.repo, w, r); task != nil {
		if err := h.repo.DeleteTask(task); err != nil {
			respondWithError(h.logger, w, http.StatusInternalServerError, err)
			return
		}

		respond(h.logger, w, http.StatusNoContent, task)
	}
}
