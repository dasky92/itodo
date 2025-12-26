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

	// 1. Header: Date + Title + Progress (3-Column Layout)
	// Left: Date
	formattedDate := m.svc.FormatDateForDisplay(m.selectedDate)
	dateStr := fmt.Sprintf("📅 %s", formattedDate)
	leftView := dateStyle.Render(dateStr)

	// Center: Title
	titleStr := "Daily Tasks"
	centerView := titleStyle.Render(titleStr)

	// Right: Progress Bar + Total Count
	total := m.completed + m.pending
	percent := 0.0
	if total > 0 {
		percent = float64(m.completed) / float64(total)
	}

	// Progress Bar Rendering
	barWidth := 10
	filled := int(percent * float64(barWidth))
	empty := barWidth - filled

	bar := progressStyle.Render(strings.Repeat(progressFullChar, filled)) +
		progressEmptyStyle.Render(strings.Repeat(progressEmptyChar, empty))

	statsStr := fmt.Sprintf("%s %d/%d", bar, m.completed, total)
	rightView := statsStyle.Render(statsStr)

	// Calculate Spacing
	// We want: [Left] ... [Center] ... [Right]
	leftWidth := lipgloss.Width(leftView)
	centerWidth := lipgloss.Width(centerView)
	rightWidth := lipgloss.Width(rightView)

	// Calculate gaps
	// Available width for gaps
	availWidth := width - leftWidth - centerWidth - rightWidth
	if availWidth < 2 {
		availWidth = 2
	}

	gapLWidth := availWidth / 2
	gapRWidth := availWidth - gapLWidth // Handle odd numbers

	gapL := strings.Repeat(" ", gapLWidth)
	gapR := strings.Repeat(" ", gapRWidth)

	headerView := lipgloss.JoinHorizontal(lipgloss.Top,
		leftView,
		gapL,
		centerView,
		gapR,
		rightView,
	)

	// 2. List
	var listView string
	if len(m.todos) == 0 {
		listView = statusStyle.Render("No todos for this day.")
	} else {
		var rows []string
		for i, todo := range m.todos {
			line := todo.Title

			// Indent child tasks
			prefix := ""
			if todo.ParentID != nil {
				prefix = "    "
			}

			// Icon selection
			icon := "○"
			if todo.IsCompleted {
				icon = "✔"
			}

			// Render item
			str := fmt.Sprintf("%s%s %s", prefix, icon, line)

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

		// Input View
		inputView := lipgloss.JoinVertical(lipgloss.Left,
			fmt.Sprintf("Title: %s", m.inputs[0].View()),
			fmt.Sprintf("Desc : %s", m.inputs[1].View()),
		)

		footerView = fmt.Sprintf("%s\n%s\n(tab to switch, esc to cancel, enter to save)", title, inputView)
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
