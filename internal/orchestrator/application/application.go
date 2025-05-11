package application

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/zalhui/calc_golang/internal/auth"
	"github.com/zalhui/calc_golang/internal/common/models"
	"github.com/zalhui/calc_golang/internal/orchestrator/repository"
	"github.com/zalhui/calc_golang/pkg/calculation"
)

type Application struct {
	repository *repository.Repository
	db         *sql.DB
}

func New(db *sql.DB) *Application {
	return &Application{
		repository: repository.NewRepository(db),
		db:         db,
	}
}

// RegisterUser регистрирует нового пользователя
func (a *Application) RegisterUser(login, password string) (*models.UserResponse, error) {
	// Проверяем существует ли пользователь
	var count int
	err := a.db.QueryRow("SELECT COUNT(*) FROM users WHERE login = ?", login).Scan(&count)
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, fmt.Errorf("user with login %s already exists", login)
	}

	// Создаем нового пользователя
	userID := uuid.New().String()
	_, err = auth.RegisterUser(a.db, login, password)
	if err != nil {
		return nil, err
	}

	// Возвращаем данные пользователя
	return &models.UserResponse{
		ID:        userID,
		Login:     login,
		CreatedAt: time.Now(),
	}, nil
}

// AuthenticateUser выполняет аутентификацию пользователя
func (a *Application) AuthenticateUser(login, password string) (*models.AuthResponse, error) {
	token, err := auth.LoginUser(a.db, login, password)
	if err != nil {
		return nil, fmt.Errorf("failed to login user: %w", err)
		//models.ErrInvalidCredentials
	}

	return &models.AuthResponse{
		Token: token,
	}, nil
}

// AddExpression добавляет новое выражение для вычисления
func (a *Application) AddExpression(expression string, userID string) (string, error) {
	expressionID := uuid.New().String()

	log.Printf("Creating expression %s for user %s", expressionID, userID)

	tasks, err := calculation.ParseExpression(expression, expressionID)
	if err != nil {
		return "", err
	}

	expr := &models.Expression{
		ID:         expressionID,
		UserID:     userID,
		Expression: expression,
		Status:     "pending",
		Tasks:      tasks,
		CreatedAt:  time.Now(),
	}

	if err := a.repository.AddExpression(expr); err != nil {
		log.Printf("Failed to save expression: %v", err)
		return "", fmt.Errorf("failed to save expression")
	}

	log.Printf("Successfully created expression %s", expressionID)
	return expressionID, nil
}

// GetExpressionByID возвращает выражение по ID
func (a *Application) GetExpressionByID(expressionID, userID string) (*models.ExpressionResponse, error) {
	expr, exists := a.repository.GetExpressionByID(expressionID, userID)
	if !exists {
		return nil, fmt.Errorf("expression with ID %s not found", expressionID)
	}

	var result *float64
	if expr.Status == "completed" {
		result = &expr.Result.Float64
	}
	return &models.ExpressionResponse{
		ID:         expr.ID,
		Expression: expr.Expression,
		Status:     expr.Status,
		Result:     result,
		CreatedAt:  expr.CreatedAt,
		FinishedAt: expr.FinishedAt,
	}, nil
}

// GetAllExpressions возвращает все выражения пользователя
func (a *Application) GetAllExpressions(userID string) ([]*models.ExpressionResponse, error) {
	expressions := a.repository.GetAllExpressions(userID)
	response := make([]*models.ExpressionResponse, 0, len(expressions))

	for _, expr := range expressions {
		var result *float64
		if expr.Status == "completed" {
			result = &expr.Result.Float64
		}

		response = append(response, &models.ExpressionResponse{
			ID:         expr.ID,
			Expression: expr.Expression,
			Status:     expr.Status,
			Result:     result,
			CreatedAt:  expr.CreatedAt,
			FinishedAt: expr.FinishedAt,
		})
	}

	return response, nil
}

// GetPendingTask возвращает следующую задачу для вычисления
func (a *Application) GetPendingTask() (*models.TaskResponse, error) {
	task, exists := a.repository.GetPendingTask()
	if !exists {
		return nil, fmt.Errorf("no pending tasks")
	}

	return &models.TaskResponse{
		ID:           task.ID,
		ExpressionID: task.ExpressionID,
		Operation:    task.Operation,
		Status:       task.Status,
		CreatedAt:    task.CreatedAt,
	}, nil
}

// UpdateTaskResult обновляет результат выполнения задачи
func (a *Application) UpdateTaskResult(taskID string, result float64, errMsg string) error {
	if errMsg != "" {
		a.repository.UpdateTaskStatus(taskID, "error", 0)
		return fmt.Errorf("task %s failed: %s", taskID, errMsg)
	}

	a.repository.UpdateTaskStatus(taskID, "completed", result)
	return nil
}

// GetUserHistory возвращает историю вычислений пользователя
func (a *Application) GetUserHistory(userID string) ([]*models.ExpressionResponse, error) {
	history, err := a.repository.GetUserHistory(userID)
	if err != nil {
		return nil, err
	}

	response := make([]*models.ExpressionResponse, 0, len(history))
	for _, expr := range history {
		var result *float64
		if expr.Status == "completed" {
			result = &expr.Result.Float64
		}

		response = append(response, &models.ExpressionResponse{
			ID:         expr.ID,
			Expression: expr.Expression,
			Status:     expr.Status,
			Result:     result,
			CreatedAt:  expr.CreatedAt,
			FinishedAt: expr.FinishedAt,
		})
	}

	return response, nil
}

// GetTaskResult возвращает результат выполнения задачи
func (a *Application) GetTaskResult(taskID string) (*float64, error) {
	task, exists := a.repository.GetTaskByID(taskID)
	if !exists {
		return nil, fmt.Errorf("task with ID %s not found", taskID)
	}

	if task.Status != "completed" {
		return nil, fmt.Errorf("task with ID %s is not completed", taskID)
	}

	return &task.Result.Float64, nil
}
