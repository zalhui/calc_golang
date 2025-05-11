package models

import (
	"database/sql"
	"time"
)

type Expression struct {
	ID         string          `json:"id"`
	UserID     string          `json:"user_id"`
	Expression string          `json:"expression"`
	Status     string          `json:"status"`
	Result     sql.NullFloat64 `json:"result,omitempty"`
	CreatedAt  time.Time       `json:"created_at"`
	FinishedAt time.Time       `json:"finished_at,omitempty"`
	Tasks      []*Task         `json:"tasks,omitempty"`
}

type Task struct {
	ID            string          `json:"id"`
	ExpressionID  string          `json:"expression_id"`
	Arg1          string          `json:"arg1"`
	Arg2          string          `json:"arg2"`
	Operation     string          `json:"operation"`
	OperationTime time.Duration   `json:"operation_time"`
	Status        string          `json:"status"`
	Result        sql.NullFloat64 `json:"result,omitempty"`
	Dependencies  []string        `json:"dependencies"`
	CreatedAt     time.Time       `json:"created_at"`
	StartedAt     time.Time       `json:"started_at,omitempty"`
	FinishedAt    time.Time       `json:"finished_at,omitempty"`
}
type ExpressionResponse struct {
	ID         string    `json:"id"`
	Expression string    `json:"expression"`
	Status     string    `json:"status"`
	Result     *float64  `json:"result,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	FinishedAt time.Time `json:"finished_at,omitempty"`
}
type TaskResponse struct {
	ID           string    `json:"id"`
	ExpressionID string    `json:"expression_id"`
	Operation    string    `json:"operation"`
	Status       string    `json:"status"`
	Result       *float64  `json:"result,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

type UserResponse struct {
	ID        string    `json:"id"`
	Login     string    `json:"login"`
	CreatedAt time.Time `json:"created_at"`
}
type AuthResponse struct {
	Token string `json:"token"`
}
