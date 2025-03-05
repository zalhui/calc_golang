package application

import (
	"testing"

	"github.com/zalhui/calc_golang/internal/orchestrator/models"
)

func TestAddExpression(t *testing.T) {
	repo := NewRepository()
	expr := &models.Expression{
		ID:     "test-id",
		Status: "pending",
		Tasks:  []*models.Task{{ID: "task-1", ExpressionID: "test-id", Status: "pending"}},
	}
	repo.AddExpression(expr)

	if _, exists := repo.expressions["test-id"]; !exists {
		t.Errorf("Expression not added to repository")
	}
	if _, exists := repo.tasks["task-1"]; !exists {
		t.Errorf("Task not added to repository")
	}
}
