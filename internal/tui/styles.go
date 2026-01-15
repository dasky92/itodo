package tui

import "github.com/charmbracelet/lipgloss"

var (
	// Styles
	appStyle                lipgloss.Style
	dateStyle               lipgloss.Style
	titleStyle              lipgloss.Style
	statsStyle              lipgloss.Style
	progressStyle           lipgloss.Style
	progressEmptyStyle      lipgloss.Style
	dividerStyle            lipgloss.Style
	selectedItemStyle       lipgloss.Style
	completedItemStyle      lipgloss.Style
	normalItemStyle         lipgloss.Style
	statusStyle             lipgloss.Style
	logStyle                lipgloss.Style
	formTitleStyle          lipgloss.Style
	labelStyle              lipgloss.Style
	calendarStyle           lipgloss.Style
	monthTitleStyle         lipgloss.Style
	weekdayStyle            lipgloss.Style
	dayStyle                lipgloss.Style
	selectedDayStyle        lipgloss.Style
	todayStyle              lipgloss.Style
	dayMarkerStyle          lipgloss.Style
	dayMarkerPendingStyle   lipgloss.Style
	dayMarkerCompletedStyle lipgloss.Style

	// Characters
	progressFullChar  = "█"
	progressEmptyChar = "·"
)

func InitStyles(theme Theme) {
	CurrentTheme = theme

	appStyle = lipgloss.NewStyle().Padding(1, 4)

	dateStyle = lipgloss.NewStyle().
		Foreground(theme.Mauve).
		Bold(true).
		Padding(0, 1)

	titleStyle = lipgloss.NewStyle().
		Foreground(theme.Text).
		Bold(true).
		Padding(0, 1)

	statsStyle = lipgloss.NewStyle().
		Foreground(theme.Green).
		Bold(true).
		Padding(0, 1)

	progressStyle = lipgloss.NewStyle().Foreground(theme.Green)
	progressEmptyStyle = lipgloss.NewStyle().Foreground(theme.Gray)

	dividerStyle = lipgloss.NewStyle().
		Foreground(theme.Surface).
		Padding(1, 0)

	selectedItemStyle = lipgloss.NewStyle().
		Foreground(theme.Pink).
		Bold(true).
		PaddingLeft(1).
		Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(theme.Border).
		MarginBottom(1)

	completedItemStyle = lipgloss.NewStyle().
		Foreground(theme.Gray).
		PaddingLeft(2).
		MarginBottom(1)

	normalItemStyle = lipgloss.NewStyle().
		Foreground(theme.Text).
		PaddingLeft(2).
		MarginBottom(1)

	statusStyle = lipgloss.NewStyle().
		Foreground(theme.Gray).
		Italic(true).
		Padding(0, 1)

	logStyle = lipgloss.NewStyle().
		Foreground(theme.Gray).
		Faint(true)

	formTitleStyle = lipgloss.NewStyle().
		Foreground(theme.Mauve).
		Bold(true).
		PaddingBottom(1)

	labelStyle = lipgloss.NewStyle().
		Foreground(theme.Blue).
		Bold(true).
		MarginTop(1)

	calendarStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(theme.Mauve).
		PaddingTop(1).
		PaddingLeft(2).
		PaddingRight(2)

	monthTitleStyle = lipgloss.NewStyle().
		Foreground(theme.Blue).
		Bold(true).
		Align(lipgloss.Center).
		Width(28)

	weekdayStyle = lipgloss.NewStyle().
		Foreground(theme.Gray).
		Width(4).
		Align(lipgloss.Center)

	dayStyle = lipgloss.NewStyle().
		Foreground(theme.Text).
		Width(4).
		Align(lipgloss.Center)

	selectedDayStyle = lipgloss.NewStyle().
		Foreground(theme.Base).
		Background(theme.Pink).
		Bold(true).
		Width(4).
		Align(lipgloss.Center)

	todayStyle = lipgloss.NewStyle().
		Foreground(theme.Blue).
		Bold(true).
		Width(4).
		Align(lipgloss.Center)

	dayMarkerStyle = lipgloss.NewStyle().
		Foreground(theme.Gray).
		Width(4).
		Align(lipgloss.Center)

	dayMarkerPendingStyle = lipgloss.NewStyle().
		Foreground(theme.Red).
		Width(4).
		Align(lipgloss.Center)

	dayMarkerCompletedStyle = lipgloss.NewStyle().
		Foreground(theme.Green).
		Width(4).
		Align(lipgloss.Center)
}
