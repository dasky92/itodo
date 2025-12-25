package model

import (
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Todo represents a task item
type Todo struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Title       string    `json:"title"`
	IsCompleted bool      `json:"is_completed"`
	Date        string    `gorm:"index" json:"date"` // Format: YYYY-MM-DD
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// DB is the database instance
var DB *gorm.DB

// InitDB initializes the SQLite database
func InitDB(dbPath string) error {
	var err error
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return err
	}

	// Auto migrate the schema
	err = DB.AutoMigrate(&Todo{})
	return err
}

// CreateTodo creates a new todo
func CreateTodo(title string, date string) (*Todo, error) {
	todo := &Todo{
		Title:       title,
		IsCompleted: false,
		Date:        date,
	}
	result := DB.Create(todo)
	return todo, result.Error
}

// GetTodosByDate retrieves todos for a specific date
func GetTodosByDate(date string) ([]Todo, error) {
	var todos []Todo
	result := DB.Where("date = ?", date).Find(&todos)
	return todos, result.Error
}

// GetTodoByID retrieves a todo by ID
func GetTodoByID(id uint) (*Todo, error) {
	var todo Todo
	if err := DB.First(&todo, id).Error; err != nil {
		return nil, err
	}
	return &todo, nil
}

// UpdateTodo updates an existing todo
func UpdateTodo(todo *Todo) error {
	return DB.Save(todo).Error
}

// DeleteTodo deletes a todo by ID
func DeleteTodo(id uint) error {
	return DB.Delete(&Todo{}, id).Error
}

// ToggleTodoStatus toggles the completion status of a todo
func ToggleTodoStatus(id uint) (*Todo, error) {
	var todo Todo
	if err := DB.First(&todo, id).Error; err != nil {
		return nil, err
	}
	todo.IsCompleted = !todo.IsCompleted
	if err := DB.Save(&todo).Error; err != nil {
		return nil, err
	}
	return &todo, nil
}

// GetStats returns the count of completed and pending todos for a specific date
func GetStats(date string) (int, int, error) {
	var completed int64
	var pending int64

	if err := DB.Model(&Todo{}).Where("date = ? AND is_completed = ?", date, true).Count(&completed).Error; err != nil {
		return 0, 0, err
	}
	if err := DB.Model(&Todo{}).Where("date = ? AND is_completed = ?", date, false).Count(&pending).Error; err != nil {
		return 0, 0, err
	}

	return int(completed), int(pending), nil
}
