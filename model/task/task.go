package task

import "time"

const (
	StatusAwaiting  = 0
	StatusExecuting = 1
	StatusCompleted = 2
)

var statuses = [3]string{"awaiting", "executing", "completed"}

type Task struct {
	Name          string    `json:"task_name"`
	CreatedAt     time.Time `json:"created_at"`
	StatusText    string    `json:"status"`
	executionTime time.Duration

	Interrupt chan struct{}
	status    int
}

func New(name string) *Task {
	t := &Task{
		Name:      name,
		CreatedAt: time.Now(),
		Interrupt: make(chan struct{}),
	}
	t.SetStatus(StatusAwaiting)

	return t
}

// ExecutionTime возвращает время работы задачи. Если задача выполняется, то время работы пересчитывается как time.Since(task.CreatedAt).
func (t *Task) ExecutionTime() time.Duration {
	// если задача не завершена, то пересчитываем время работы
	if t.status == StatusExecuting {
		t.executionTime = time.Since(t.CreatedAt)
	}
	return t.executionTime
}

func (t *Task) Dto() Dto {
	return Dto{
		Name:          t.Name,
		CreatedAt:     t.CreatedAt,
		StatusText:    t.StatusText,
		ExecutionTime: t.ExecutionTime(),
	}
}

func (t *Task) Status() int {
	return t.status
}

func (t *Task) SetStatus(status int) {
	t.status = status
	t.StatusText = statuses[status]
}
