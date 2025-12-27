package tui

import (
	"itodo/internal/config"

	"github.com/charmbracelet/bubbles/key"
)

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
	Save     key.Binding
	Cancel   key.Binding
	Today    key.Binding
	Calendar key.Binding
	PrevView key.Binding
	NextView key.Binding
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

// Global default keys, will be overwritten by config
var Keys = KeyMap{}

// InitKeys initializes the key bindings from configuration
func InitKeys(cfg *config.Config) KeyMap {
	k := KeyMap{
		Up: key.NewBinding(
			key.WithKeys(cfg.Keys.Up...),
			key.WithHelp(formatHelp(cfg.Keys.Up, "up")),
		),
		Down: key.NewBinding(
			key.WithKeys(cfg.Keys.Down...),
			key.WithHelp(formatHelp(cfg.Keys.Down, "down")),
		),
		Left: key.NewBinding(
			key.WithKeys(cfg.Keys.Left...),
			key.WithHelp(formatHelp(cfg.Keys.Left, "prev day")),
		),
		Right: key.NewBinding(
			key.WithKeys(cfg.Keys.Right...),
			key.WithHelp(formatHelp(cfg.Keys.Right, "next day")),
		),
		New: key.NewBinding(
			key.WithKeys(cfg.Keys.New...),
			key.WithHelp(formatHelp(cfg.Keys.New, "new todo")),
		),
		Edit: key.NewBinding(
			key.WithKeys(cfg.Keys.Edit...),
			key.WithHelp(formatHelp(cfg.Keys.Edit, "edit todo")),
		),
		Delete: key.NewBinding(
			key.WithKeys(cfg.Keys.Delete...),
			key.WithHelp(formatHelp(cfg.Keys.Delete, "delete")),
		),
		Toggle: key.NewBinding(
			key.WithKeys(cfg.Keys.Toggle...),
			key.WithHelp(formatHelp(cfg.Keys.Toggle, "toggle")),
		),
		Indent: key.NewBinding(
			key.WithKeys(cfg.Keys.Indent...),
			key.WithHelp(formatHelp(cfg.Keys.Indent, "indent")),
		),
		Outdent: key.NewBinding(
			key.WithKeys(cfg.Keys.Outdent...),
			key.WithHelp(formatHelp(cfg.Keys.Outdent, "outdent")),
		),
		Help: key.NewBinding(
			key.WithKeys(cfg.Keys.Help...),
			key.WithHelp(formatHelp(cfg.Keys.Help, "toggle help")),
		),
		Quit: key.NewBinding(
			key.WithKeys(cfg.Keys.Quit...),
			key.WithHelp(formatHelp(cfg.Keys.Quit, "quit")),
		),
		Save: key.NewBinding(
			key.WithKeys(cfg.Keys.Save...),
			key.WithHelp(formatHelp(cfg.Keys.Save, "save")),
		),
		Cancel: key.NewBinding(
			key.WithKeys(cfg.Keys.Cancel...),
			key.WithHelp(formatHelp(cfg.Keys.Cancel, "cancel")),
		),
		Today: key.NewBinding(
			key.WithKeys(cfg.Keys.Today...),
			key.WithHelp(formatHelp(cfg.Keys.Today, "go to today")),
		),
		Calendar: key.NewBinding(
			key.WithKeys(cfg.Keys.Calendar...),
			key.WithHelp(formatHelp(cfg.Keys.Calendar, "calendar")),
		),
		PrevView: key.NewBinding(
			key.WithKeys(cfg.Keys.PrevView...),
			key.WithHelp(formatHelp(cfg.Keys.PrevView, "daily view")),
		),
		NextView: key.NewBinding(
			key.WithKeys(cfg.Keys.NextView...),
			key.WithHelp(formatHelp(cfg.Keys.NextView, "weekly view")),
		),
	}
	Keys = k // Update global keys for fallback or direct access if needed
	return k
}

func formatHelp(keys []string, desc string) (string, string) {
	if len(keys) == 0 {
		return "", desc
	}
	return keys[0], desc
}
