package application

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func (a *Application) AddExpressionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "only POST method allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Expression string `json:"expression"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	expressionID, err := a.AddExpression(req.Expression)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": expressionID})
}

func (a *Application) GetExpressionByIDHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "only GET method allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	expressionID := vars["id"]

	expression, exists := a.repository.GetExpressionByID(expressionID)
	if !exists {
		http.Error(w, "expression not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"expression": expression})
}

func (a *Application) GetAllExpressionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "only GET method allowed", http.StatusMethodNotAllowed)
		return
	}

	expressions := a.repository.GetAllExpressions()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"expressions": expressions})
}

func (a *Application) GetPendingTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "only GET method allowed", http.StatusMethodNotAllowed)
		return
	}

	task, exists := a.repository.GetPendingTask()
	if !exists {
		log.Printf("No pending tasks available")
		http.Error(w, "no pending task", http.StatusNotFound)
		return
	}

	log.Printf("Sending task to agent: ID=%s, Arg1=%s, Arg2=%s, Operation=%s, Status=%s",
		task.ID, task.Arg1, task.Arg2, task.Operation, task.Status)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]interface{}{"task": task}); err != nil {
		log.Printf("Error encoding task %s: %v", task.ID, err)
	}
}

func (a *Application) SubmitTaskResultHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// Обработка GET-запроса для получения результата задачи
		taskID := r.URL.Query().Get("id")
		if taskID == "" {
			http.Error(w, "Missing task id", http.StatusBadRequest)
			return
		}

		task, found := a.repository.GetTaskByID(taskID)
		if !found || task.Status != "completed" {
			http.Error(w, "Task result not ready", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"result": task.Result,
		})
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "only POST method allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ID     string  `json:"id"`
		Result float64 `json:"result,omitempty"`
		Error  string  `json:"error,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Error != "" {
		a.repository.UpdateTaskStatus(req.ID, "error", 0)
	} else {
		a.repository.UpdateTaskStatus(req.ID, "completed", req.Result)
	}

	w.WriteHeader(http.StatusOK)
}
