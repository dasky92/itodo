package tui

import (
	"fmt"
	"strings"
)

func (m Model) View() string {
	var s strings.Builder

	// Header
	header := titleStyle.Render(m.selectedDate)
	stats := statsStyle.Render(fmt.Sprintf("✔ %d  •  ☐ %d", m.completed, m.pending))
	s.WriteString(header + stats + "\n\n")

	// List
	if len(m.todos) == 0 {
		s.WriteString(statusStyle.Render("No todos for this day."))
	} else {
		for i, todo := range m.todos {
			cursor := "  "
			if m.cursor == i {
				cursor = "> "
			}

			line := todo.Title
			if todo.IsCompleted {
				line = fmt.Sprintf("[x] %s", line)
				if m.cursor == i {
					s.WriteString(selectedItemStyle.Render(cursor + line))
				} else {
					s.WriteString(completedItemStyle.Render(cursor + line))
				}
			} else {
				line = fmt.Sprintf("[ ] %s", line)
				if m.cursor == i {
					s.WriteString(selectedItemStyle.Render(cursor + line))
				} else {
					s.WriteString(normalItemStyle.Render(cursor + line))
				}
			}
			s.WriteString("\n")
		}
	}

	s.WriteString("\n\n")

	// Input or Logs
	if m.mode == Adding || m.mode == Editing {
		title := "New Todo:"
		if m.mode == Editing {
			title = "Edit Todo:"
		}
		s.WriteString(fmt.Sprintf("%s\n%s\n(esc to cancel, enter to save)", title, m.textInput.View()))
	} else {
		// Logs (show last 3)
		logCount := len(m.logs)
		start := 0
		if logCount > 3 {
			start = logCount - 3
		}
		for i := start; i < logCount; i++ {
			s.WriteString(logStyle.Render(m.logs[i]) + "\n")
		}
		
		s.WriteString("\n")
		s.WriteString(statusStyle.Render("h/l: date • j/k: nav • n: new • e: edit • d: del • space: toggle • q: quit"))
	}

	return appStyle.Render(s.String())
}
