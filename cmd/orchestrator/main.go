package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"github.com/zalhui/calc_golang/internal/db"
	"github.com/zalhui/calc_golang/internal/middleware"
	"github.com/zalhui/calc_golang/internal/orchestrator/application"
)

func main() {
	// Инициализация базы данных
	database, err := db.NewDB("calc.db")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	// Создание экземпляра приложения
	app := application.New(database)

	// Настройка маршрутизатора
	router := mux.NewRouter()

	// Публичные эндпоинты
	publicRouter := router.PathPrefix("/api/v1").Subrouter()
	publicRouter.HandleFunc("/register", app.RegisterHandler).Methods("POST")
	publicRouter.HandleFunc("/login", app.LoginHandler).Methods("POST")

	// Защищенные эндпоинты
	protectedRouter := router.PathPrefix("/api/v1").Subrouter()
	protectedRouter.Use(middleware.JWTAuthMiddleware)

	protectedRouter.HandleFunc("/calculate", app.AddExpressionHandler).Methods("POST")
	protectedRouter.HandleFunc("/expressions", app.GetAllExpressionsHandler).Methods("GET")
	protectedRouter.HandleFunc("/expressions/{id}", app.GetExpressionByIDHandler).Methods("GET")
	protectedRouter.HandleFunc("/history", app.GetUserHistoryHandler).Methods("GET")

	// Внутренние эндпоинты для агентов
	internalRouter := router.PathPrefix("/internal").Subrouter()
	internalRouter.HandleFunc("/task", app.GetPendingTaskHandler).Methods("GET")
	internalRouter.HandleFunc("/task/result", app.SubmitTaskResultHandler).Methods("POST", "GET")

	// Настройка HTTP сервера
	server := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		log.Println("Shutting down server...")

		shutdownCtx, shutdownCancel := context.WithTimeout(ctx, 5*time.Second)
		defer shutdownCancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Printf("Server shutdown error: %v", err)
		}
	}()

	// Запуск сервера
	log.Println("Starting orchestrator on :8080")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed: %v", err)
	}

	log.Println("Server stopped gracefully")
}
