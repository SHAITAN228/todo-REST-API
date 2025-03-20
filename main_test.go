package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"todo-app/database"
	"todo-app/handlers"
)

// Тест для проверки успешного создания задачи
func TestPostTaskHandler(t *testing.T) {
	database.InitDB()
	taskJSON := `{"title": "Test Task", "content": "Test description", "completed": false}`
	req := httptest.NewRequest("POST", "/tasks", bytes.NewBufferString(taskJSON))
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	handlers.PostTaskHandler(recorder, req)

	if recorder.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, recorder.Code)
	}

	var response map[string]int64
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	_, ok := response["id"]
	if !ok || response["id"] <= 0 {
		t.Errorf("Expected valid id in response, got %+v", response)
	}

}
