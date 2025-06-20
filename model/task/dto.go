package task

import "time"

type Dto struct {
	Name          string        `json:"task_name"`
	CreatedAt     time.Time     `json:"created_at"`
	Status        string        `json:"status"`
	ExecutionTime time.Duration `json:"execution_time"`
}
