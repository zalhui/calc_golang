package worker

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
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

			result, err := performOperation(task.Arg1, task.Arg2, task.Operation)
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
