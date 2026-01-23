# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

itodo is a keyboard-driven terminal-based todo list manager built with Go and Bubble Tea (TUI framework). It provides daily and weekly views, task management with hierarchical support, and multiple themes.

## Architecture

The application follows a layered architecture:

### Core Layers
1. **Model Layer** (`internal/model/`) - Data persistence with GORM/SQLite
2. **Service Layer** (`internal/service/`) - Business logic and validation
3. **TUI Layer** (`internal/tui/`) - User interface using Bubble Tea components

### Key Components

#### Model (`internal/model/todo.go`)
- Uses GORM with SQLite for persistence
- `Todo` struct with composite key (Date + ID)
- Manual ID generation per date (starts at 1 each day)
- Supports hierarchical tasks via `ParentID`
- Sorting: Pending → Completed, then by Position → ID
- Position-based ordering within siblings

#### Service (`internal/service/todo.go`)
- Business logic and validation
- Prevents past-date task creation
- Enforces daily task limits (50 tasks max)
- Handles task relationships (indent/outdent)
- Date manipulation utilities

#### TUI (`internal/tui/`)
- **Model**: State management, cursor positioning, view switching
- **Update**: Event handling, key bindings, form management
- **View**: Rendering, themes, calendar integration
- **Keys**: Configurable keyboard bindings
- **Styles**: Lipgloss theming system (Monokai, OneDark, Catppuccin, etc.)

#### Configuration (`internal/config/`)
- JSON-based config at `$HOME/.config/itodo/config.json`
- Key binding customization
- Theme selection
- UI settings (default view, line numbers, etc.)

## Development Commands

### Build & Run
```bash
# Build the binary
make build

# Install to ~/.local/bin
make install

# Run directly
go run .

# Clean build artifacts
make clean
```

### Testing
```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests for specific package
go test ./internal/model/
go test ./internal/service/
```

### Demo
```bash
# Generate demo.gif (requires vhs and ffmpeg)
make demo
```

### Release
```bash
# Test release build
make publish-test

# Actual release
make publish
```

## Tech Stack

- **Go 1.24+** - Main language (managed by mise)
- **Bubble Tea** - TUI framework
- **Lipgloss** - Styling
- **GORM** - ORM
- **SQLite** - Database
- **GoReleaser** - Build automation

## Database Schema

The `Todo` model has these key fields:
- `ID`: Manual ID per day (resets daily)
- `Date`: YYYY-MM-DD format (part of composite primary key)
- `ParentID`: Hierarchical relationships (within same date)
- `Position`: Ordering within siblings
- `IsCompleted`: Task status

## Important Patterns

### Key Bindings
- Configurable via JSON config
- Supports multiple keys per action
- Fallback to defaults if not configured

### View System
- **Daily View**: Shows tasks for selected date
- **Weekly View**: Shows 7-day range with expandable sections
- **Calendar View**: Monthly overview with task markers

### Task Management
- Tasks can be indented/outdent for hierarchy
- Movement restricted within same completion status
- Automatic position management on reordering

### Theme System
- Multiple built-in themes (Monokai, OneDark, Catppuccin, etc.)
- Configured via `cfg.General.Theme`
- Uses Lipgloss styles throughout

## Configuration File Location

Default config: `$HOME/.config/itodo/config.json`
Fallback config: `config.default.json` in project root

## Go Version

Uses Go 1.24.12 (managed by mise.toml)

## Dependencies

Key dependencies:
- `github.com/charmbracelet/bubbletea` - TUI framework
- `github.com/charmbracelet/lipgloss` - Styling
- `github.com/glebarez/sqlite` - SQLite driver
- `gorm.io/gorm` - ORM