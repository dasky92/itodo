package tui

import (
	"itodo/model"
	"itodo/service"

	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Mode int

const (
	Normal Mode = iota
	Adding
	Editing
	Helping
	CalendarMode
)

type Model struct {
	svc          *service.TodoService
	todos        []model.Todo
	cursor       int
	selectedDate string
	mode         Mode

	// Calendar State
	calendarCursor   time.Time
	calendarViewDate time.Time // Month being viewed

	// Form Inputs
	titleInput textinput.Model
	descInput  textarea.Model
	focusIndex int // 0: Title, 1: Description

	help      help.Model
	keys      KeyMap
	logs      []string
	completed int
	pending   int
	width     int
	height    int
	err       error
	editID    uint // ID of the todo being edited
}

func NewModel(svc *service.TodoService) Model {
	// Initialize Title Input
	ti := textinput.New()
	ti.Placeholder = "Task Title"
	ti.Focus()
	ti.CharLimit = 100
	ti.Width = 50

	// Initialize Description Textarea
	ta := textarea.New()
	ta.Placeholder = "Description"
	ta.ShowLineNumbers = false
	ta.SetHeight(10)
	ta.SetWidth(50)

	h := help.New()
	h.ShowAll = false // Default to short help

	m := Model{
		svc:          svc,
		selectedDate: svc.GetCurrentDate(),
		mode:         Normal,
		titleInput:   ti,
		descInput:    ta,
		focusIndex:   0,
		help:         h,
		keys:         Keys,
		logs:         []string{},
	}
	m.refreshData()
	return m
}

func (m *Model) refreshData() {
	var err error
	m.todos, err = m.svc.List(m.selectedDate)
	if err != nil {
		m.logs = append(m.logs, "Error fetching todos: "+err.Error())
	}
	m.completed, m.pending, err = m.svc.GetStats(m.selectedDate)
	if err != nil {
		m.logs = append(m.logs, "Error fetching stats: "+err.Error())
	}
	// Adjust cursor if out of bounds
	if m.cursor >= len(m.todos) && len(m.todos) > 0 {
		m.cursor = len(m.todos) - 1
	} else if len(m.todos) == 0 {
		m.cursor = 0
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}
