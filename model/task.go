package model

import "time"

type Task struct {
	Name       string    `json:"task_name"`
	CreatedAt  time.Time `json:"created_at"`
	StatusText string    `json:"status"`

	Interrupt chan struct{} `json:"-"`
	status    int           `json:"-"`
}

func (t *Task) Status() int {
	return t.status
}

func (t *Task) SetStatus(status int) {
	statuses := [2]string{"executing", "completed"}
	t.status = status
	t.StatusText = statuses[status-1]
}
