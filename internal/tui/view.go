package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	if m.mode == CalendarMode {
		return m.calendarView()
	}

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
	var leftView, centerView string

	if m.viewType == DailyView {
		leftView = dateStyle.Render(m.selectedDate)
		t, _ := time.Parse("2006-01-02", m.selectedDate)
		centerView = titleStyle.Render(t.Weekday().String())
	} else {
		// Weekly View
		t, _ := time.Parse("2006-01-02", m.weeklyStart)
		_, week := t.ISOWeek()
		leftView = dateStyle.Render(fmt.Sprintf("%s - %s", m.weeklyStart, m.weeklyEnd))
		centerView = titleStyle.Render(fmt.Sprintf("Week %d", week))
	}

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

	if m.viewType == DailyView {
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
					prefix = "  "
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
					renderedRow = selectedItemStyle.Render(str)
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
	} else {
		// Weekly View
		if len(m.weeklyViewItems) == 0 {
			listView = statusStyle.Render("No todos for this week.")
			currentH := lipgloss.Height(listView)
			if targetListHeight > currentH {
				listView += strings.Repeat("\n", targetListHeight-currentH)
			}
		} else {
			// Calculate visible range based on cursor and height
			start, end := 0, len(m.weeklyViewItems)
			if len(m.weeklyViewItems) > targetListHeight {
				half := targetListHeight / 2
				if m.cursor < half {
					start = 0
					end = targetListHeight
				} else if m.cursor >= len(m.weeklyViewItems)-half {
					start = len(m.weeklyViewItems) - targetListHeight
					end = len(m.weeklyViewItems)
				} else {
					start = m.cursor - half
					end = start + targetListHeight
				}
			}

			var rows []string
			for i := start; i < end; i++ {
				item := m.weeklyViewItems[i]
				var line string

				if item.Type == WeeklyHeader {
					// Date Header with Progress
					dateTodos := m.weeklyTodos[item.Date]
					c, p := 0, 0
					for _, t := range dateTodos {
						if t.IsCompleted {
							c++
						} else {
							p++
						}
					}
					total := c + p
					percent := 0.0
					if total > 0 {
						percent = float64(c) / float64(total)
					}

					bw := 5
					f := int(percent * float64(bw))
					e := bw - f
					bar := "[" + strings.Repeat("=", f) + strings.Repeat("-", e) + "]"

					// Add expansion indicator
					indicator := "▶"
					if m.weeklyExpanded[item.Date] {
						indicator = "▼"
					}

					headerText := fmt.Sprintf("%s %s %s %d/%d", indicator, item.Date, bar, c, total)
					line = lipgloss.NewStyle().Foreground(CurrentTheme.Pink).Bold(true).Render(headerText)
				} else {
					// Task Item
					icon := "○"
					if item.Todo.IsCompleted {
						icon = "✔"
					}
					// Tree structure guide line
					// Using └─ for visual hierarchy
					taskText := fmt.Sprintf("  └─ %s %s", icon, item.Todo.Title)

					if item.Todo.IsCompleted {
						line = completedItemStyle.Copy().MarginBottom(0).Render(taskText)
					} else {
						line = normalItemStyle.Copy().MarginBottom(0).Render(taskText)
					}
				}

				// Apply Cursor Selection with Fixed Width Prefix to avoid Jitter
				// prefix := "  " // Default 2 spaces
				if m.cursor == i {
					// prefix = "> " // Selected 2 chars
					// Highlight line logic if needed, but prefix is enough for now
					// To fix jitter, ensure prefix is prepended OUTSIDE the styled render or consistently

					// Re-render line with selection style if needed, OR just prepend prefix
					// For simple fix:
					line = selectedItemStyle.Copy().MarginBottom(0).Render(line)
				} else {
					// Unselected
					line = lipgloss.NewStyle().PaddingLeft(1).Render(line) // Match padding of selectedItemStyle
				}

				// Construct final row: Prefix + Line
				// Note: selectedItemStyle has PaddingLeft(1) and Border left.
				// We need to match visual width.

				// Let's simplify:
				// Use a dedicated cursor column.
				cursorStr := "  "
				if m.cursor == i {
					cursorStr = "> "
				}

				// Override line rendering to be cleaner
				if item.Type == WeeklyHeader {
					// Re-render header to be clean
					baseStyle := lipgloss.NewStyle().Foreground(CurrentTheme.Pink).Bold(true)
					if m.cursor == i {
						baseStyle = baseStyle.Background(CurrentTheme.Highlight) // Highlight background
					}
					line = baseStyle.Render(fmt.Sprintf("%s %s", cursorStr, item.Date)) // Simplified for now, let's restore content

					// Restore content logic
					dateTodos := m.weeklyTodos[item.Date]
					c, p := 0, 0
					for _, t := range dateTodos {
						if t.IsCompleted {
							c++
						} else {
							p++
						}
					}
					total := c + p
					percent := 0.0
					if total > 0 {
						percent = float64(c) / float64(total)
					}
					bw := 5
					f := int(percent * float64(bw))
					e := bw - f
					bar := "[" + strings.Repeat("=", f) + strings.Repeat("-", e) + "]"
					indicator := "▶"
					if m.weeklyExpanded[item.Date] {
						indicator = "▼"
					}

					content := fmt.Sprintf("%s %s %s %d/%d", indicator, item.Date, bar, c, total)
					line = baseStyle.Render(cursorStr + content)

				} else {
					// Task Item
					icon := "○"
					if item.Todo.IsCompleted {
						icon = "✔"
					}

					// Indent task more: cursor(2) + indent(4) + tree(3) = 9 chars total prefix?
					// cursorStr is 2 chars.
					// We want task to be indented relative to header.
					// Header: "> ▼ 2023-..."
					// Task:   "      └─ ○ Title"

					taskContent := fmt.Sprintf("      └─ %s %s", icon, item.Todo.Title)

					style := normalItemStyle.Copy().MarginBottom(0).PaddingLeft(0)
					if item.Todo.IsCompleted {
						style = completedItemStyle.Copy().MarginBottom(0).PaddingLeft(0)
					}

					if m.cursor == i {
						style = style.Background(CurrentTheme.Highlight) // Highlight background
					}

					line = style.Render(cursorStr + taskContent)
				}
				rows = append(rows, line)
			}
			content := lipgloss.JoinVertical(lipgloss.Left, rows...)
			// Fill empty space
			currentHeight := lipgloss.Height(content)
			if currentHeight < targetListHeight {
				diff := targetListHeight - currentHeight
				content += strings.Repeat("\n", diff)
			}
			listView = content
		}
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

func (m Model) calendarView() string {
	t := m.calendarViewDate
	month := t.Month().String()
	year := t.Year()
	title := fmt.Sprintf("%s %d", month, year)

	header := monthTitleStyle.Render(title)

	// Weekdays
	weekdays := []string{"Su", "Mo", "Tu", "We", "Th", "Fr", "Sa"}
	var weekdayViews []string
	for _, w := range weekdays {
		weekdayViews = append(weekdayViews, weekdayStyle.Render(w))
	}
	headerRow := lipgloss.JoinHorizontal(lipgloss.Center, weekdayViews...)

	// Days
	firstDay := time.Date(year, t.Month(), 1, 0, 0, 0, 0, t.Location())
	startWeekday := int(firstDay.Weekday()) // 0=Sun
	daysInMonth := time.Date(year, t.Month()+1, 0, 0, 0, 0, 0, t.Location()).Day()

	var rows []string
	var currentRow []string

	// Padding for first row
	for i := 0; i < startWeekday; i++ {
		currentRow = append(currentRow, dayStyle.Render(" "))
	}

	for d := 1; d <= daysInMonth; d++ {
		date := time.Date(year, t.Month(), d, 0, 0, 0, 0, t.Location())
		dateStr := fmt.Sprintf("%d", d)

		var style lipgloss.Style

		isCursor := date.Year() == m.calendarCursor.Year() && date.Month() == m.calendarCursor.Month() && date.Day() == m.calendarCursor.Day()
		isToday := date.Year() == time.Now().Year() && date.Month() == time.Now().Month() && date.Day() == time.Now().Day()

		if isCursor {
			style = selectedDayStyle
		} else if isToday {
			style = todayStyle
		} else {
			style = dayStyle
		}

		currentRow = append(currentRow, style.Render(dateStr))

		if len(currentRow) == 7 {
			rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Center, currentRow...))
			currentRow = []string{}
		}
	}

	// Fill last row
	if len(currentRow) > 0 {
		for len(currentRow) < 7 {
			currentRow = append(currentRow, dayStyle.Render(" "))
		}
		rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Center, currentRow...))
	}

	calendarBody := lipgloss.JoinVertical(lipgloss.Center, rows...)

	cal := lipgloss.JoinVertical(lipgloss.Center, header, headerRow, calendarBody)

	// Footer instructions
	footer := statusStyle.Render("Arrows: Move • Space: Today • Enter: Select • Esc: Cancel")

	// Center the calendar in the screen
	// Use m.width and m.height for placing
	content := lipgloss.JoinVertical(lipgloss.Center, calendarStyle.Render(cal), footer)

	return appStyle.Render(
		lipgloss.Place(m.width-10, m.height-5,
			lipgloss.Center, lipgloss.Center,
			content,
		),
	)
}
