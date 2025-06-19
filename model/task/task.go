package task

import "time"

const (
	StatusAwaiting    = 0
	StatusExecuting   = 1
	StatusCompleted   = 2
	StatusInterrupted = 3
)

var statuses = [4]string{"awaiting", "executing", "completed", "interrupted"}

type Task struct {
	Name       string    `json:"task_name"`
	CreatedAt  time.Time `json:"created_at"`
	StatusText string    `json:"status"`

	Interrupt chan struct{} `json:"-"`
	status    int           `json:"-"`
}

func New(name string) Task {
	t := Task{
		Name:      name,
		CreatedAt: time.Now(),
		Interrupt: make(chan struct{}),
	}
	t.SetStatus(StatusAwaiting)

	return t
}

func (t *Task) Status() int {
	return t.status
}

func (t *Task) SetStatus(status int) {
	t.status = status
	t.StatusText = statuses[status]
}
