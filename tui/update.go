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
				m.textInput.SetValue("")
				m.textInput.Focus()
				return m, textinput.Blink

			case key.Matches(msg, m.keys.Edit):
				if len(m.todos) > 0 {
					m.mode = Editing
					todo := m.todos[m.cursor]
					m.editID = todo.ID
					m.textInput.SetValue(todo.Title)
					m.textInput.Focus()
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
