package taskscontroller

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/aachex/taskrunner/model/task"
)

type TasksController struct {
	tasks  map[string]*task.Task
	logger *slog.Logger
}

func New(logger *slog.Logger) *TasksController {
	return &TasksController{
		tasks:  make(map[string]*task.Task),
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

	type response struct {
		Message string `json:"message"`
	}
	writeJson(w, response{fmt.Sprintf("created task %s", taskName)}, http.StatusCreated)
}

func (tc *TasksController) TaskStatus(w http.ResponseWriter, r *http.Request) {
	taskName := r.PathValue("name")
	t, ok := tc.tasks[taskName]
	if !ok {
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}

	writeJson(w, t.Dto(), http.StatusOK)
}

func (tc *TasksController) DeleteTask(w http.ResponseWriter, r *http.Request) {
	taskName := r.PathValue("name")

	t, ok := tc.tasks[taskName]
	if !ok {
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}

	// если задача выполняется в данный момент - останавливаем её
	if t.Status() == task.StatusExecuting {
		t.Interrupt <- struct{}{}
	}
	delete(tc.tasks, taskName)

	type response struct {
		Message string `json:"message"`
	}
	writeJson(w, response{fmt.Sprintf("task %s deleted", taskName)}, http.StatusOK)
}

func (tc *TasksController) runTask(name string) error {
	// нельзя запускать уже существующую задачу
	if _, ok := tc.tasks[name]; ok {
		return fmt.Errorf("task %s already exists", name)
	}

	go func() {
		tc.tasks[name] = task.New(name)
		t := tc.tasks[name]

		tc.logger.Debug("running new task", "name", t.Name, "createdAt", t.CreatedAt, "status", t.Status())
		t.Run()

		switch t.Status() {
		case task.StatusInterrupted:
			tc.logger.Info("task interrupted", "name", name)
			delete(tc.tasks, name)
			return

		case task.StatusCompleted:
			tc.logger.Info("task completed", "name", name, "execution_time", t.ExecutionTime())
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
