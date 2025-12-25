package tui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch m.mode {
		case Normal:
			switch msg.String() {
			case "q", "ctrl+c":
				return m, tea.Quit

			case "k", "up":
				if m.cursor > 0 {
					m.cursor--
				}

			case "j", "down":
				if m.cursor < len(m.todos)-1 {
					m.cursor++
				}

			case "h", "left":
				m.selectedDate = m.svc.GetPrevDate(m.selectedDate)
				m.refreshData()

			case "l", "right":
				m.selectedDate = m.svc.GetNextDate(m.selectedDate)
				m.refreshData()

			case "n":
				m.mode = Adding
				m.textInput.SetValue("")
				m.textInput.Focus()
				return m, textinput.Blink

			case "e":
				if len(m.todos) > 0 {
					m.mode = Editing
					todo := m.todos[m.cursor]
					m.editID = todo.ID
					m.textInput.SetValue(todo.Title)
					m.textInput.Focus()
					return m, textinput.Blink
				}

			case "d":
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

			case " ":
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
			}

		case Adding:
			switch msg.String() {
			case "enter":
				text := m.textInput.Value()
				if text != "" {
					log, err := m.svc.Add(text, m.selectedDate)
					if err != nil {
						m.logs = append(m.logs, "Error: "+err.Error())
					} else {
						m.logs = append(m.logs, log)
						m.refreshData()
						m.mode = Normal
						m.textInput.Blur()
					}
				}
			case "esc":
				m.mode = Normal
				m.textInput.Blur()
			default:
				m.textInput, cmd = m.textInput.Update(msg)
			}

		case Editing:
			switch msg.String() {
			case "enter":
				text := m.textInput.Value()
				if text != "" {
					log, err := m.svc.Edit(m.editID, text)
					if err != nil {
						m.logs = append(m.logs, "Error: "+err.Error())
					} else {
						m.logs = append(m.logs, log)
						m.refreshData()
						m.mode = Normal
						m.textInput.Blur()
					}
				}
			case "esc":
				m.mode = Normal
				m.textInput.Blur()
			default:
				m.textInput, cmd = m.textInput.Update(msg)
			}
		}
	}

	return m, cmd
}
