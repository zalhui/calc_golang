package application

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/zalhui/calc_golang/internal/auth"
)

func (a *Application) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	if req.Login == "" || req.Password == "" {
		http.Error(w, "Login and password are required", http.StatusBadRequest)
		return
	}

	_, err := auth.RegisterUser(a.db, req.Login, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (a *Application) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	token, err := auth.LoginUser(a.db, req.Login, req.Password)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (a *Application) AddExpressionHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		Expression string `json:"expression"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	expressionID, err := a.AddExpression(req.Expression, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"id":      expressionID,
		"status":  "pending",
		"message": "Expression accepted for processing",
	})
}

func (a *Application) GetExpressionByIDHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	expressionID := vars["id"]
	if expressionID == "" {
		http.Error(w, "Missing expression ID", http.StatusBadRequest)
		return
	}

	expression, exists := a.repository.GetExpressionByID(expressionID, userID)
	if !exists {
		http.Error(w, "Expression not found", http.StatusNotFound)
		return
	}
	var result interface{}
	if expression.Result.Valid {
		result = expression.Result.Float64
	} else {
		result = nil
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":      expression.ID,
		"status":  expression.Status,
		"result":  result,
		"created": expression.CreatedAt,
	})
}

func (a *Application) GetAllExpressionsHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	expressions := a.repository.GetAllExpressions(userID)

	w.Header().Set("Content-Type", "application/json")
	response := make([]map[string]interface{}, 0, len(expressions))
	for _, expr := range expressions {
		response = append(response, map[string]interface{}{
			"id":      expr.ID,
			"status":  expr.Status,
			"result":  expr.Result,
			"created": expr.CreatedAt,
		})
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"expressions": response})
}

func (a *Application) GetPendingTaskHandler(w http.ResponseWriter, r *http.Request) {
	task, exists := a.repository.GetPendingTask()
	if !exists {
		http.Error(w, "No tasks available", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"task": map[string]interface{}{
			"id":            task.ID,
			"expression_id": task.ExpressionID,
			"arg1":          task.Arg1,
			"arg2":          task.Arg2,
			"operation":     task.Operation,
		},
	})
}

func (a *Application) SubmitTaskResultHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		taskID := r.URL.Query().Get("id")
		if taskID == "" {
			http.Error(w, "Missing task ID", http.StatusBadRequest)
			return
		}

		task, found := a.repository.GetTaskByID(taskID)
		if !found || task.Status != "completed" {
			http.Error(w, "Result not ready", http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"result": task.Result,
		})
		return
	}

	var req struct {
		ID     string  `json:"id"`
		Result float64 `json:"result,omitempty"`
		Error  string  `json:"error,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	if req.Error != "" {
		a.repository.UpdateTaskStatus(req.ID, "error", 0)
	} else {
		a.repository.UpdateTaskStatus(req.ID, "completed", req.Result)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (a *Application) GetUserHistoryHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	history, err := a.repository.GetUserHistory(userID)
	if err != nil {
		http.Error(w, "Failed to get history", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"history": history})

}
