package application

import (
	"github.com/google/uuid"
	"github.com/zalhui/calc_golang/internal/orchestrator/models"
	"github.com/zalhui/calc_golang/pkg/calculation"
)

type Application struct {
	repository *Repository
}

func New() *Application {
	return &Application{
		repository: NewRepository(),
	}
}

func (a *Application) AddExpression(expression string) (string, error) {
	expressionID := uuid.New().String()

	tasks, err := calculation.ParseExpression(expression, expressionID)
	if err != nil {
		return "", err
	}

	a.repository.AddExpression(&models.Expression{
		ID:     expressionID,
		Status: "pending",
		Tasks:  tasks,
	})

	return expressionID, nil
}

func (a *Application) GetExpressionByID(expressionID string) (*models.Expression, bool) {
	return a.repository.GetExpressionByID(expressionID)
}

func (a *Application) GetAllExpressions() []*models.Expression {
	return a.repository.GetAllExpressions()
}

func (a *Application) GetPendingTask() (*models.Task, bool) {
	return a.repository.GetPendingTask()
}

func (a *Application) UpdateTaskStatus(taskID string, status string, result float64) {
	a.repository.UpdateTaskStatus(taskID, status, result)

	r := a.repository
	r.mu.Lock()
	defer r.mu.Unlock()

	if task, found := r.tasks[taskID]; found {
		exprID := task.ExpressionID
		if _, exists := r.expressions[exprID]; exists {
			a.repository.UpdateExpressionStatus(exprID)
		}
	}
}
