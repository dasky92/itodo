package main

import (
	"fmt"
	"os"

	"github.com/dasky92/itodo/internal/config"
	"github.com/dasky92/itodo/internal/model"
	"github.com/dasky92/itodo/internal/service"
	"github.com/dasky92/itodo/internal/tui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Load Configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Warning: Failed to load config: %v\nUsing defaults.\n", err)
	}

	// Initialize Database
	if err := model.InitDB(cfg.General.DBPath); err != nil {
		fmt.Printf("Failed to initialize database: %v\n", err)
		os.Exit(1)
	}

	// Initialize Service
	svc := service.NewTodoService()

	// Initialize TUI
	p := tea.NewProgram(tui.NewModel(svc, cfg), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
