package model

import (
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Todo represents a task item
type Todo struct {
	ID          uint      `gorm:"primaryKey;autoIncrement:false" json:"id"` // Manually managed ID per day
	Title       string    `json:"title"`
	Description string    `json:"description"`
	IsCompleted bool      `json:"is_completed"`
	ParentID    *uint     `json:"parent_id"`              // Refers to ID within the same Date
	Date        string    `gorm:"primaryKey" json:"date"` // Format: YYYY-MM-DD
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

// GetNextID calculates the next ID for a given date
func GetNextID(date string) (uint, error) {
	var maxID *uint
	err := DB.Model(&Todo{}).Where("date = ?", date).Select("MAX(id)").Scan(&maxID).Error
	if err != nil {
		return 0, err
	}
	if maxID == nil {
		return 1, nil
	}
	return *maxID + 1, nil
}

// CreateTodo creates a new todo with manual ID generation
func CreateTodo(title string, description string, date string) (*Todo, error) {
	id, err := GetNextID(date)
	if err != nil {
		return nil, err
	}

	todo := &Todo{
		ID:          id,
		Title:       title,
		Description: description,
		IsCompleted: false,
		Date:        date,
	}
	result := DB.Create(todo)
	return todo, result.Error
}

// GetTodosByDate retrieves todos for a specific date with hierarchical sorting
func GetTodosByDate(date string) ([]Todo, error) {
	var todos []Todo
	// Fetch all todos for the date
	if err := DB.Where("date = ?", date).Find(&todos).Error; err != nil {
		return nil, err
	}

	// Separate root todos and child todos
	var roots []Todo
	childrenMap := make(map[uint][]Todo)

	for _, t := range todos {
		if t.ParentID == nil {
			roots = append(roots, t)
		} else {
			pid := *t.ParentID
			childrenMap[pid] = append(childrenMap[pid], t)
		}
	}

	// Sort roots: Pending first, then Completed. Secondary sort by ID.
	sortTodos(roots)

	// Sort children map
	for pid := range childrenMap {
		sortTodos(childrenMap[pid])
	}

	// Reconstruct the list
	var result []Todo
	for _, root := range roots {
		result = append(result, root)
		if children, ok := childrenMap[root.ID]; ok {
			result = append(result, children...)
		}
	}

	return result, nil
}

func sortTodos(todos []Todo) {
	// Simple bubble sort for stability and small lists
	for i := 0; i < len(todos); i++ {
		for j := 0; j < len(todos)-i-1; j++ {
			// Sort criteria:
			// 1. IsCompleted (false < true)
			// 2. ID (asc)

			swap := false
			if todos[j].IsCompleted && !todos[j+1].IsCompleted {
				swap = true
			} else if todos[j].IsCompleted == todos[j+1].IsCompleted {
				if todos[j].ID > todos[j+1].ID {
					swap = true
				}
			}

			if swap {
				todos[j], todos[j+1] = todos[j+1], todos[j]
			}
		}
	}
}

// GetTodoByID retrieves a todo by Date and ID
func GetTodoByID(date string, id uint) (*Todo, error) {
	var todo Todo
	if err := DB.Where("date = ? AND id = ?", date, id).First(&todo).Error; err != nil {
		return nil, err
	}
	return &todo, nil
}

// UpdateTodo updates an existing todo
func UpdateTodo(todo *Todo) error {
	// Ensure we save using the composite key
	return DB.Save(todo).Error
}

// DeleteTodo deletes a todo by Date and ID
func DeleteTodo(date string, id uint) error {
	return DB.Where("date = ? AND id = ?", date, id).Delete(&Todo{}).Error
}

// ToggleTodoStatus toggles the completion status of a todo
func ToggleTodoStatus(date string, id uint) (*Todo, error) {
	var todo Todo
	if err := DB.Where("date = ? AND id = ?", date, id).First(&todo).Error; err != nil {
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

// GetTotalCount returns the total number of todos for a specific date
func GetTotalCount(date string) (int64, error) {
	var count int64
	if err := DB.Model(&Todo{}).Where("date = ?", date).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// HasChildren checks if a todo has any child tasks (within the same date)
func HasChildren(date string, id uint) (bool, error) {
	var count int64
	if err := DB.Model(&Todo{}).Where("date = ? AND parent_id = ?", date, id).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
