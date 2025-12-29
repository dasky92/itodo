# itodo

<p align="center">
  <img src="demo.gif" alt="itodo demo" width="600" />
</p>

<p align="center">
  A modern, keyboard-driven terminal-based todo list manager built with <a href="https://github.com/charmbracelet/bubbletea">Bubble Tea</a>.
</p>

<p align="center">
  <a href="#features">Features</a> •
  <a href="#installation">Installation</a> •
  <a href="#usage">Usage</a> •
  <a href="#configuration">Configuration</a>
</p>

## ✨ Features

*   **Keyboard-centric workflow**: Manage your tasks without leaving the keyboard.
*   **Multiple Views**: Switch between Daily and Weekly views to plan your time effectively.
*   **Themable**: Comes with built-in themes like Monokai, OneDark, OneLight, Hacker, and Catppuccin.
*   **Persistent Storage**: Tasks are saved in a local SQLite database.
*   **Configurable**: Customize key bindings and UI settings via a JSON configuration file.
*   **Cross-Platform**: Runs on Linux, macOS, and Windows.

## 🚀 Installation

### Homebrew (macOS/Linux)

```bash
# Add the tap
brew tap dasky92/tap

# Install itodo
brew install --cask itodo
```

### Chocolatey (Windows)

```powershell
choco install itodo
```

### Curl (Linux/macOS)

You can use the installation script to download and install the binary directly:

```bash
curl -sL https://raw.githubusercontent.com/dasky92/itodo/main/install.sh | bash
```

### Go Install

If you have Go installed, you can build from source:

```bash
go install github.com/dasky92/itodo@latest
```

## ⌨️ Usage

### Key Bindings

| Action | Key(s) |
| :--- | :--- |
| **Navigation** | |
| Move Up | `k` / `↑` |
| Move Down | `j` / `↓` |
| Move Left | `h` / `←` |
| Move Right | `l` / `→` |
| Previous View | `H` |
| Next View | `L` |
| **Task Management** | |
| New Task | `n` |
| Edit Task | `i` |
| Delete Task | `d` |
| Toggle Complete | `Enter` / `Tab` |
| Indent Task | `>` / `.` |
| Outdent Task | `<` / `,` |
| **General** | |
| Help | `?` |
| Quit | `q` / `Ctrl+c` |
| Save | `Ctrl+s` |
| Cancel | `Esc` |
| Go to Today | `Space` |
| Calendar | `;` |

## ⚙️ Configuration

`itodo` creates a configuration file at `$HOME/.config/itodo/config.json` on first run. You can also find a `config.default.json` in the same directory for reference.

### Example Configuration

```json
{
  "general": {
    "db_name": "itodo.db",
    "theme": "Monokai"
  },
  "ui": {
    "default_view": "daily",
    "show_line_numbers": false
  },
  "keys": {
    "up": ["k", "up"],
    "down": ["j", "down"],
    "new": ["n"],
    "quit": ["q", "ctrl+c"]
  }
}
```

### Available Themes

*   `Monokai` (Default)
*   `OneDark`
*   `OneLight`
*   `Hacker`
*   `Catppuccin`

## 📝 License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
