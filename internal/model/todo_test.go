package model

import (
	"errors"
	"path/filepath"
	"testing"

	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) {
	t.Helper()
	if err := InitDB(":memory:"); err != nil {
		t.Fatalf("Failed to init test db: %v", err)
	}
}

func createTodo(t *testing.T, title, desc, date string) *Todo {
	t.Helper()
	todo, err := CreateTodo(title, desc, date)
	if err != nil {
		t.Fatalf("CreateTodo: %v", err)
	}
	return todo
}

func TestInitDB(t *testing.T) {
	t.Run("memory", func(t *testing.T) {
		if err := InitDB(":memory:"); err != nil {
			t.Fatalf("InitDB(:memory:): %v", err)
		}
		_, err := CreateTodo("x", "", "2024-01-01")
		if err != nil {
			t.Errorf("CreateTodo after InitDB(:memory:): %v", err)
		}
	})

	t.Run("file_path_creates_dirs", func(t *testing.T) {
		dir := t.TempDir()
		dbPath := filepath.Join(dir, "a", "b", "db.sqlite")
		if err := InitDB(dbPath); err != nil {
			t.Fatalf("InitDB(file): %v", err)
		}
		todo, err := CreateTodo("y", "", "2024-01-02")
		if err != nil {
			t.Fatalf("CreateTodo after InitDB(file): %v", err)
		}
		got, err := GetTodoByID("2024-01-02", todo.ID)
		if err != nil {
			t.Fatalf("GetTodoByID: %v", err)
		}
		if got.Title != "y" {
			t.Errorf("got Title %q, want %q", got.Title, "y")
		}
	})
}

func TestGetNextID(t *testing.T) {
	setupTestDB(t)
	dateA := "2024-01-01"
	dateB := "2024-01-02"

	t.Run("empty_returns_1", func(t *testing.T) {
		got, err := GetNextID(dateA)
		if err != nil {
			t.Fatalf("GetNextID: %v", err)
		}
		if got != 1 {
			t.Errorf("GetNextID(empty) = %d, want 1", got)
		}
	})

	t.Run("after_one_returns_2", func(t *testing.T) {
		createTodo(t, "a", "", dateA)
		got, err := GetNextID(dateA)
		if err != nil {
			t.Fatalf("GetNextID: %v", err)
		}
		if got != 2 {
			t.Errorf("GetNextID(after 1) = %d, want 2", got)
		}
	})

	t.Run("per_date_independence", func(t *testing.T) {
		createTodo(t, "x", "", dateA)
		createTodo(t, "y", "", dateA)
		gotA, err := GetNextID(dateA)
		if err != nil {
			t.Fatalf("GetNextID(dateA): %v", err)
		}
		if gotA != 4 {
			t.Errorf("GetNextID(dateA) = %d, want 4 (dateA has 3 todos from prior subtests)", gotA)
		}
		gotB, err := GetNextID(dateB)
		if err != nil {
			t.Fatalf("GetNextID(dateB): %v", err)
		}
		if gotB != 1 {
			t.Errorf("GetNextID(dateB) = %d, want 1 (dateB has no todos)", gotB)
		}
	})
}

func TestGetNextPosition(t *testing.T) {
	setupTestDB(t)
	date := "2024-01-01"

	t.Run("root_empty_returns_0", func(t *testing.T) {
		got, err := GetNextPosition(date, nil)
		if err != nil {
			t.Fatalf("GetNextPosition: %v", err)
		}
		if got != 0 {
			t.Errorf("GetNextPosition(date, nil) = %d, want 0", got)
		}
	})

	t.Run("root_after_one_returns_1", func(t *testing.T) {
		createTodo(t, "r", "", date)
		got, err := GetNextPosition(date, nil)
		if err != nil {
			t.Fatalf("GetNextPosition: %v", err)
		}
		if got != 1 {
			t.Errorf("GetNextPosition(date, nil) = %d, want 1", got)
		}
	})

	t.Run("parent_empty_returns_0", func(t *testing.T) {
		// A is root, B is still root. Get next position under A (no children yet).
		a := createTodo(t, "A", "", date)
		got, err := GetNextPosition(date, &a.ID)
		if err != nil {
			t.Fatalf("GetNextPosition: %v", err)
		}
		if got != 0 {
			t.Errorf("GetNextPosition(date, &A.ID) with no children = %d, want 0", got)
		}
	})

	t.Run("parent_after_one_child_returns_1", func(t *testing.T) {
		a := createTodo(t, "A2", "", date)
		b := createTodo(t, "B2", "", date)
		b.ParentID = &a.ID
		b.Position = 0
		if err := UpdateTodo(b); err != nil {
			t.Fatalf("UpdateTodo: %v", err)
		}
		got, err := GetNextPosition(date, &a.ID)
		if err != nil {
			t.Fatalf("GetNextPosition: %v", err)
		}
		if got != 1 {
			t.Errorf("GetNextPosition(date, &A.ID) with 1 child = %d, want 1", got)
		}
	})
}

func TestCreateTodo(t *testing.T) {
	setupTestDB(t)
	date := "2024-01-01"

	todo, err := CreateTodo("T", "D", date)
	if err != nil {
		t.Fatalf("CreateTodo: %v", err)
	}
	if todo.ID != 1 {
		t.Errorf("ID = %d, want 1", todo.ID)
	}
	if todo.Position != 0 {
		t.Errorf("Position = %d, want 0", todo.Position)
	}
	if todo.Date != date {
		t.Errorf("Date = %q, want %q", todo.Date, date)
	}
	if todo.Title != "T" {
		t.Errorf("Title = %q, want %q", todo.Title, "T")
	}
	if todo.Description != "D" {
		t.Errorf("Description = %q, want %q", todo.Description, "D")
	}
	if todo.IsCompleted {
		t.Error("IsCompleted = true, want false")
	}
	if todo.ParentID != nil {
		t.Error("ParentID != nil, want nil")
	}
}

func TestGetTodoByID(t *testing.T) {
	setupTestDB(t)
	date := "2024-01-01"
	created := createTodo(t, "F", "d", date)

	t.Run("found", func(t *testing.T) {
		got, err := GetTodoByID(date, created.ID)
		if err != nil {
			t.Fatalf("GetTodoByID: %v", err)
		}
		if got.ID != created.ID || got.Title != "F" || got.Description != "d" {
			t.Errorf("GetTodoByID: got %+v", got)
		}
	})

	t.Run("not_found", func(t *testing.T) {
		_, err := GetTodoByID(date, 999)
		if err == nil {
			t.Error("GetTodoByID(999): expected error")
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			t.Errorf("expected gorm.ErrRecordNotFound, got %v", err)
		}
	})
}

func TestUpdateTodo(t *testing.T) {
	setupTestDB(t)
	date := "2024-01-01"
	a := createTodo(t, "A", "", date)
	b := createTodo(t, "B", "d1", date)

	b.Title = "New"
	b.Description = "d2"
	b.ParentID = &a.ID
	if err := UpdateTodo(b); err != nil {
		t.Fatalf("UpdateTodo: %v", err)
	}

	got, err := GetTodoByID(date, b.ID)
	if err != nil {
		t.Fatalf("GetTodoByID: %v", err)
	}
	if got.Title != "New" || got.Description != "d2" || got.ParentID == nil || *got.ParentID != a.ID {
		t.Errorf("GetTodoByID: got Title=%q Desc=%q ParentID=%v", got.Title, got.Description, got.ParentID)
	}
}

func TestDeleteTodo(t *testing.T) {
	setupTestDB(t)
	date := "2024-01-01"
	todo := createTodo(t, "X", "", date)

	if err := DeleteTodo(date, todo.ID); err != nil {
		t.Fatalf("DeleteTodo: %v", err)
	}
	_, err := GetTodoByID(date, todo.ID)
	if err == nil {
		t.Error("GetTodoByID after Delete: expected error")
	}
	c, _ := GetTotalCount(date)
	if c != 0 {
		t.Errorf("GetTotalCount after Delete = %d, want 0", c)
	}
}

func TestToggleTodoStatus(t *testing.T) {
	setupTestDB(t)
	date := "2024-01-01"
	todo := createTodo(t, "T", "", date)

	toggled, err := ToggleTodoStatus(date, todo.ID)
	if err != nil {
		t.Fatalf("ToggleTodoStatus: %v", err)
	}
	if !toggled.IsCompleted {
		t.Error("after first toggle: IsCompleted = false, want true")
	}
	got, _ := GetTodoByID(date, todo.ID)
	if !got.IsCompleted {
		t.Error("GetTodoByID: IsCompleted = false, want true")
	}

	_, err = ToggleTodoStatus(date, todo.ID)
	if err != nil {
		t.Fatalf("ToggleTodoStatus(2): %v", err)
	}
	got, _ = GetTodoByID(date, todo.ID)
	if got.IsCompleted {
		t.Error("after second toggle: IsCompleted = true, want false")
	}
}

func TestGetStats(t *testing.T) {
	setupTestDB(t)
	date := "2024-01-01"

	t.Run("empty_0_0", func(t *testing.T) {
		c, p, err := GetStats(date)
		if err != nil {
			t.Fatalf("GetStats: %v", err)
		}
		if c != 0 || p != 0 {
			t.Errorf("GetStats(empty) = %d completed, %d pending; want 0, 0", c, p)
		}
	})

	t.Run("two_roots_one_completed", func(t *testing.T) {
		a := createTodo(t, "A", "", date)
		createTodo(t, "B", "", date)
		_, _ = ToggleTodoStatus(date, a.ID)
		c, p, err := GetStats(date)
		if err != nil {
			t.Fatalf("GetStats: %v", err)
		}
		if c != 1 || p != 1 {
			t.Errorf("GetStats = %d completed, %d pending; want 1, 1", c, p)
		}
	})

	t.Run("root_and_child_excludes_child", func(t *testing.T) {
		setupTestDB(t)
		date := "2024-01-03"
		root := createTodo(t, "R", "", date)
		child := createTodo(t, "C", "", date)
		child.ParentID = &root.ID
		_ = UpdateTodo(child)
		_, _ = ToggleTodoStatus(date, root.ID)
		// Root completed, child pending. GetStats only counts roots → 1 completed, 0 pending.
		c, p, err := GetStats(date)
		if err != nil {
			t.Fatalf("GetStats: %v", err)
		}
		if c != 1 || p != 0 {
			t.Errorf("GetStats(root+child) = %d completed, %d pending; want 1, 0", c, p)
		}
	})
}

func TestGetTotalCount(t *testing.T) {
	setupTestDB(t)
	date := "2024-01-01"

	c, err := GetTotalCount(date)
	if err != nil {
		t.Fatalf("GetTotalCount: %v", err)
	}
	if c != 0 {
		t.Errorf("GetTotalCount(empty) = %d, want 0", c)
	}

	r := createTodo(t, "R", "", date)
	child := createTodo(t, "C", "", date)
	child.ParentID = &r.ID
	_ = UpdateTodo(child)

	c, err = GetTotalCount(date)
	if err != nil {
		t.Fatalf("GetTotalCount: %v", err)
	}
	if c != 2 {
		t.Errorf("GetTotalCount(root+child) = %d, want 2", c)
	}
}

func TestHasChildren(t *testing.T) {
	setupTestDB(t)
	date := "2024-01-01"
	a := createTodo(t, "A", "", date)
	b := createTodo(t, "B", "", date)

	ok, err := HasChildren(date, a.ID)
	if err != nil {
		t.Fatalf("HasChildren: %v", err)
	}
	if ok {
		t.Error("HasChildren(A) = true, want false (no children yet)")
	}

	b.ParentID = &a.ID
	_ = UpdateTodo(b)

	ok, err = HasChildren(date, a.ID)
	if err != nil {
		t.Fatalf("HasChildren: %v", err)
	}
	if !ok {
		t.Error("HasChildren(A) = false, want true (B is child)")
	}
}

func TestGetRootTodosByDateRange(t *testing.T) {
	setupTestDB(t)

	createTodo(t, "x", "", "2024-01-01")
	createTodo(t, "y", "", "2024-01-02")
	createTodo(t, "z", "", "2024-01-05")

	todos, err := GetRootTodosByDateRange("2024-01-02", "2024-01-04")
	if err != nil {
		t.Fatalf("GetRootTodosByDateRange: %v", err)
	}
	if len(todos) != 1 {
		t.Fatalf("len(todos) = %d, want 1", len(todos))
	}
	if todos[0].Date != "2024-01-02" || todos[0].Title != "y" {
		t.Errorf("got %+v", todos[0])
	}

	// Ordering: date ASC, is_completed ASC, id ASC
	a := createTodo(t, "a", "", "2024-01-10")
	createTodo(t, "b", "", "2024-01-10")
	_, _ = ToggleTodoStatus("2024-01-10", a.ID)
	all, _ := GetRootTodosByDateRange("2024-01-10", "2024-01-10")
	if len(all) != 2 {
		t.Fatalf("len = %d, want 2", len(all))
	}
	// Pending first (b=id 5), then completed (a=id 4)
	if all[0].Title != "b" || all[1].Title != "a" {
		t.Errorf("order: got %q, %q; want pending first", all[0].Title, all[1].Title)
	}
}

func TestGetTodosByDate(t *testing.T) {
	setupTestDB(t)
	date := "2024-01-01"

	t.Run("ordering", func(t *testing.T) {
		a := createTodo(t, "A", "", date)
		createTodo(t, "B", "", date)
		createTodo(t, "C", "", date)
		_, _ = ToggleTodoStatus(date, a.ID)
		todos, err := GetTodosByDate(date)
		if err != nil {
			t.Fatalf("GetTodosByDate: %v", err)
		}
		if len(todos) != 3 {
			t.Fatalf("len = %d, want 3", len(todos))
		}
		// Pending first (B,C), then completed (A). Among same status: Position then ID.
		if todos[0].Title != "B" || todos[1].Title != "C" || todos[2].Title != "A" {
			t.Errorf("order: got %q, %q, %q; want B, C, A", todos[0].Title, todos[1].Title, todos[2].Title)
		}
	})

	t.Run("tree", func(t *testing.T) {
		setupTestDB(t)
		date := "2024-01-15"
		root := createTodo(t, "R", "", date)
		c1 := createTodo(t, "C1", "", date)
		c2 := createTodo(t, "C2", "", date)
		c1.ParentID = &root.ID
		c2.ParentID = &root.ID
		_ = UpdateTodo(c1)
		_ = UpdateTodo(c2)

		todos, err := GetTodosByDate(date)
		if err != nil {
			t.Fatalf("GetTodosByDate: %v", err)
		}
		if len(todos) != 3 {
			t.Fatalf("len = %d, want 3 (root, c1, c2)", len(todos))
		}
		if todos[0].Title != "R" {
			t.Errorf("todos[0] = %q, want R", todos[0].Title)
		}
		// Children order by sortTodos (Position, ID)
		if todos[1].Title != "C1" || todos[2].Title != "C2" {
			t.Errorf("children order: got %q, %q", todos[1].Title, todos[2].Title)
		}
	})
}

func TestMoveTodo(t *testing.T) {
	setupTestDB(t)
	date := "2024-01-01"

	t.Run("move_down", func(t *testing.T) {
		a := createTodo(t, "A", "", date)
		createTodo(t, "B", "", date)
		createTodo(t, "C", "", date)
		if err := MoveTodo(date, a.ID, 1); err != nil {
			t.Fatalf("MoveTodo(down): %v", err)
		}
		todos, _ := GetTodosByDate(date)
		if len(todos) != 3 {
			t.Fatalf("len = %d, want 3", len(todos))
		}
		if todos[0].Title != "B" || todos[1].Title != "A" || todos[2].Title != "C" {
			t.Errorf("after move A down: got %q, %q, %q; want B, A, C", todos[0].Title, todos[1].Title, todos[2].Title)
		}
	})

	t.Run("move_up", func(t *testing.T) {
		setupTestDB(t)
		createTodo(t, "A", "", date)
		createTodo(t, "B", "", date)
		c := createTodo(t, "C", "", date)
		if err := MoveTodo(date, c.ID, -1); err != nil {
			t.Fatalf("MoveTodo(up): %v", err)
		}
		todos, _ := GetTodosByDate(date)
		if todos[0].Title != "A" || todos[1].Title != "C" || todos[2].Title != "B" {
			t.Errorf("after move C up: got %q, %q, %q; want A, C, B", todos[0].Title, todos[1].Title, todos[2].Title)
		}
	})

	t.Run("bounds_no_op", func(t *testing.T) {
		setupTestDB(t)
		first := createTodo(t, "A", "", date)
		createTodo(t, "B", "", date)
		last := createTodo(t, "C", "", date)
		if err := MoveTodo(date, first.ID, -1); err != nil {
			t.Fatalf("MoveTodo(first up): %v", err)
		}
		todos, _ := GetTodosByDate(date)
		if todos[0].Title != "A" {
			t.Errorf("move first up should be no-op: got %q first", todos[0].Title)
		}
		if err := MoveTodo(date, last.ID, 1); err != nil {
			t.Fatalf("MoveTodo(last down): %v", err)
		}
		todos, _ = GetTodosByDate(date)
		if todos[2].Title != "C" {
			t.Errorf("move last down should be no-op: got %q last", todos[2].Title)
		}
	})

	t.Run("cross_completion_no_op", func(t *testing.T) {
		setupTestDB(t)
		a := createTodo(t, "A", "", date)
		b := createTodo(t, "B", "", date)
		_, _ = ToggleTodoStatus(date, b.ID)
		// A pending, B completed. Move A down would swap with B -> not allowed (different completion).
		if err := MoveTodo(date, a.ID, 1); err != nil {
			t.Fatalf("MoveTodo: %v", err)
		}
		todos, _ := GetTodosByDate(date)
		if todos[0].Title != "A" || todos[1].Title != "B" {
			t.Errorf("cross-completion move should be no-op: got %q, %q; want A, B", todos[0].Title, todos[1].Title)
		}
	})

	t.Run("equal_positions_renormalize", func(t *testing.T) {
		setupTestDB(t)
		a := createTodo(t, "A", "", date)
		b := createTodo(t, "B", "", date)
		a.Position = 0
		b.Position = 0
		_ = UpdateTodo(a)
		_ = UpdateTodo(b)
		if err := MoveTodo(date, b.ID, -1); err != nil {
			t.Fatalf("MoveTodo: %v", err)
		}
		todos, _ := GetTodosByDate(date)
		if todos[0].Title != "B" || todos[1].Title != "A" {
			t.Errorf("equal positions move second up: got %q, %q; want B, A", todos[0].Title, todos[1].Title)
		}
	})
}
