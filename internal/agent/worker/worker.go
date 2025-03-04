package worker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/zalhui/calc_golang/config"
	"github.com/zalhui/calc_golang/internal/orchestrator/models"
	"github.com/zalhui/calc_golang/pkg/calculation"
)

var cfg = config.LoadConfig()

func StartWorker() {
	for {
		resp, err := http.Get("http://localhost:8080/internal/task")
		if err != nil {
			log.Printf("Error getting task: %v\n", err)
			time.Sleep(time.Second)
			continue
		}

		if resp.StatusCode == http.StatusOK {
			var task models.Task
			err := json.NewDecoder(resp.Body).Decode(&task)
			if err != nil {
				log.Printf("Error decoding task: %v\n", err)
				continue
			}
			fmt.Printf("Получена задача: %+v\n", task)

			// Обработка зависимостей
			arg1, err := resolveArg(task.Arg1)
			if err != nil {
				log.Printf("Error resolving arg1: %v\n", err)
				submitError(task.ID, err.Error())
				continue
			}

			arg2, err := resolveArg(task.Arg2)
			if err != nil {
				log.Printf("Error resolving arg2: %v\n", err)
				submitError(task.ID, err.Error())
				continue
			}

			// Выполнение операции
			result, err := performOperation(arg1, arg2, task.Operation)
			if err != nil {
				log.Printf("Error performing operation: %v\n", err)
				submitError(task.ID, err.Error())
			} else {
				submitResult(task.ID, result)
			}
		} else if resp.StatusCode == http.StatusNoContent {
			log.Println("No tasks found, waiting for 1 second...")
			time.Sleep(time.Second)
		} else {
			log.Printf("Unexpected status code: %d\n", resp.StatusCode)
			time.Sleep(time.Second)
		}
	}
}

// resolveArg обрабатывает аргумент задачи
func resolveArg(arg string) (float64, error) {
	if isPlaceholder(arg) {
		// Если аргумент — это плейсхолдер, ждем результат задачи
		taskID := strings.TrimSuffix(strings.TrimPrefix(arg, "task_"), "_result")
		return waitForTaskResult(taskID)
	}

	// Если аргумент — это число, преобразуем его
	return strconv.ParseFloat(arg, 64)
}

// isPlaceholder проверяет, является ли аргумент плейсхолдером
func isPlaceholder(arg string) bool {
	return strings.HasPrefix(arg, "task_") && strings.HasSuffix(arg, "_result")
}
func waitForTaskResult(taskID string) (float64, error) {
	for {
		resp, err := http.Get("http://localhost:8080/internal/task/result?id=" + taskID)
		if err != nil {
			log.Printf("Error getting task result: %v\n", err)
			time.Sleep(500 * time.Millisecond) // Ждем перед повторной попыткой
			continue
		}

		if resp.StatusCode == http.StatusOK {
			var resultData struct {
				Result float64 `json:"result"`
			}
			err := json.NewDecoder(resp.Body).Decode(&resultData)
			resp.Body.Close()
			if err != nil {
				log.Printf("Error decoding result: %v\n", err)
				return 0, err
			}
			return resultData.Result, nil
		} else if resp.StatusCode == http.StatusNotFound {
			log.Printf("Task %s result not ready, waiting...\n", taskID)
			time.Sleep(500 * time.Millisecond)
		} else {
			log.Printf("Unexpected status code when getting result: %d\n", resp.StatusCode)
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func performOperation(arg1, arg2 float64, operation string) (float64, error) {
	var result float64
	var err error

	switch operation {
	case "+":
		result = arg1 + arg2
		<-time.After(cfg.TimeAddition)
	case "-":
		result = arg1 - arg2
		<-time.After(cfg.TimeSubtraction)
	case "*":
		result = arg1 * arg2
		<-time.After(cfg.TimeMultiplication)
	case "/":
		if arg2 == 0 {
			return 0, calculation.ErrDivisionByZero
		}
		result = arg1 / arg2
		<-time.After(cfg.TimeDivision)
	default:
		return 0, calculation.ErrAllowed
	}

	return result, err
}

func submitResult(taskID string, result float64) {
	data := map[string]interface{}{
		"id":     taskID,
		"result": result,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshaling result: %v\n", err)
		return
	}

	resp, err := http.Post("http://localhost:8080/internal/task/result", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error submitting result: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Unexpected status code when submitting result: %d\n", resp.StatusCode)
	}
}

func submitError(taskID string, errorMsg string) {
	data := map[string]interface{}{
		"id":    taskID,
		"error": errorMsg,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshaling error: %v\n", err)
		return
	}

	resp, err := http.Post("http://localhost:8080/internal/task/result", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error submitting error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Unexpected status code when submitting error: %d\n", resp.StatusCode)
	}
}
