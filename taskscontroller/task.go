package taskscontroller

import "time"

const (
	statusExecuting = 1
	statusCompleted = 2
)

type task struct {
	Name      string    `json:"task_name"`
	CreatedAt time.Time `json:"created_at"`
	Status    int
	Interrupt chan struct{}
}
