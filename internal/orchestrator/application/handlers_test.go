package application

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAddExpressionHandler(t *testing.T) {
	app := New()
	reqBody, _ := json.Marshal(map[string]string{"expression": "2+2"})
	req := httptest.NewRequest("POST", "/api/v1/calculate", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	app.AddExpressionHandler(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", rr.Code)
	}

	var resp map[string]string
	json.Unmarshal(rr.Body.Bytes(), &resp)
	if _, ok := resp["id"]; !ok {
		t.Errorf("Expected response to contain 'id', got %v", resp)
	}
}
