package repository

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/zalhui/calc_golang/internal/common/models"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(`
		CREATE TABLE expressions (
			id TEXT PRIMARY KEY,
			user_id TEXT,
			expression TEXT,
			status TEXT,
			result REAL
		);
		CREATE TABLE tasks (
			id TEXT PRIMARY KEY,
			expression_id TEXT,
			arg1 TEXT,
			arg2 TEXT,
			operation TEXT,
			status TEXT,
			result REAL,
			dependencies TEXT
		);
	`)
	if err != nil {
		t.Fatal(err)
	}

	return db
}

func TestAddAndGetExpression(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepository(db)

	expr := &models.Expression{
		ID:         "test-id",
		UserID:     "user1",
		Expression: "2+2",
		Status:     "pending",
		Tasks: []*models.Task{
			{
				ID:        "task1",
				Arg1:      "2",
				Arg2:      "2",
				Operation: "+",
				Status:    "pending",
			},
		},
	}

	t.Run("Add expression", func(t *testing.T) {
		err := repo.AddExpression(expr)
		if err != nil {
			t.Errorf("AddExpression failed: %v", err)
		}
	})

	t.Run("Get expression", func(t *testing.T) {
		found, exists := repo.GetExpressionByID("test-id", "user1")
		if !exists {
			t.Error("Expression not found")
		}
		if found.Expression != "2+2" {
			t.Errorf("Unexpected expression: %s", found.Expression)
		}
	})
}
