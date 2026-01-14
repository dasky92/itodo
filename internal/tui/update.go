package tui

import (
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.help.Width = msg.Width
		m.descInput.SetWidth(msg.Width - 10) // Adjust textarea width

	case tea.KeyMsg:
		switch m.mode {
		case CalendarMode:
			switch {
			case key.Matches(msg, m.keys.Cancel): // Esc
				m.mode = Normal

			case key.Matches(msg, m.keys.Toggle): // Enter (Select)
				m.selectedDate = m.calendarCursor.Format("2006-01-02")
				m.mode = Normal
				m.refreshData()

			case key.Matches(msg, m.keys.Today): // Space (Go to Today)
				now := time.Now()
				m.calendarCursor = now
				m.calendarViewDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
				m.refreshCalendarMarks()

			case key.Matches(msg, m.keys.Left):
				m.calendarCursor = m.calendarCursor.AddDate(0, 0, -1)
				if m.calendarCursor.Month() != m.calendarViewDate.Month() || m.calendarCursor.Year() != m.calendarViewDate.Year() {
					m.calendarViewDate = time.Date(m.calendarCursor.Year(), m.calendarCursor.Month(), 1, 0, 0, 0, 0, m.calendarCursor.Location())
					m.refreshCalendarMarks()
				}

			case key.Matches(msg, m.keys.Right):
				m.calendarCursor = m.calendarCursor.AddDate(0, 0, 1)
				if m.calendarCursor.Month() != m.calendarViewDate.Month() || m.calendarCursor.Year() != m.calendarViewDate.Year() {
					m.calendarViewDate = time.Date(m.calendarCursor.Year(), m.calendarCursor.Month(), 1, 0, 0, 0, 0, m.calendarCursor.Location())
					m.refreshCalendarMarks()
				}

			case key.Matches(msg, m.keys.Up):
				m.calendarCursor = m.calendarCursor.AddDate(0, 0, -7)
				if m.calendarCursor.Month() != m.calendarViewDate.Month() || m.calendarCursor.Year() != m.calendarViewDate.Year() {
					m.calendarViewDate = time.Date(m.calendarCursor.Year(), m.calendarCursor.Month(), 1, 0, 0, 0, 0, m.calendarCursor.Location())
					m.refreshCalendarMarks()
				}

			case key.Matches(msg, m.keys.Down):
				m.calendarCursor = m.calendarCursor.AddDate(0, 0, 7)
				if m.calendarCursor.Month() != m.calendarViewDate.Month() || m.calendarCursor.Year() != m.calendarViewDate.Year() {
					m.calendarViewDate = time.Date(m.calendarCursor.Year(), m.calendarCursor.Month(), 1, 0, 0, 0, 0, m.calendarCursor.Location())
					m.refreshCalendarMarks()
				}
			}

		case Normal:
			switch {
			case key.Matches(msg, m.keys.Quit):
				return m, tea.Quit

			case key.Matches(msg, m.keys.PrevView):
				m.viewType = DailyView
				m.cursor = 0
				m.refreshData()

			case key.Matches(msg, m.keys.NextView):
				m.viewType = WeeklyView
				m.cursor = 0
				m.refreshData()

			case key.Matches(msg, m.keys.Up):
				if m.cursor > 0 {
					m.cursor--
				}

			case key.Matches(msg, m.keys.Down):
				maxLen := len(m.todos)
				if m.viewType == WeeklyView {
					maxLen = len(m.weeklyViewItems)
				}
				if m.cursor < maxLen-1 {
					m.cursor++
				}

			case key.Matches(msg, m.keys.Left):
				if m.viewType == DailyView {
					m.selectedDate = m.svc.GetPrevDate(m.selectedDate)
				} else {
					t, _ := time.Parse("2006-01-02", m.selectedDate)
					m.selectedDate = t.AddDate(0, 0, -7).Format("2006-01-02")
				}
				m.refreshData()

			case key.Matches(msg, m.keys.Right):
				if m.viewType == DailyView {
					m.selectedDate = m.svc.GetNextDate(m.selectedDate)
				} else {
					t, _ := time.Parse("2006-01-02", m.selectedDate)
					m.selectedDate = t.AddDate(0, 0, 7).Format("2006-01-02")
				}
				m.refreshData()

			case key.Matches(msg, m.keys.New):
				if m.viewType == WeeklyView {
					return m, nil // Read-only
				}
				// Check if selected date is in the past
				today := m.svc.GetCurrentDate()
				if m.selectedDate < today {
					m.logs = append(m.logs, "Error: Cannot create tasks for past dates")
					return m, nil
				}

				m.mode = Adding
				m.titleInput.SetValue("")
				m.descInput.SetValue("")
				m.titleInput.Focus()
				m.descInput.Blur()
				m.focusIndex = 0
				return m, textinput.Blink

			case key.Matches(msg, m.keys.Edit):
				if m.viewType == WeeklyView {
					return m, nil
				}
				if len(m.todos) > 0 {
					m.mode = Editing
					todo := m.todos[m.cursor]
					m.editID = todo.ID
					m.titleInput.SetValue(todo.Title)
					m.descInput.SetValue(todo.Description)
					m.titleInput.Focus()
					m.descInput.Blur()
					m.focusIndex = 0
					return m, textinput.Blink
				}

			case key.Matches(msg, m.keys.Delete):
				if m.viewType == WeeklyView {
					return m, nil
				}
				if len(m.todos) > 0 {
					todo := m.todos[m.cursor]
					log, err := m.svc.Delete(m.selectedDate, todo.ID)
					if err != nil {
						m.logs = append(m.logs, "Error: "+err.Error())
					} else {
						m.logs = append(m.logs, log)
						m.refreshData()
					}
				}

			case key.Matches(msg, m.keys.Toggle):
				if m.viewType == WeeklyView {
					// Toggle Expansion if header
					if len(m.weeklyViewItems) > 0 {
						item := m.weeklyViewItems[m.cursor]
						if item.Type == WeeklyHeader {
							m.weeklyExpanded[item.Date] = !m.weeklyExpanded[item.Date]
							m.buildWeeklyItems()
						}
					}
					return m, nil
				}
				if len(m.todos) > 0 {
					todo := m.todos[m.cursor]
					log, err := m.svc.Toggle(m.selectedDate, todo.ID)
					if err != nil {
						m.logs = append(m.logs, "Error: "+err.Error())
					} else {
						m.logs = append(m.logs, log)
						m.refreshData()
					}
				}

			case key.Matches(msg, m.keys.Indent):
				if m.viewType == WeeklyView {
					return m, nil
				}
				if len(m.todos) > 0 && m.cursor > 0 {
					todo := m.todos[m.cursor]
					prevTodo := m.todos[m.cursor-1]
					log, err := m.svc.IndentTodo(m.selectedDate, todo.ID, prevTodo.ID)
					if err != nil {
						m.logs = append(m.logs, "Error: "+err.Error())
					} else {
						m.logs = append(m.logs, log)
						m.refreshData()
					}
				}

			case key.Matches(msg, m.keys.Outdent):
				if m.viewType == WeeklyView {
					return m, nil
				}
				if len(m.todos) > 0 {
					todo := m.todos[m.cursor]
					log, err := m.svc.OutdentTodo(m.selectedDate, todo.ID)
					if err != nil {
						m.logs = append(m.logs, "Error: "+err.Error())
					} else {
						m.logs = append(m.logs, log)
						m.refreshData()
					}
				}

			case key.Matches(msg, m.keys.MoveUp):
				if m.viewType == WeeklyView {
					return m, nil
				}
				if len(m.todos) > 0 {
					todo := m.todos[m.cursor]
					log, err := m.svc.MoveTodo(m.selectedDate, todo.ID, -1)
					if err != nil {
						m.logs = append(m.logs, "Error: "+err.Error())
					} else {
						m.logs = append(m.logs, log)
						m.refreshData()
						// Cursor follows the item. If successful move up, decrement cursor
						if m.cursor > 0 {
							// Check if visual order actually changed (sometimes sort criteria prevents move)
							// But blindly following direction is usually fine if move succeeded.
							// Wait, if we are at top, we can't move up.
							// The service check will return nil if no move happened.
							// However, simpler to just re-find the ID?
							// Or simpler: just decrement if successful.
							// But we need to know if it ACTUALLY moved.
							// The current implementation of MoveTodo returns "Moved todo up/down" on success.
							// Let's assume if err == nil, it moved.
							// BUT: MoveTodo returns nil err if it hits bounds or different completion status.
							// The log message is "Moved todo up/down".
							// Let's re-find the cursor position by ID to be robust.
							for i, t := range m.todos {
								if t.ID == todo.ID {
									m.cursor = i
									break
								}
							}
						}
					}
				}

			case key.Matches(msg, m.keys.MoveDown):
				if m.viewType == WeeklyView {
					return m, nil
				}
				if len(m.todos) > 0 {
					todo := m.todos[m.cursor]
					log, err := m.svc.MoveTodo(m.selectedDate, todo.ID, 1)
					if err != nil {
						m.logs = append(m.logs, "Error: "+err.Error())
					} else {
						m.logs = append(m.logs, log)
						m.refreshData()
						// Cursor follows the item
						// Re-find the cursor position by ID to be robust.
						for i, t := range m.todos {
							if t.ID == todo.ID {
								m.cursor = i
								break
							}
						}
					}
				}

			case key.Matches(msg, m.keys.Today):
				m.selectedDate = m.svc.GetCurrentDate()
				m.refreshData()

			case key.Matches(msg, m.keys.Calendar):
				m.mode = CalendarMode
				t, err := time.Parse("2006-01-02", m.selectedDate)
				if err != nil {
					t = time.Now()
				}
				m.calendarCursor = t
				m.calendarViewDate = time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
				m.refreshCalendarMarks()

			case key.Matches(msg, m.keys.Help):
				m.mode = Helping
				m.help.ShowAll = true
			}

		case Helping:
			switch {
			case key.Matches(msg, m.keys.Help), msg.String() == "esc", key.Matches(msg, m.keys.Quit):
				m.mode = Normal
				m.help.ShowAll = false
			}

		case Adding, Editing:
			// Form Logic
			switch {
			case key.Matches(msg, m.keys.Cancel): // Esc
				m.mode = Normal
				m.titleInput.Blur()
				m.descInput.Blur()

			case key.Matches(msg, m.keys.Save): // Ctrl+S
				title := m.titleInput.Value()
				desc := m.descInput.Value()

				if title != "" {
					var log string
					var err error

					if m.mode == Adding {
						log, err = m.svc.Add(title, desc, m.selectedDate)
					} else {
						log, err = m.svc.Edit(m.selectedDate, m.editID, title, desc)
					}

					if err != nil {
						m.logs = append(m.logs, "Error: "+err.Error())
					} else {
						m.logs = append(m.logs, log)
						m.refreshData()
						m.mode = Normal
						m.titleInput.Blur()
						m.descInput.Blur()
					}
				}

			case msg.String() == "tab":
				m.focusIndex = (m.focusIndex + 1) % 2
				if m.focusIndex == 0 {
					m.descInput.Blur()
					m.titleInput.Focus()
				} else {
					m.titleInput.Blur()
					m.descInput.Focus()
				}
				return m, nil

			case msg.String() == "enter":
				if m.focusIndex == 0 {
					// In title, Enter moves to Description
					m.focusIndex = 1
					m.titleInput.Blur()
					m.descInput.Focus()
					return m, nil
				}
				// In description, Enter is handled by textarea for newline
			}
		}
	}

	// Update inputs if in Form mode
	if m.mode == Adding || m.mode == Editing {
		if m.focusIndex == 0 {
			m.titleInput, cmd = m.titleInput.Update(msg)
			cmds = append(cmds, cmd)
		} else {
			m.descInput, cmd = m.descInput.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}
