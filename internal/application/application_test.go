package application_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/zalhui/calc_golang/internal/application"
)

func TestCalcHandler(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		body       interface{}
		statusCode int
		response   interface{}
	}{
		{
			name:       "valid",
			method:     "POST",
			body:       map[string]string{"expression": "3+5*2"},
			statusCode: http.StatusOK,
			response:   map[string]interface{}{"result": 13.0},
		},
		{
			name:       "brackets",
			method:     "POST",
			body:       map[string]string{"expression": "3+(5*2"},
			statusCode: http.StatusUnprocessableEntity,
			response:   map[string]interface{}{"error": "Expression is not valid. Number of brackets doesn't match"},
		},
		{
			name:       "values",
			method:     "POST",
			body:       map[string]string{"expression": "3+"},
			statusCode: http.StatusUnprocessableEntity,
			response:   map[string]interface{}{"error": "Expression is not valid. Not enough values"},
		},
		{
			name:       "division by zero",
			method:     "POST",
			body:       map[string]string{"expression": "3/0"},
			statusCode: http.StatusUnprocessableEntity,
			response:   map[string]interface{}{"error": "Expression is not valid. Division by zero"},
		},
		{
			name:       "letters",
			method:     "POST",
			body:       map[string]string{"expression": "3+abc"},
			statusCode: http.StatusUnprocessableEntity,
			response:   map[string]interface{}{"error": "Expression is not valid. Only numbers and ( ) + - * / allowed"},
		},
		{
			name:       "method",
			method:     "GET",
			body:       nil,
			statusCode: http.StatusMethodNotAllowed,
			response:   map[string]interface{}{"error": "only POST method allowed"},
		},
		{
			name:       "bad request",
			method:     "POST",
			body:       "invalid json",
			statusCode: http.StatusBadRequest,
			response:   map[string]interface{}{"error": "Bad request"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var reqBody []byte
			if test.body != nil {
				reqBody, _ = json.Marshal(test.body)
			}
			req := httptest.NewRequest(test.method, "/", bytes.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")

			rec := httptest.NewRecorder()
			handler := http.HandlerFunc(application.CalcHandler)
			handler.ServeHTTP(rec, req)

			if rec.Code != test.statusCode {
				t.Errorf("Expected status code %d, got %d", test.statusCode, rec.Code)
			}

			var actualResponse map[string]interface{}
			json.Unmarshal(rec.Body.Bytes(), &actualResponse)
			expectedResponse := test.response.(map[string]interface{})
			for key, value := range expectedResponse {
				if actualResponse[key] != value {
					t.Errorf("Expected %s: %v, got %v", key, value, actualResponse[key])
				}
			}
		})
	}
}
