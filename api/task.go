package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/ankeesler/anwork/task"
	"github.com/tedsuo/rata"
)

func findTask(log *log.Logger, repo task.Repo, w http.ResponseWriter, r *http.Request) (*task.Task, int) {
	id := rata.Param(r, "id")
	idN, err := strconv.Atoi(id)
	if err != nil {
		respondWithError(log, w, http.StatusBadRequest, err)
		return nil, 0
	}

	task, err := repo.FindTaskByID(idN)
	if err != nil {
		respondWithError(log, w, http.StatusInternalServerError, err)
		return nil, 0
	}

	if task == nil {
		log.Printf("unknown task with id %d", idN)
		respondWithError(log, w, http.StatusNotFound, fmt.Errorf("unknown task with ID %d", idN))
		return nil, 0
	}

	log.Printf("found task with id %d: %+v", idN, task)
	return task, idN
}

type getTaskHandler struct {
	log  *log.Logger
	repo task.Repo
}

func (h *getTaskHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if task, _ := findTask(h.log, h.repo, w, r); task != nil {
		respond(h.log, w, http.StatusOK, task)
	}
}

type updateTaskHandler struct {
	log  *log.Logger
	repo task.Repo
}

func (h *updateTaskHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respond(h.log, w, http.StatusInternalServerError, err)
		return
	}

	var newTask task.Task
	if err := json.Unmarshal(data, &newTask); err != nil {
		respond(h.log, w, http.StatusBadRequest, err)
		return
	}

	if task, idN := findTask(h.log, h.repo, w, r); task != nil {
		newTask.ID = idN
		if err := h.repo.UpdateTask(&newTask); err != nil {
			respondWithError(h.log, w, http.StatusInternalServerError, err)
			return
		}

		respond(h.log, w, http.StatusNoContent, nil)
	}
}

type deleteTaskHandler struct {
	log  *log.Logger
	repo task.Repo
}

func (h *deleteTaskHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if task, _ := findTask(h.log, h.repo, w, r); task != nil {
		if err := h.repo.DeleteTask(task); err != nil {
			respondWithError(h.log, w, http.StatusInternalServerError, err)
			return
		}

		respond(h.log, w, http.StatusNoContent, task)
	}
}
