package tui

import "github.com/charmbracelet/lipgloss"

type Theme struct {
	Base      lipgloss.Color
	Text      lipgloss.Color
	Blue      lipgloss.Color
	Green     lipgloss.Color
	Pink      lipgloss.Color
	Gray      lipgloss.Color
	DarkGray  lipgloss.Color
	Mauve     lipgloss.Color
	Red       lipgloss.Color
	Surface   lipgloss.Color
	Highlight lipgloss.Color
	Border    lipgloss.Color
}

var (
	// Current Theme Palette
	CurrentTheme Theme

	// Styles
	appStyle           lipgloss.Style
	dateStyle          lipgloss.Style
	titleStyle         lipgloss.Style
	statsStyle         lipgloss.Style
	progressStyle      lipgloss.Style
	progressEmptyStyle lipgloss.Style
	dividerStyle       lipgloss.Style
	selectedItemStyle  lipgloss.Style
	completedItemStyle lipgloss.Style
	normalItemStyle    lipgloss.Style
	statusStyle        lipgloss.Style
	logStyle           lipgloss.Style
	formTitleStyle     lipgloss.Style
	labelStyle         lipgloss.Style
	calendarStyle      lipgloss.Style
	monthTitleStyle    lipgloss.Style
	weekdayStyle       lipgloss.Style
	dayStyle           lipgloss.Style
	selectedDayStyle   lipgloss.Style
	todayStyle         lipgloss.Style

	// Characters
	progressFullChar  = "█"
	progressEmptyChar = "·"
)

// Define Themes
var Themes = map[string]Theme{
	"Catppuccin": {
		Base:      lipgloss.Color("#1e1e2e"),
		Text:      lipgloss.Color("#cdd6f4"),
		Blue:      lipgloss.Color("#89b4fa"),
		Green:     lipgloss.Color("#a6e3a1"),
		Pink:      lipgloss.Color("#f5c2e7"),
		Gray:      lipgloss.Color("#6c7086"),
		DarkGray:  lipgloss.Color("#45475a"),
		Mauve:     lipgloss.Color("#cba6f7"),
		Red:       lipgloss.Color("#f38ba8"),
		Surface:   lipgloss.Color("#313244"),
		Highlight: lipgloss.Color("#313244"),
		Border:    lipgloss.Color("#f5c2e7"),
	},
	"Monokai": {
		Base:      lipgloss.Color("#272822"),
		Text:      lipgloss.Color("#F8F8F2"),
		Blue:      lipgloss.Color("#66D9EF"),
		Green:     lipgloss.Color("#A6E22E"),
		Pink:      lipgloss.Color("#F92672"), // Using Pink slot for Red/Pink
		Gray:      lipgloss.Color("#75715E"),
		DarkGray:  lipgloss.Color("#49483E"),
		Mauve:     lipgloss.Color("#AE81FF"),
		Red:       lipgloss.Color("#F92672"),
		Surface:   lipgloss.Color("#3E3D32"),
		Highlight: lipgloss.Color("#3E3D32"),
		Border:    lipgloss.Color("#F92672"),
	},
	"OneDark": {
		Base:      lipgloss.Color("#282c34"),
		Text:      lipgloss.Color("#abb2bf"),
		Blue:      lipgloss.Color("#61afef"),
		Green:     lipgloss.Color("#98c379"),
		Pink:      lipgloss.Color("#e06c75"),
		Gray:      lipgloss.Color("#5c6370"),
		DarkGray:  lipgloss.Color("#4b5263"),
		Mauve:     lipgloss.Color("#c678dd"),
		Red:       lipgloss.Color("#e06c75"),
		Surface:   lipgloss.Color("#3e4452"),
		Highlight: lipgloss.Color("#3e4452"),
		Border:    lipgloss.Color("#61afef"),
	},
	"OneLight": {
		Base:      lipgloss.Color("#fafafa"),
		Text:      lipgloss.Color("#383a42"),
		Blue:      lipgloss.Color("#4078f2"),
		Green:     lipgloss.Color("#50a14f"),
		Pink:      lipgloss.Color("#e45649"),
		Gray:      lipgloss.Color("#a0a1a7"),
		DarkGray:  lipgloss.Color("#d0d0d0"),
		Mauve:     lipgloss.Color("#a626a4"),
		Red:       lipgloss.Color("#e45649"),
		Surface:   lipgloss.Color("#e5e5e6"),
		Highlight: lipgloss.Color("#e5e5e6"),
		Border:    lipgloss.Color("#4078f2"),
	},
	"Hacker": {
		Base:      lipgloss.Color("#000000"),
		Text:      lipgloss.Color("#00FF00"),
		Blue:      lipgloss.Color("#008800"),
		Green:     lipgloss.Color("#00FF00"),
		Pink:      lipgloss.Color("#00AA00"),
		Gray:      lipgloss.Color("#004400"),
		DarkGray:  lipgloss.Color("#002200"),
		Mauve:     lipgloss.Color("#00DD00"),
		Red:       lipgloss.Color("#FF0000"), // Alert
		Surface:   lipgloss.Color("#001100"),
		Highlight: lipgloss.Color("#003300"),
		Border:    lipgloss.Color("#00FF00"),
	},
}

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
		Padding(1, 2)

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
}
