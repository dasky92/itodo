package service

import (
	"itodo/model"
	"testing"
)

func setupTestDB(t *testing.T) {
	err := model.InitDB(":memory:")
	if err != nil {
		t.Fatalf("Failed to init test db: %v", err)
	}
}

func TestTodoService_Add(t *testing.T) {
	setupTestDB(t)
	svc := NewTodoService()
	date := "2023-10-01"

	msg, err := svc.Add("Test Todo", date)
	if err != nil {
		t.Errorf("Add failed: %v", err)
	}
	if msg != "Created todo: Test Todo" {
		t.Errorf("Unexpected log message: %s", msg)
	}

	todos, err := svc.List(date)
	if err != nil {
		t.Errorf("List failed: %v", err)
	}
	if len(todos) != 1 {
		t.Errorf("Expected 1 todo, got %d", len(todos))
	}
	if todos[0].Title != "Test Todo" {
		t.Errorf("Expected title 'Test Todo', got '%s'", todos[0].Title)
	}
}

func TestTodoService_Toggle(t *testing.T) {
	setupTestDB(t)
	svc := NewTodoService()
	date := "2023-10-01"

	svc.Add("Test Todo", date)
	todos, _ := svc.List(date)
	id := todos[0].ID

	// Toggle to completed
	msg, err := svc.Toggle(id)
	if err != nil {
		t.Errorf("Toggle failed: %v", err)
	}
	// Check log message contains "completed"
	if msg == "" {
		t.Error("Expected log message")
	}

	todos, _ = svc.List(date)
	if !todos[0].IsCompleted {
		t.Error("Todo should be completed")
	}

	// Toggle back to pending
	svc.Toggle(id)
	todos, _ = svc.List(date)
	if todos[0].IsCompleted {
		t.Error("Todo should be pending")
	}
}

func TestTodoService_Delete(t *testing.T) {
	setupTestDB(t)
	svc := NewTodoService()
	date := "2023-10-01"

	svc.Add("To Delete", date)
	todos, _ := svc.List(date)
	id := todos[0].ID

	msg, err := svc.Delete(id)
	if err != nil {
		t.Errorf("Delete failed: %v", err)
	}
	if msg == "" {
		t.Error("Expected log message")
	}

	todos, _ = svc.List(date)
	if len(todos) != 0 {
		t.Error("Todo should be deleted")
	}
}

func TestTodoService_Edit(t *testing.T) {
	setupTestDB(t)
	svc := NewTodoService()
	date := "2023-10-01"

	svc.Add("Old Title", date)
	todos, _ := svc.List(date)
	id := todos[0].ID

	msg, err := svc.Edit(id, "New Title")
	if err != nil {
		t.Errorf("Edit failed: %v", err)
	}
	if msg == "" {
		t.Error("Expected log message")
	}

	todos, _ = svc.List(date)
	if todos[0].Title != "New Title" {
		t.Errorf("Expected 'New Title', got '%s'", todos[0].Title)
	}
}

func TestTodoService_GetStats(t *testing.T) {
	setupTestDB(t)
	svc := NewTodoService()
	date := "2023-10-01"

	svc.Add("Todo 1", date)
	svc.Add("Todo 2", date)

	todos, _ := svc.List(date)
	svc.Toggle(todos[0].ID) // Complete one

	completed, pending, err := svc.GetStats(date)
	if err != nil {
		t.Errorf("GetStats failed: %v", err)
	}

	if completed != 1 {
		t.Errorf("Expected 1 completed, got %d", completed)
	}
	if pending != 1 {
		t.Errorf("Expected 1 pending, got %d", pending)
	}
}
