package main

import (
	"fmt"
	"os"

	"itodo/model"
	"itodo/service"
	"itodo/tui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Initialize Database
	if err := model.InitDB("itodo.db"); err != nil {
		fmt.Printf("Failed to initialize database: %v\n", err)
		os.Exit(1)
	}

	// Initialize Service
	svc := service.NewTodoService()

	// Initialize TUI
	p := tea.NewProgram(tui.NewModel(svc), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
