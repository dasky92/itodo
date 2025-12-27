package tui

import "github.com/charmbracelet/lipgloss"

var (
	// Catppuccin Mocha Palette
	colBase     = lipgloss.Color("#1e1e2e")
	colText     = lipgloss.Color("#cdd6f4")
	colBlue     = lipgloss.Color("#89b4fa")
	colGreen    = lipgloss.Color("#a6e3a1")
	colPink     = lipgloss.Color("#f5c2e7")
	colGray     = lipgloss.Color("#6c7086")
	colDarkGray = lipgloss.Color("#45475a")
	colMauve    = lipgloss.Color("#cba6f7")
	colRed      = lipgloss.Color("#f38ba8")
	colSurface0 = lipgloss.Color("#313244")

	// Global App Padding
	appStyle = lipgloss.NewStyle().Padding(1, 4)

	// Section Headers
	dateStyle = lipgloss.NewStyle().
			Foreground(colMauve).
			Bold(true).
			Padding(0, 1)

	titleStyle = lipgloss.NewStyle().
			Foreground(colText).
			Bold(true).
			Padding(0, 1)

	statsStyle = lipgloss.NewStyle().
			Foreground(colGreen).
			Bold(true).
			Padding(0, 1)

	progressFullChar   = "█"
	progressEmptyChar  = "░"
	progressStyle      = lipgloss.NewStyle().Foreground(colGreen)
	progressEmptyStyle = lipgloss.NewStyle().Foreground(colSurface0)

	// Divider
	dividerStyle = lipgloss.NewStyle().
			Foreground(colSurface0).
			Padding(1, 0) // Vertical breathing room

	// List Items
	selectedItemStyle = lipgloss.NewStyle().
				Foreground(colPink).
				Bold(true).
				PaddingLeft(1).
				Border(lipgloss.NormalBorder(), false, false, false, true).
				BorderForeground(colPink).
				MarginBottom(1) // Add breathing room

	completedItemStyle = lipgloss.NewStyle().
				Foreground(colGray).
				Strikethrough(true).
				PaddingLeft(2).
				MarginBottom(1) // Add breathing room

	normalItemStyle = lipgloss.NewStyle().
			Foreground(colText).
			PaddingLeft(2).
			MarginBottom(1) // Add breathing room

	// Footer / Logs
	statusStyle = lipgloss.NewStyle().
			Foreground(colGray).
			Italic(true).
			Padding(0, 1)

	logStyle = lipgloss.NewStyle().
			Foreground(colGray).
			Faint(true)

	// Form
	formTitleStyle = lipgloss.NewStyle().
			Foreground(colMauve).
			Bold(true).
			PaddingBottom(1)

	labelStyle = lipgloss.NewStyle().
			Foreground(colBlue).
			Bold(true).
			MarginTop(1)
)
