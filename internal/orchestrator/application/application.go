package application

import (
	"log"

	"github.com/google/uuid"
	"github.com/zalhui/calc_golang/internal/orchestrator/models"
	"github.com/zalhui/calc_golang/pkg/calculation"
)

type Application struct {
	repository *Repository
}

type ExpressionDTO struct {
	ID     string  `json:"id"`
	Status string  `json:"status"`
	Result float64 `json:"result,omitempty"`
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
		Status: "pending", // Устанавливаем начальный статус
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
	r.mu.RLock()
	defer r.mu.RUnlock()

	for exprID, expr := range r.expressions {
		for _, task := range expr.Tasks {
			if task.ID == taskID {
				log.Printf("Task %s updated to %s with result %f, triggering update for expression %s", taskID, status, result, exprID)
				a.repository.UpdateExpressionStatus(exprID)
				return // Выходим после первого совпадения
			}
		}
		log.Printf("Task %s not found in tasks of expression %s", taskID, exprID)
	}
	log.Printf("No expression found containing task %s", taskID)
}
