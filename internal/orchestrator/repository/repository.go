package repository

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/zalhui/calc_golang/internal/common/models"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) AddExpression(expr *models.Expression) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	_, err = tx.Exec(
		"INSERT INTO expressions (id, user_id, expression, status, created_at) VALUES (?, ?, ?, ?, ?)",
		expr.ID, expr.UserID, expr.Expression, expr.Status, time.Now(),
	)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to insert expression: %w", err)
	}

	// Вставляем задачи
	for _, task := range expr.Tasks {
		deps := strings.Join(task.Dependencies, ",")
		_, err = tx.Exec(
			"INSERT INTO tasks (id, expression_id, arg1, arg2, operation, status, dependencies, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
			task.ID, expr.ID, task.Arg1, task.Arg2, task.Operation, task.Status, deps, time.Now(),
		)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to insert task: %w", err)
		}
	}

	return tx.Commit()
}

func (r *Repository) GetExpressionByID(expressionID, userID string) (*models.Expression, bool) {
	row := r.db.QueryRow(
		`SELECT id, user_id, expression, 
		status, result, created_at FROM 
		expressions WHERE id = ? AND user_id = ?`,
		expressionID, userID,
	)

	var expr models.Expression
	var createdAt time.Time
	err := row.Scan(
		&expr.ID,
		&expr.UserID,
		&expr.Expression,
		&expr.Status,
		&expr.Result,
		&createdAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, false
		}
		log.Printf("Error getting expression: %v", err)
		return nil, false
	}
	expr.CreatedAt = createdAt

	// Получаем связанные задачи
	tasks, err := r.getTasksForExpression(expr.ID)
	if err != nil {
		log.Printf("Error getting tasks: %v", err)
		return nil, false
	}
	expr.Tasks = tasks

	return &expr, true
}

func (r *Repository) GetAllExpressions(userID string) []*models.Expression {
	rows, err := r.db.Query(
		`SELECT id, expression, status, result, 
		created_at FROM expressions WHERE user_id = ?
		 ORDER BY created_at DESC`,
		userID,
	)
	if err != nil {
		log.Printf("Error querying expressions: %v", err)
		return nil
	}
	defer rows.Close()

	var expressions []*models.Expression
	for rows.Next() {
		var expr models.Expression
		var createdAt time.Time
		err := rows.Scan(
			&expr.ID,
			&expr.Expression,
			&expr.Status,
			&expr.Result,
			&createdAt,
		)
		if err != nil {
			log.Printf("Error scanning expression: %v", err)
			continue
		}
		expr.UserID = userID
		expr.CreatedAt = createdAt
		expressions = append(expressions, &expr)
	}

	return expressions
}

func (r *Repository) GetTaskByID(taskID string) (*models.Task, bool) {
	row := r.db.QueryRow(
		`SELECT id, expression_id, arg1, arg2, 
		operation, status, result, dependencies 
		FROM tasks WHERE id = ?`,
		taskID,
	)

	var task models.Task
	var deps string
	err := row.Scan(
		&task.ID,
		&task.ExpressionID,
		&task.Arg1,
		&task.Arg2,
		&task.Operation,
		&task.Status,
		&task.Result,
		&deps,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, false
		}
		log.Printf("Error getting task: %v", err)
		return nil, false
	}

	task.Dependencies = strings.Split(deps, ",")
	return &task, true
}

func (r *Repository) GetPendingTask() (*models.Task, bool) {
	rows, err := r.db.Query(
		`SELECT id, expression_id, arg1, arg2, 
		operation, dependencies FROM tasks 
		WHERE status = 'pending'`,
	)
	if err != nil {
		log.Printf("Error querying pending tasks: %v", err)
		return nil, false
	}
	defer rows.Close()

	for rows.Next() {
		var task models.Task
		var deps string
		err := rows.Scan(
			&task.ID,
			&task.ExpressionID,
			&task.Arg1,
			&task.Arg2,
			&task.Operation,
			&deps,
		)
		if err != nil {
			log.Printf("Error scanning task: %v", err)
			continue
		}
		task.Dependencies = strings.Split(deps, ",")
		task.Status = "pending"

		//if r.allDependenciesCompleted(task.Dependencies) {
		return &task, true
		//}
	}

	return nil, false
}

func (r *Repository) allDependenciesCompleted(dependencies []string) bool {
	for _, depID := range dependencies {
		var status string
		err := r.db.QueryRow(
			"SELECT status FROM tasks WHERE id = ?",
			depID,
		).Scan(&status)

		if err != nil || status != "completed" {
			return false
		}
	}
	return true
}

func (r *Repository) UpdateTaskStatus(taskID string, status string, result float64) {
	tx, err := r.db.Begin()
	if err != nil {
		log.Printf("Error starting transaction: %v", err)
		return
	}

	// Обновляем статус задачи
	_, err = tx.Exec(
		"UPDATE tasks SET status = ?, result = ? WHERE id = ?",
		status, result, taskID,
	)
	if err != nil {
		tx.Rollback()
		log.Printf("Error updating task status: %v", err)
		return
	}

	// Получаем expression_id для обновления статуса выражения
	var expressionID string
	err = tx.QueryRow(
		"SELECT expression_id FROM tasks WHERE id = ?",
		taskID,
	).Scan(&expressionID)
	if err != nil {
		tx.Rollback()
		log.Printf("Error getting expression ID: %v", err)
		return
	}

	// Проверяем все ли задачи выражения выполнены
	var pendingTasks int
	err = tx.QueryRow(
		`SELECT COUNT(*) FROM tasks WHERE expression_id = ? 
		AND status NOT IN ('completed', 'error')`,
		expressionID,
	).Scan(&pendingTasks)

	if err != nil {
		tx.Rollback()
		log.Printf("Error checking pending tasks: %v", err)
		return
	}

	if pendingTasks == 0 {
		exprStatus := "completed"
		var finalResult float64

		// Проверяем наличие ошибок в задачах
		var errorTasks int
		err = tx.QueryRow(
			`SELECT COUNT(*) FROM tasks WHERE 
        expression_id = ? AND status = 'error'`,
			expressionID,
		).Scan(&errorTasks)
		if err != nil {
			tx.Rollback()
			log.Printf("Error checking error tasks: %v", err)
			return
		}

		if errorTasks > 0 {
			exprStatus = "error"
			finalResult = 0 // или другое значение по умолчанию
		} else {
			finalResult = result
		}

		_, err = tx.Exec(
			"UPDATE expressions SET status = ?, result = ? WHERE id = ?",
			exprStatus, finalResult, expressionID,
		)
		if err != nil {
			tx.Rollback()
			log.Printf("Error updating expression status: %v", err)
			return
		}
	}

	if err = tx.Commit(); err != nil {
		log.Printf("Error committing transaction: %v", err)
	}
}

func (r *Repository) getTasksForExpression(expressionID string) ([]*models.Task, error) {
	rows, err := r.db.Query(
		`SELECT id, arg1, arg2, operation, status, 
		result, dependencies FROM tasks WHERE expression_id = ?`,
		expressionID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*models.Task
	for rows.Next() {
		var task models.Task
		var deps string
		err := rows.Scan(
			&task.ID,
			&task.Arg1,
			&task.Arg2,
			&task.Operation,
			&task.Status,
			&task.Result,
			&deps,
		)
		if err != nil {
			return nil, err
		}
		task.ExpressionID = expressionID
		task.Dependencies = strings.Split(deps, ",")
		tasks = append(tasks, &task)
	}

	return tasks, nil
}

func (r *Repository) GetUserHistory(userID string) ([]*models.Expression, error) {
	rows, err := r.db.Query(
		`SELECT id, expression, status, result, created_at 
		FROM expressions WHERE user_id = ? ORDER BY created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []*models.Expression
	for rows.Next() {
		var expr models.Expression
		var createdAt time.Time
		err := rows.Scan(
			&expr.ID,
			&expr.Expression,
			&expr.Status,
			&expr.Result,
			&createdAt,
		)
		if err != nil {
			return nil, err
		}
		expr.UserID = userID
		expr.CreatedAt = createdAt
		history = append(history, &expr)
	}

	return history, nil
}
