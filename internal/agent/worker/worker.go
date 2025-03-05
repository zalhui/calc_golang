package worker

import (
	"bytes"
	"encoding/json"
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

type TaskResponse struct {
	Task models.Task `json:"task"`
}

func StartWorker() {
	for {
		resp, err := http.Get("http://localhost:8080/internal/task")
		if err != nil {
			log.Printf("Error getting task: %v", err)
			time.Sleep(time.Second)
			continue
		}

		if resp.StatusCode == http.StatusOK {
			var response TaskResponse
			if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
				log.Printf("Error decoding task response: %v", err)
				continue
			}
			task := response.Task
			log.Printf("Received task: ID=%s, ExpressionID=%s, Arg1=%s, Arg2=%s, Operation=%s, Status=%s",
				task.ID, task.ExpressionID, task.Arg1, task.Arg2, task.Operation, task.Status)

			if task.ID == "" || task.Operation == "" || task.Arg1 == "" || task.Arg2 == "" {
				log.Printf("Received invalid task with empty fields: %+v", task)
				continue
			}

			arg1, err := resolveArg(task.Arg1)
			if err != nil {
				log.Printf("Error resolving arg1 for task %s: %v", task.ID, err)
				submitError(task.ID, err.Error())
				continue
			}
			log.Printf("Resolved arg1 for task %s: %f", task.ID, arg1)

			arg2, err := resolveArg(task.Arg2)
			if err != nil {
				log.Printf("Error resolving arg2 for task %s: %v", task.ID, err)
				submitError(task.ID, err.Error())
				continue
			}
			log.Printf("Resolved arg2 for task %s: %f", task.ID, arg2)

			result, err := performOperation(arg1, arg2, task.Operation)
			if err != nil {
				log.Printf("Error performing operation for task %s: %v", task.ID, err)
				submitError(task.ID, err.Error())
			} else {
				log.Printf("Operation completed for task %s: %s %f %s %f = %f",
					task.ID, task.Arg1, arg1, task.Operation, arg2, result)
				submitResult(task.ID, result)
			}
		} else if resp.StatusCode == http.StatusNotFound {
			log.Println("No tasks found, waiting for 1 second...")
			time.Sleep(time.Second)
		} else {
			log.Printf("Unexpected status code: %d", resp.StatusCode)
			time.Sleep(time.Second)
		}
	}
}

func resolveArg(arg string) (float64, error) {
	if isPlaceholder(arg) {
		taskID := strings.TrimSuffix(strings.TrimPrefix(arg, "task_"), "_result")
		log.Printf("Resolving placeholder %s for task %s", arg, taskID)
		result, err := waitForTaskResult(taskID)
		if err != nil {
			log.Printf("Failed to resolve placeholder %s: %v", arg, err)
			return 0, err
		}
		return result, nil
	}

	result, err := strconv.ParseFloat(arg, 64)
	if err != nil {
		log.Printf("Failed to parse argument %s as float64: %v", arg, err)
		return 0, err
	}
	return result, nil
}

func isPlaceholder(arg string) bool {
	return strings.HasPrefix(arg, "task_") && strings.HasSuffix(arg, "_result")
}
func waitForTaskResult(taskID string) (float64, error) {
	for {
		resp, err := http.Get("http://localhost:8080/internal/task/result?id=" + taskID)
		if err != nil {
			log.Printf("Error getting task result: %v\n", err)
			time.Sleep(500 * time.Millisecond)
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
		result = arg2 - arg1
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
		log.Printf("Error marshaling result: %v", err)
		return
	}

	resp, err := http.Post("http://localhost:8080/internal/task/result", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error submitting result for task %s: %v", taskID, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to submit result for task %s, status code: %d", taskID, resp.StatusCode)
	} else {
		log.Printf("Successfully submitted result for task %s: %f", taskID, result)
	}
}

func submitError(taskID string, errorMsg string) {
	data := map[string]interface{}{
		"id":    taskID,
		"error": errorMsg,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshaling error: %v", err)
		return
	}

	resp, err := http.Post("http://localhost:8080/internal/task/result", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error submitting error for task %s: %v", taskID, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to submit error for task %s, status code: %d", taskID, resp.StatusCode)
	} else {
		log.Printf("Successfully submitted error for task %s: %s", taskID, errorMsg)
	}
}
