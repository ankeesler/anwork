package api

import (
	"net/http"

	"github.com/ankeesler/anwork/task2"
)

func tasksGetHandler(hd *handlerData, repo task2.Repo) (int, interface{}, error) {
	tasks, err := repo.Tasks()
	if err != nil {
		return http.StatusInternalServerError, nil, err
	} else {
		return http.StatusOK, tasks, nil
	}
}
