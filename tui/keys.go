package tui

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Up       key.Binding
	Down     key.Binding
	Left     key.Binding
	Right    key.Binding
	New      key.Binding
	Edit     key.Binding
	Delete   key.Binding
	Toggle   key.Binding
	Indent   key.Binding
	Outdent  key.Binding
	Help     key.Binding
	Quit     key.Binding
	Save     key.Binding // New
	Cancel   key.Binding // New
	Today    key.Binding // New: Go to Today
	Calendar key.Binding // New: Open Calendar
	PrevView key.Binding // New: Switch to Daily
	NextView key.Binding // New: Switch to Weekly
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit, k.Calendar, k.PrevView, k.NextView}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right},
		{k.New, k.Edit, k.Delete, k.Toggle},
		{k.Indent, k.Outdent, k.Help, k.Quit},
		{k.Save, k.Cancel, k.Today, k.Calendar},
		{k.PrevView, k.NextView},
	}
}

var Keys = KeyMap{
	Up: key.NewBinding(
		key.WithKeys("k", "up"),
		key.WithHelp("k/↑", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("j", "down"),
		key.WithHelp("j/↓", "down"),
	),
	Left: key.NewBinding(
		key.WithKeys("h", "left"),
		key.WithHelp("h/←", "prev day"),
	),
	Right: key.NewBinding(
		key.WithKeys("l", "right"),
		key.WithHelp("l/→", "next day"),
	),
	New: key.NewBinding(
		key.WithKeys("n"),
		key.WithHelp("n", "new todo"),
	),
	Edit: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "edit todo"),
	),
	Delete: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "delete"),
	),
	Toggle: key.NewBinding(
		key.WithKeys("enter", "tab"),
		key.WithHelp("enter", "toggle"),
	),
	Indent: key.NewBinding(
		key.WithKeys(">", "."),
		key.WithHelp(">", "indent"),
	),
	Outdent: key.NewBinding(
		key.WithKeys("<", ","),
		key.WithHelp("<", "outdent"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Save: key.NewBinding(
		key.WithKeys("ctrl+s"),
		key.WithHelp("ctrl+s", "save"),
	),
	Cancel: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "cancel"),
	),
	Today: key.NewBinding(
		key.WithKeys(" "),
		key.WithHelp("space", "go to today"),
	),
	Calendar: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "calendar"),
	),
	PrevView: key.NewBinding(
		key.WithKeys("H"),
		key.WithHelp("H", "daily"),
	),
	NextView: key.NewBinding(
		key.WithKeys("L"),
		key.WithHelp("L", "weekly"),
	),
}
