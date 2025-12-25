package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	// Calculate width for divider
	// We subtract some padding to be safe, though appStyle handles outer padding
	// Default to 40 if width is not set yet
	width := m.width - 10
	if width < 0 {
		width = 40
	}
	divider := dividerStyle.Render(strings.Repeat("┈", width))

	// 1. Header: Date + Stats (Horizontal Layout)
	// Format: [ 📅 YYYY-MM-DD Weekday ] ....... [ ✔ N  ☐ M ]

	formattedDate := m.svc.FormatDateForDisplay(m.selectedDate)
	dateStr := fmt.Sprintf("📅 %s", formattedDate)
	statsStr := fmt.Sprintf("✔ %d  ☐ %d", m.completed, m.pending)

	// Calculate spacing
	dateWidth := lipgloss.Width(dateStr) + 2   // +2 for padding
	statsWidth := lipgloss.Width(statsStr) + 2 // +2 for padding

	gapWidth := width - dateWidth - statsWidth
	if gapWidth < 1 {
		gapWidth = 1
	}
	gap := strings.Repeat(" ", gapWidth)

	headerView := lipgloss.JoinHorizontal(lipgloss.Top,
		dateStyle.Render(dateStr),
		gap,
		statsStyle.Render(statsStr),
	)

	// 2. List
	var listView string
	if len(m.todos) == 0 {
		listView = statusStyle.Render("No todos for this day.")
	} else {
		var rows []string
		for i, todo := range m.todos {
			line := todo.Title

			// Icon selection
			icon := "○"
			if todo.IsCompleted {
				icon = "✔"
			}

			// Render item
			str := fmt.Sprintf("%s %s", icon, line)

			var renderedRow string
			if m.cursor == i {
				// Selected
				if todo.IsCompleted {
					renderedRow = selectedItemStyle.Copy().Strikethrough(true).Render(str)
				} else {
					renderedRow = selectedItemStyle.Render(str)
				}
			} else {
				// Not Selected
				if todo.IsCompleted {
					renderedRow = completedItemStyle.Render(str)
				} else {
					renderedRow = normalItemStyle.Render(str)
				}
			}
			rows = append(rows, renderedRow)
		}
		listView = lipgloss.JoinVertical(lipgloss.Left, rows...)
	}

	// 3. Footer: Logs or Input
	var footerView string
	if m.mode == Adding || m.mode == Editing {
		title := "New Todo:"
		if m.mode == Editing {
			title = "Edit Todo:"
		}
		footerView = fmt.Sprintf("%s\n%s\n(esc to cancel, enter to save)", title, m.textInput.View())
	} else if m.mode == Helping {
		footerView = m.help.View(m.keys)
	} else {
		// Logs (show last 3, faded)
		logCount := len(m.logs)
		start := 0
		if logCount > 3 {
			start = logCount - 3
		}
		var logRows []string
		for i := start; i < logCount; i++ {
			logRows = append(logRows, logStyle.Render(m.logs[i]))
		}
		if len(logRows) > 0 {
			footerView = lipgloss.JoinVertical(lipgloss.Left, logRows...)
		} else {
			footerView = logStyle.Render("Ready.")
		}
	}

	// Combine all sections with dividers
	// Layout:
	// Header (Date ... Stats)
	// ----
	// List
	// ----
	// Footer

	return appStyle.Render(
		lipgloss.JoinVertical(lipgloss.Left,
			headerView,
			divider,
			listView,
			divider,
			footerView,
		),
	)
}
