package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/zalhui/calc_golang/internal/orchestrator/application"
)

func main() {
	app := application.New()

	// Создаем маршрутизатор
	router := mux.NewRouter()

	// Регистрируем обработчики с маршрутами
	router.HandleFunc("/api/v1/calculate", app.AddExpressionHandler).Methods("POST")
	router.HandleFunc("/api/v1/expressions", app.GetAllExpressionsHandler).Methods("GET")
	router.HandleFunc("/api/v1/expressions/{id}", app.GetExpressionByIDHandler).Methods("GET") // Используем {id} для динамического параметра
	router.HandleFunc("/internal/task", app.GetPendingTaskHandler).Methods("GET")
	router.HandleFunc("/internal/task/result", app.SubmitTaskResultHandler).Methods("POST", "GET")

	log.Printf("Orchestrator started on :8080")
	log.Fatal(http.ListenAndServe(":8080", router)) // Используем маршрутизатор вместо nil
}
