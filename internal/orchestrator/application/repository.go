package application

import (
	"log"
	"sync"

	"github.com/zalhui/calc_golang/internal/orchestrator/models"
)

type Repository struct {
	expressions map[string]*models.Expression
	tasks       map[string]*models.Task
	mu          sync.RWMutex
}

func NewRepository() *Repository {
	return &Repository{
		expressions: make(map[string]*models.Expression),
		tasks:       make(map[string]*models.Task),
	}
}

func (r *Repository) AddExpression(expression *models.Expression) {
	log.Printf("Добавление выражения: ID=%s, Задачи=%v\n", expression.ID, expression.Tasks)

	r.mu.Lock()
	r.expressions[expression.ID] = expression
	r.mu.Unlock()

	for _, task := range expression.Tasks {
		log.Printf("Добавление задачи: ID=%s, Arg1=%s, Arg2=%s, Operation=%s\n", task.ID, task.Arg1, task.Arg2, task.Operation)
		r.tasks[task.ID] = task
	}
}

func (r *Repository) GetExpressionByID(expressionID string) (*models.Expression, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	expression, exists := r.expressions[expressionID]

	return expression, exists
}
func (r *Repository) GetAllExpressions() []*models.Expression {
	r.mu.RLock()
	defer r.mu.RUnlock()

	expressions := make([]*models.Expression, 0, len(r.expressions))

	for _, expression := range r.expressions {
		expressions = append(expressions, expression)
	}

	return expressions
}

func (r *Repository) GetTaskByID(taskID string) (*models.Task, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	task, exists := r.tasks[taskID]

	return task, exists
}

func (r *Repository) GetPendingTask() (*models.Task, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, task := range r.tasks {
		if task.Status == "pending" && r.allDependenciesCompleted(task.Dependencies) {
			return task, true
		}
	}
	return nil, false
}

func (r *Repository) allDependenciesCompleted(dependencies []string) bool {
	for _, depID := range dependencies {
		depTask, exists := r.tasks[depID]
		if !exists || depTask.Status != "completed" {
			return false
		}
	}
	return true
}

func (r *Repository) UpdateTaskStatus(taskID string, status string, result float64) {
	r.mu.Lock()
	defer r.mu.Unlock()

	task, exists := r.tasks[taskID]
	if exists {
		task.Status = status
		if status == "completed" {
			task.Result = result
		} else if status == "error" {

			for _, expr := range r.expressions {
				for _, t := range expr.Tasks {
					if t.ID == taskID {
						expr.Status = "error"
						break
					}
				}
			}
		}
	}
}

func (r *Repository) UpdateExpressionStatus(expressionID string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	expression, exists := r.expressions[expressionID]
	if !exists {
		log.Printf("Expression %s not found in repository", expressionID)
		return
	}

	allCompleted := true

	for _, task := range expression.Tasks {
		storedTask, found := r.tasks[task.ID]
		if !found || storedTask.Status != "completed" {
			allCompleted = false
			break
		}
	}

	if allCompleted {

		lastTaskID := expression.Tasks[len(expression.Tasks)-1].ID
		if lastTask, found := r.tasks[lastTaskID]; found {
			expression.Status = "completed"
			expression.Result = lastTask.Result
			log.Printf("Expression %s updated to completed with result: %f", expressionID, lastTask.Result)
		} else {
			log.Printf("Last task %s not found for expression %s", lastTaskID, expressionID)
		}
	} else {
		expression.Status = "pending"
		log.Printf("Expression %s remains pending", expressionID)
	}

	r.expressions[expressionID] = expression
}
