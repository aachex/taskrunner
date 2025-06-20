package task

import "time"

const (
	StatusAwaiting  = "awaiting"
	StatusExecuting = "executing"
	StatusCompleted = "completed"
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

func (t *Task) SetStatus(status string) {
	t.status = status
}
