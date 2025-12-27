package service

import (
	"fmt"
	"time"

	"itodo/internal/model"
)

// TodoService handles business logic for todos
type TodoService struct{}

// NewTodoService creates a new TodoService
func NewTodoService() *TodoService {
	return &TodoService{}
}

// Add creates a new todo and returns a log message
func (s *TodoService) Add(title string, description string, date string) (string, error) {
	if title == "" {
		return "", fmt.Errorf("title cannot be empty")
	}

	// 1. Validate Date: Cannot create todos for past dates
	today := s.GetCurrentDate()
	if date < today {
		return "", fmt.Errorf("cannot create tasks for past dates")
	}

	// 2. Validate Count: Limit daily tasks (e.g., 50)
	count, err := model.GetTotalCount(date)
	if err != nil {
		return "", err
	}
	if count >= 50 {
		return "", fmt.Errorf("daily task limit (50) reached")
	}

	_, err = model.CreateTodo(title, description, date)
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
func (s *TodoService) Toggle(date string, id uint) (string, error) {
	todo, err := model.ToggleTodoStatus(date, id)
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
func (s *TodoService) Delete(date string, id uint) (string, error) {
	todo, err := model.GetTodoByID(date, id)
	if err != nil {
		return "", err
	}
	if err := model.DeleteTodo(date, id); err != nil {
		return "", err
	}
	return fmt.Sprintf("Deleted todo: %s", todo.Title), nil
}

// Edit updates a todo title and description
func (s *TodoService) Edit(date string, id uint, newTitle string, newDescription string) (string, error) {
	todo, err := model.GetTodoByID(date, id)
	if err != nil {
		return "", err
	}
	oldTitle := todo.Title
	todo.Title = newTitle
	todo.Description = newDescription
	if err := model.UpdateTodo(todo); err != nil {
		return "", err
	}
	return fmt.Sprintf("Updated todo '%s'", oldTitle), nil
}

// GetWeekRange returns the start (Monday) and end (Sunday) dates for the week of the given date
func (s *TodoService) GetWeekRange(dateStr string) (string, string, error) {
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return "", "", err
	}

	// Calculate offset to Monday (0 = Monday, ..., 6 = Sunday)
	// time.Weekday: Sunday=0, Monday=1, ...
	weekday := t.Weekday()
	offset := int(weekday) - 1
	if weekday == time.Sunday {
		offset = 6
	}

	start := t.AddDate(0, 0, -offset)
	end := start.AddDate(0, 0, 6)

	return start.Format("2006-01-02"), end.Format("2006-01-02"), nil
}

// GetTasksByRange returns root todos within the given date range
func (s *TodoService) GetTasksByRange(startDate, endDate string) ([]model.Todo, error) {
	return model.GetRootTodosByDateRange(startDate, endDate)
}

// IndentTodo makes the current todo a child of the previous todo
func (s *TodoService) IndentTodo(date string, currentID uint, prevID uint) (string, error) {
	if prevID == 0 {
		return "", fmt.Errorf("no previous todo to indent under")
	}

	curr, err := model.GetTodoByID(date, currentID)
	if err != nil {
		return "", err
	}

	// Check if current todo has children. If so, prevent indentation.
	hasChildren, err := model.HasChildren(date, currentID)
	if err != nil {
		return "", err
	}
	if hasChildren {
		return "", fmt.Errorf("cannot indent a task that has sub-tasks")
	}

	prev, err := model.GetTodoByID(date, prevID)
	if err != nil {
		return "", err
	}

	// Logic for Indentation:
	// 1. If 'prev' is a Root task (ParentID == nil):
	//    'curr' becomes a child of 'prev'.
	// 2. If 'prev' is already a Child task (ParentID != nil):
	//    'curr' becomes a sibling of 'prev' (adopts prev.ParentID).
	//    This prevents nesting deeper than 1 level, but allows multiple children under one parent.

	if prev.ParentID == nil {
		// Case 1: Indent under root
		curr.ParentID = &prev.ID
	} else {
		// Case 2: Indent as sibling (same parent as prev)
		curr.ParentID = prev.ParentID
	}

	if err := model.UpdateTodo(curr); err != nil {
		return "", err
	}

	return fmt.Sprintf("Indented todo '%s'", curr.Title), nil
}

// OutdentTodo makes the current todo a root item
func (s *TodoService) OutdentTodo(date string, id uint) (string, error) {
	todo, err := model.GetTodoByID(date, id)
	if err != nil {
		return "", err
	}

	if todo.ParentID == nil {
		return "", fmt.Errorf("todo is already at root level")
	}

	todo.ParentID = nil
	if err := model.UpdateTodo(todo); err != nil {
		return "", err
	}

	return fmt.Sprintf("Outdented todo '%s'", todo.Title), nil
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
