package task

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	StatusAwaiting    = "awaiting"
	StatusExecuting   = "executing"
	StatusCompleted   = "completed"
	StatusInterrupted = "interrupted"
)

type Task struct {
	Name      string    `json:"task_name"`
	CreatedAt time.Time `json:"created_at"`

	Interrupt chan struct{}

	executionTime time.Duration
	status        string
}

func New(name string) *Task {
	t := &Task{
		Name:      name,
		CreatedAt: time.Now(),
		Interrupt: make(chan struct{}),
		status:    StatusAwaiting,
	}

	return t
}

func (t *Task) Run() {
	t.status = StatusExecuting

	defer close(t.Interrupt)

	select {
	case <-t.Interrupt:
		fmt.Printf("%s interrupted\n", t.Name)
		t.status = StatusInterrupted
		return

	case <-time.After(time.Minute * time.Duration(rand.Intn(3)+3)):
		fmt.Printf("%s completed\n", t.Name)
		t.ExecutionTime() // фиксируем итоговое время работы
		t.status = StatusCompleted
	}
}

// ExecutionTime возвращает время работы задачи. Если задача выполняется, то время работы пересчитывается как time.Since(task.CreatedAt).
func (t *Task) ExecutionTime() time.Duration {
	// если задача не завершена, то пересчитываем время работы
	if t.Status() == StatusExecuting {
		t.executionTime = time.Since(t.CreatedAt)
	}
	return t.executionTime
}

func (t *Task) Dto() Dto {
	return Dto{
		Name:          t.Name,
		CreatedAt:     t.CreatedAt,
		Status:        t.Status(),
		ExecutionTime: t.ExecutionTime(),
	}
}

func (t *Task) Status() string {
	return t.status
}
