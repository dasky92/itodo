package service

import (
	"fmt"
	"time"

	"itodo/model"
)

// TodoService handles business logic for todos
type TodoService struct{}

// NewTodoService creates a new TodoService
func NewTodoService() *TodoService {
	return &TodoService{}
}

// Add creates a new todo and returns a log message
func (s *TodoService) Add(title string, date string) (string, error) {
	if title == "" {
		return "", fmt.Errorf("title cannot be empty")
	}
	_, err := model.CreateTodo(title, date)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Created todo: %s", title), nil
}

// List returns todos for a date
func (s *TodoService) List(date string) ([]model.Todo, error) {
	return model.GetTodosByDate(date)
}

// Toggle toggles the status of a todo
func (s *TodoService) Toggle(id uint) (string, error) {
	todo, err := model.ToggleTodoStatus(id)
	if err != nil {
		return "", err
	}
	status := "pending"
	if todo.IsCompleted {
		status = "completed"
	}
	return fmt.Sprintf("Toggled todo '%s' to %s", todo.Title, status), nil
}

// Delete removes a todo
func (s *TodoService) Delete(id uint) (string, error) {
	todo, err := model.GetTodoByID(id)
	if err != nil {
		return "", err
	}
	if err := model.DeleteTodo(id); err != nil {
		return "", err
	}
	return fmt.Sprintf("Deleted todo: %s", todo.Title), nil
}

// Edit updates a todo title
func (s *TodoService) Edit(id uint, newTitle string) (string, error) {
	todo, err := model.GetTodoByID(id)
	if err != nil {
		return "", err
	}
	oldTitle := todo.Title
	todo.Title = newTitle
	if err := model.UpdateTodo(todo); err != nil {
		return "", err
	}
	return fmt.Sprintf("Updated todo '%s' to '%s'", oldTitle, newTitle), nil
}

// GetStats returns stats for a date
func (s *TodoService) GetStats(date string) (int, int, error) {
	return model.GetStats(date)
}

// GetCurrentDate returns today's date in YYYY-MM-DD format
func (s *TodoService) GetCurrentDate() string {
	return time.Now().Format("2006-01-02")
}

// GetNextDate returns the next day
func (s *TodoService) GetNextDate(dateStr string) string {
	t, _ := time.Parse("2006-01-02", dateStr)
	return t.AddDate(0, 0, 1).Format("2006-01-02")
}

// GetPrevDate returns the previous day
func (s *TodoService) GetPrevDate(dateStr string) string {
	t, _ := time.Parse("2006-01-02", dateStr)
	return t.AddDate(0, 0, -1).Format("2006-01-02")
}

// FormatDateForDisplay returns a friendly date string (e.g. "2023-10-27 Friday")
func (s *TodoService) FormatDateForDisplay(dateStr string) string {
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return dateStr
	}
	return t.Format("2006-01-02 Monday")
}
