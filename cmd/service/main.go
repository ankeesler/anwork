package main

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/ankeesler/anwork/cmd/anwork/command"
	task "github.com/ankeesler/anwork/tasknew"
)

type handler interface {
	http.Handler
	Manager() (*task.Manager, error)
}

type tasksHandler struct {
	m *task.Manager
}

func newTasksHandler(m *task.Manager) *tasksHandler {
	return &tasksHandler{m: m}
}

func (h *tasksHandler) Manager() (*task.Manager, error) {
	if h.m == nil {
		return nil, errors.New("Manager is nil!")
	}
	return h.m, nil
}

func (h *tasksHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cmd := command.FindCommand("show")
	if cmd == nil {
		http.Error(w, "Cannot find command 'show'", http.StatusInternalServerError)
	} else {
		cmd.Run(flag.NewFlagSet("", flag.ContinueOnError), w, h.m)
	}
}

func main() {
	m := task.NewManager()
	h := newTasksHandler(m)
	http.Handle("/api/v1/tasks", h)

	port := os.Getenv("PORT")
	if port == "" {
		port = "54321"
	}
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil); err != nil {
		fmt.Println("ERROR!", err)
		os.Exit(1)
	}
}
