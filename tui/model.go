package tui

import (
	"itodo/model"
	"itodo/service"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Mode int

const (
	Normal Mode = iota
	Adding
	Editing
	Helping
)

type Model struct {
	svc          *service.TodoService
	todos        []model.Todo
	cursor       int
	selectedDate string
	mode         Mode
	textInput    textinput.Model
	help         help.Model
	keys         KeyMap
	logs         []string
	completed    int
	pending      int
	width        int
	height       int
	err          error
	editID       uint // ID of the todo being edited
}

func NewModel(svc *service.TodoService) Model {
	ti := textinput.New()
	ti.Placeholder = "Enter todo title..."
	ti.CharLimit = 156
	ti.Width = 30

	h := help.New()
	h.ShowAll = false // Default to short help

	m := Model{
		svc:          svc,
		selectedDate: svc.GetCurrentDate(),
		mode:         Normal,
		textInput:    ti,
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
	return nil
}
