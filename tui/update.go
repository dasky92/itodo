package tui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.help.Width = msg.Width

	case tea.KeyMsg:
		switch m.mode {
		case Normal:
			switch {
			case key.Matches(msg, m.keys.Quit):
				return m, tea.Quit

			case key.Matches(msg, m.keys.Up):
				if m.cursor > 0 {
					m.cursor--
				}

			case key.Matches(msg, m.keys.Down):
				if m.cursor < len(m.todos)-1 {
					m.cursor++
				}

			case key.Matches(msg, m.keys.Left):
				m.selectedDate = m.svc.GetPrevDate(m.selectedDate)
				m.refreshData()

			case key.Matches(msg, m.keys.Right):
				m.selectedDate = m.svc.GetNextDate(m.selectedDate)
				m.refreshData()

			case key.Matches(msg, m.keys.New):
				m.mode = Adding
				m.inputs[0].SetValue("")
				m.inputs[1].SetValue("")
				m.inputs[0].Focus()
				m.focusIndex = 0
				return m, textinput.Blink

			case key.Matches(msg, m.keys.Edit):
				if len(m.todos) > 0 {
					m.mode = Editing
					todo := m.todos[m.cursor]
					m.editID = todo.ID
					m.inputs[0].SetValue(todo.Title)
					m.inputs[1].SetValue(todo.Description)
					m.inputs[0].Focus()
					m.focusIndex = 0
					return m, textinput.Blink
				}

			case key.Matches(msg, m.keys.Delete):
				if len(m.todos) > 0 {
					todo := m.todos[m.cursor]
					log, err := m.svc.Delete(todo.ID)
					if err != nil {
						m.logs = append(m.logs, "Error: "+err.Error())
					} else {
						m.logs = append(m.logs, log)
						m.refreshData()
					}
				}

			case key.Matches(msg, m.keys.Toggle):
				if len(m.todos) > 0 {
					todo := m.todos[m.cursor]
					log, err := m.svc.Toggle(todo.ID)
					if err != nil {
						m.logs = append(m.logs, "Error: "+err.Error())
					} else {
						m.logs = append(m.logs, log)
						m.refreshData()
					}
				}

			case key.Matches(msg, m.keys.Indent):
				if len(m.todos) > 0 && m.cursor > 0 {
					todo := m.todos[m.cursor]
					prevTodo := m.todos[m.cursor-1]
					log, err := m.svc.IndentTodo(todo.ID, prevTodo.ID)
					if err != nil {
						m.logs = append(m.logs, "Error: "+err.Error())
					} else {
						m.logs = append(m.logs, log)
						m.refreshData()
					}
				}

			case key.Matches(msg, m.keys.Outdent):
				if len(m.todos) > 0 {
					todo := m.todos[m.cursor]
					log, err := m.svc.OutdentTodo(todo.ID)
					if err != nil {
						m.logs = append(m.logs, "Error: "+err.Error())
					} else {
						m.logs = append(m.logs, log)
						m.refreshData()
					}
				}

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

		case Adding:
			switch msg.String() {
			case "enter":
				title := m.inputs[0].Value()
				desc := m.inputs[1].Value()
				if title != "" {
					log, err := m.svc.Add(title, desc, m.selectedDate)
					if err != nil {
						m.logs = append(m.logs, "Error: "+err.Error())
					} else {
						m.logs = append(m.logs, log)
						m.refreshData()
						m.mode = Normal
						m.inputs[0].Blur()
						m.inputs[1].Blur()
					}
				}
			case "esc":
				m.mode = Normal
				m.inputs[0].Blur()
				m.inputs[1].Blur()
			case "tab":
				m.focusIndex = (m.focusIndex + 1) % len(m.inputs)
				if m.focusIndex == 0 {
					m.inputs[1].Blur()
					m.inputs[0].Focus()
				} else {
					m.inputs[0].Blur()
					m.inputs[1].Focus()
				}
			default:
				m.inputs[m.focusIndex], cmd = m.inputs[m.focusIndex].Update(msg)
			}

		case Editing:
			switch msg.String() {
			case "enter":
				title := m.inputs[0].Value()
				desc := m.inputs[1].Value()
				if title != "" {
					log, err := m.svc.Edit(m.editID, title, desc)
					if err != nil {
						m.logs = append(m.logs, "Error: "+err.Error())
					} else {
						m.logs = append(m.logs, log)
						m.refreshData()
						m.mode = Normal
						m.inputs[0].Blur()
						m.inputs[1].Blur()
					}
				}
			case "esc":
				m.mode = Normal
				m.inputs[0].Blur()
				m.inputs[1].Blur()
			case "tab":
				m.focusIndex = (m.focusIndex + 1) % len(m.inputs)
				if m.focusIndex == 0 {
					m.inputs[1].Blur()
					m.inputs[0].Focus()
				} else {
					m.inputs[0].Blur()
					m.inputs[1].Focus()
				}
			default:
				m.inputs[m.focusIndex], cmd = m.inputs[m.focusIndex].Update(msg)
			}
		}
	}

	return m, cmd
}
