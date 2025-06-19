package taskscontroller

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"taskrunner/model"
	"time"
)

const (
	statusExecuting = 1
	statusCompleted = 2
)

type TasksController struct {
	tasks  map[string]model.Task
	logger *slog.Logger
}

func New(logger *slog.Logger) *TasksController {
	return &TasksController{
		tasks:  make(map[string]model.Task),
		logger: logger,
	}
}

func (tc *TasksController) RegisterEndpoints(mux *http.ServeMux) {
	mux.HandleFunc("POST /task/{name}/run", tc.RunTask)
	mux.HandleFunc("GET /task/{name}/status", tc.TaskStatus)
	mux.HandleFunc("DELETE /task/{name}/rm", tc.DeleteTask)
}

// HTTP POST /task/{name}/run
func (tc *TasksController) RunTask(w http.ResponseWriter, r *http.Request) {
	taskName := r.PathValue("name")

	if err := tc.runTask(taskName); err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	tc.logger.Info("created task", "name", taskName)

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

	writeJson(w, task, http.StatusOK)
}

func (tc *TasksController) DeleteTask(w http.ResponseWriter, r *http.Request) {
	taskName := r.PathValue("name")

	task, ok := tc.tasks[taskName]
	if !ok {
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}

	if task.Status() == statusExecuting {
		tc.tasks[taskName].Interrupt <- struct{}{}
		return
	}

	delete(tc.tasks, taskName)
}

func (tc *TasksController) runTask(name string) error {
	task := tc.tasks[name]

	// нельзя запускать уже запущенную задачу
	if task.Status() == statusExecuting {
		return fmt.Errorf("task %s is already executing", name)
	}

	go func() {
		// инициализация задачи
		task.Name = name
		task.SetStatus(statusExecuting)
		task.CreatedAt = time.Now()
		task.Interrupt = make(chan struct{})
		defer close(task.Interrupt)

		// т.к. нельзя напрямую изменять значение в мапе, приходится присваивать изменённую копию
		tc.tasks[name] = task

		select {
		case <-task.Interrupt:
			tc.logger.Info("task interrupted", "name", task.Name)
			delete(tc.tasks, name)

		case <-time.After(time.Second * 15):
			tc.logger.Info("task completed", "name", task.Name)
			task.SetStatus(statusCompleted)
			tc.tasks[name] = task // обновление статуса в мапе
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
