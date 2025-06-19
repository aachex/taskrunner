package taskscontroller

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

var statusText [3]string = [3]string{"none", "executing", "completed"}

type TasksController struct {
	tasks  map[string]task
	logger *slog.Logger
}

func New(logger *slog.Logger) *TasksController {
	return &TasksController{
		tasks:  make(map[string]task),
		logger: logger,
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
	task, ok := tc.tasks[taskName]
	if !ok {
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}

	result := struct {
		Name   string `json:"task_name"`
		Status string `json:"status"`
	}{}

	result.Name = taskName
	result.Status = statusText[task.Status]
	writeJson(w, result, http.StatusOK)
}

func (tc *TasksController) DeleteTask(w http.ResponseWriter, r *http.Request) {
	taskName := r.PathValue("name")

	task, ok := tc.tasks[taskName]
	if !ok {
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}

	if task.Status == statusExecuting {
		tc.tasks[taskName].Interrupt <- struct{}{}
		return
	}

	delete(tc.tasks, taskName)
}

func (tc *TasksController) runTask(name string) error {
	task := tc.tasks[name]
	// нельзя запускать уже запущенную задачу
	if task.Status == statusExecuting {
		return fmt.Errorf("task %s is already executing", name)
	}

	go func() {
		task.Status = statusExecuting
		task.Interrupt = make(chan struct{})
		defer close(task.Interrupt)

		tc.tasks[name] = task

		// либо задача успешно выполняется, либо её останавливают
		select {
		case <-task.Interrupt:
			delete(tc.tasks, name)
		case <-time.After(time.Second * 15):
			task.Status = statusCompleted
			tc.tasks[name] = task
		}
	}()

	return nil
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
