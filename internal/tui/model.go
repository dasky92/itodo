package tui

import (
	"itodo/internal/model"
	"itodo/internal/service"

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

type ViewType int

const (
	DailyView ViewType = iota
	WeeklyView
)

type WeeklyItemType int

const (
	WeeklyHeader WeeklyItemType = iota
	WeeklyTask
)

type WeeklyItem struct {
	Type WeeklyItemType
	Date string
	Todo model.Todo
}

type Model struct {
	svc          *service.TodoService
	todos        []model.Todo
	cursor       int
	selectedDate string
	mode         Mode
	viewType     ViewType

	// Weekly View State
	weeklyTodos     map[string][]model.Todo
	weeklyViewItems []WeeklyItem
	weeklyStart     string
	weeklyEnd       string
	weeklyExpanded  map[string]bool // New: Track expanded/collapsed state

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
		svc:            svc,
		selectedDate:   svc.GetCurrentDate(),
		mode:           Normal,
		viewType:       DailyView,
		weeklyTodos:    make(map[string][]model.Todo),
		weeklyExpanded: make(map[string]bool),
		titleInput:     ti,
		descInput:      ta,
		focusIndex:     0,
		help:           h,
		keys:           Keys,
		logs:           []string{},
	}

	// Initialize default theme
	InitStyles(Themes["Monokai"])

	m.refreshData()
	return m
}

func (m *Model) refreshData() {
	var err error

	switch m.viewType {
	case DailyView:
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
	case WeeklyView:
		start, end, err := m.svc.GetWeekRange(m.selectedDate)
		if err != nil {
			m.logs = append(m.logs, "Error calculating week range: "+err.Error())
			return
		}
		m.weeklyStart = start
		m.weeklyEnd = end

		todos, err := m.svc.GetTasksByRange(start, end)
		if err != nil {
			m.logs = append(m.logs, "Error fetching weekly todos: "+err.Error())
		}

		// Group by date
		m.weeklyTodos = make(map[string][]model.Todo)
		m.completed = 0
		m.pending = 0
		for _, t := range todos {
			m.weeklyTodos[t.Date] = append(m.weeklyTodos[t.Date], t)
			if t.IsCompleted {
				m.completed++
			} else {
				m.pending++
			}
		}

		m.buildWeeklyItems()
	}
}

func (m *Model) buildWeeklyItems() {
	m.weeklyViewItems = []WeeklyItem{}
	current, _ := time.Parse("2006-01-02", m.weeklyStart)
	endT, _ := time.Parse("2006-01-02", m.weeklyEnd)

	for !current.After(endT) {
		d := current.Format("2006-01-02")
		// Add Header
		m.weeklyViewItems = append(m.weeklyViewItems, WeeklyItem{
			Type: WeeklyHeader,
			Date: d,
		})

		// Add Tasks if Expanded
		if m.weeklyExpanded[d] {
			if tasks, ok := m.weeklyTodos[d]; ok {
				for _, t := range tasks {
					m.weeklyViewItems = append(m.weeklyViewItems, WeeklyItem{
						Type: WeeklyTask,
						Todo: t,
					})
				}
			}
		}
		current = current.AddDate(0, 0, 1)
	}

	// Reset cursor for weekly view if out of bounds
	if m.cursor >= len(m.weeklyViewItems) && len(m.weeklyViewItems) > 0 {
		m.cursor = len(m.weeklyViewItems) - 1
	} else if len(m.weeklyViewItems) == 0 {
		m.cursor = 0
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}
