package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	// 1. Form View (Adding/Editing) - Full Screen
	if m.mode == Adding || m.mode == Editing {
		title := "Create New Task"
		if m.mode == Editing {
			title = "Edit Task"
		}

		// Ensure width is safe
		width := m.width - 10
		if width < 0 {
			width = 40
		}
		divider := dividerStyle.Render(strings.Repeat("┈", width))

		return appStyle.Render(lipgloss.JoinVertical(lipgloss.Left,
			formTitleStyle.Render(title),

			labelStyle.Render("Title"),
			m.titleInput.View(),

			labelStyle.Render("Description"),
			m.descInput.View(),

			divider,
			statusStyle.Render("Ctrl+S: Save • Esc: Cancel • Tab: Switch Field • Enter: Next/Newline"),
		))
	}

	// 2. Normal View (List)

	// Calculate width for divider
	width := m.width - 10
	if width < 0 {
		width = 40
	}
	divider := dividerStyle.Render(strings.Repeat("┈", width))

	// Header: Date + Title + Progress (3-Column Layout)
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

	statsStr := fmt.Sprintf("%s%s%s %d/%d",
		statusStyle.Render("["),
		bar,
		statusStyle.Render("]"),
		m.completed,
		total,
	)
	rightView := statsStyle.Render(statsStr)

	// Calculate Spacing for Absolute Centering
	leftWidth := lipgloss.Width(leftView)
	centerWidth := lipgloss.Width(centerView)
	rightWidth := lipgloss.Width(rightView)

	// Goal: Place centerView at the exact center of the screen
	// Center of screen is at m.width / 2
	// centerView starts at (m.width - centerWidth) / 2

	targetCenterStart := (width - centerWidth) / 2

	gapLWidth := targetCenterStart - leftWidth
	if gapLWidth < 0 {
		gapLWidth = 0 // Not enough space to center absolutely
	}

	// Calculate right gap
	// Right starts at targetCenterStart + centerWidth + gapRWidth
	// We want rightView to end at width
	// So gapRWidth = width - rightWidth - (targetCenterStart + centerWidth)

	gapRWidth := width - rightWidth - (targetCenterStart + centerWidth)
	if gapRWidth < 0 {
		gapRWidth = 0
	}

	// Fallback to equal spacing if absolute centering is impossible (screen too small)
	if gapLWidth == 0 || gapRWidth == 0 {
		availWidth := width - leftWidth - centerWidth - rightWidth
		if availWidth > 0 {
			gapLWidth = availWidth / 2
			gapRWidth = availWidth - gapLWidth
		} else {
			gapLWidth = 0
			gapRWidth = 0
		}
	}

	gapL := strings.Repeat(" ", gapLWidth)
	gapR := strings.Repeat(" ", gapRWidth)

	headerView := lipgloss.JoinHorizontal(lipgloss.Top,
		leftView,
		gapL,
		centerView,
		gapR,
		rightView,
	)

	// Footer: Logs or Help
	var footerView string
	if m.mode == Helping {
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

	// List
	var listView string

	// Calculate available height for the list
	// Fixed Elements:
	// 1. App Padding: 2 lines (1 top, 1 bottom) - handled by appStyle
	// 2. Header: 1 line (assumed)
	// 3. Divider 1: 3 lines (Padding 1,0 -> 1 top, 1 bottom + 1 content)
	// 4. Divider 2: 3 lines
	// 5. Footer: Dynamic Height

	headerHeight := lipgloss.Height(headerView)
	dividerHeight := lipgloss.Height(divider) // Should be 3
	footerHeight := lipgloss.Height(footerView)
	appPaddingHeight := 2 // appStyle.Padding(1, 4)

	fixedHeight := appPaddingHeight + headerHeight + dividerHeight*2 + footerHeight

	targetListHeight := m.height - fixedHeight
	if targetListHeight < 0 {
		targetListHeight = 0
	}

	if len(m.todos) == 0 {
		listView = statusStyle.Render("No todos for this day.")
		// Fill remaining height
		currentH := lipgloss.Height(listView)
		if targetListHeight > currentH {
			listView += strings.Repeat("\n", targetListHeight-currentH)
		}
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

		// Join rows
		content := lipgloss.JoinVertical(lipgloss.Left, rows...)

		// Ensure fixed height using lipgloss.Place or manual padding
		currentHeight := lipgloss.Height(content)
		if currentHeight < targetListHeight {
			diff := targetListHeight - currentHeight
			content += strings.Repeat("\n", diff)
		}

		listView = content
	}

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
