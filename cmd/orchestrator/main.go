package main

import (
	"log"
	"net/http"

	"github.com/zalhui/calc_golang/internal/orchestrator/application"
)

func main() {

	app := application.New()
	http.HandleFunc("/api/v1/calculate", app.AddExpressionHandler)
	http.HandleFunc("/api/v1/expressions", app.GetAllExpressionsHandler)
	http.HandleFunc("/api/v1/expressions/", app.GetExpressionByIDHandler)
	http.HandleFunc("/internal/task", app.GetPendingTaskHandler)
	http.HandleFunc("/internal/task/result", app.SubmitTaskResultHandler)

	log.Printf("Orchestrator started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
