package taskscontroller

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

const (
	statusExecuting = 1
	statusCompleted = 2
)

var statusText [3]string = [3]string{"none", "executing", "completed"}

type TasksController struct {
	status    map[string]int
	interrupt map[string]chan struct{}

	logger *slog.Logger
}

func New(logger *slog.Logger) *TasksController {
	return &TasksController{
		status:    make(map[string]int),
		interrupt: make(map[string]chan struct{}),
		logger:    logger,
	}
}

func (tc *TasksController) RegisterEndpoints(mux *http.ServeMux) {
	mux.HandleFunc("/task/{name}/run", tc.RunTask)
	mux.HandleFunc("/task/{name}/status", tc.TaskStatus)
	mux.HandleFunc("/task/{name}/rm", tc.DeleteTask)
}

// HTTP POST /task/{name}/run
func (tc *TasksController) RunTask(w http.ResponseWriter, r *http.Request) {
	taskName := r.PathValue("name")

	if err := tc.runTask(taskName); err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	tc.logger.Info("successfully created task", "name", taskName)

	result := struct {
		Name string `json:"task_name"`
	}{taskName}
	writeJson(w, result, http.StatusCreated)
}

func (tc *TasksController) TaskStatus(w http.ResponseWriter, r *http.Request) {
	taskName := r.PathValue("name")
	status, ok := tc.status[taskName]
	if !ok {
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}

	result := struct {
		Name   string `json:"task_name"`
		Status string `json:"status"`
	}{}

	result.Name = taskName
	result.Status = statusText[status]
	writeJson(w, result, http.StatusOK)
}

func (tc *TasksController) DeleteTask(w http.ResponseWriter, r *http.Request) {
	taskName := r.PathValue("name")

	status, ok := tc.status[taskName]
	if !ok {
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}

	if status == statusExecuting {
		tc.interrupt[taskName] <- struct{}{}
		return
	}

	tc.deleteTask(taskName)
	fmt.Println(tc.status)
}

func (tc *TasksController) runTask(name string) error {
	// нельзя запускать уже запущенную задачу
	if tc.status[name] == statusExecuting {
		return fmt.Errorf("task %s is already executing", name)
	}

	go func() {
		tc.status[name] = statusExecuting
		tc.interrupt[name] = make(chan struct{})

		// либо задача успешно выполняется, либо её останавливают
		select {
		case <-tc.interrupt[name]:
			tc.deleteTask(name)
		case <-time.After(time.Second * 15):
			tc.status[name] = statusCompleted
		}
	}()

	return nil
}

func (tc *TasksController) deleteTask(name string) {
	close(tc.interrupt[name])
	delete(tc.interrupt, name)
	delete(tc.status, name)
}

func writeJson(w http.ResponseWriter, obj any, statusCode int) {
	b, err := json.Marshal(obj)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(b)
}
