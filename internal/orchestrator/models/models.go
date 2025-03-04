package models

import "time"

type Expression struct {
	ID     string  `json:"id"`
	Status string  `json:"status"`
	Result float64 `json:"result,omitempty"`
	Tasks  []Task  `json:"tasks,omitempty"`
}

type Task struct {
	ID            string        `json:"id"`
	ExpressionID  string        `json:"expression_id"`
	Arg1          string        `json:"arg1"`
	Arg2          string        `json:"arg2"`
	Operation     string        `json:"operation"`
	OperationTime time.Duration `json:"operation_time"`
	Status        string        `json:"status"`
	Result        float64       `json:"result,omitempty"`
}
