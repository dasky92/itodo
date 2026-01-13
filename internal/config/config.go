package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
)

type GeneralConfig struct {
	DBName string `json:"db_name"`
	Theme  string `json:"theme"`

	// Internal field, not exported to JSON
	DBPath string `json:"-"`
}

type UIConfig struct {
	DefaultView     string `json:"default_view"`
	ShowLineNumbers bool   `json:"show_line_numbers"`
	ShowFullHelp    bool   `json:"show_full_help"`
}

type InputConfig struct {
	TitleCharLimit int `json:"title_char_limit"`
	TitleWidth     int `json:"title_width"`
	DescHeight     int `json:"desc_height"`
	DescWidth      int `json:"desc_width"`
}

type KeysConfig struct {
	Up       []string `json:"up"`
	Down     []string `json:"down"`
	Left     []string `json:"left"`
	Right    []string `json:"right"`
	New      []string `json:"new"`
	Edit     []string `json:"edit"`
	Delete   []string `json:"delete"`
	Toggle   []string `json:"toggle"`
	Indent   []string `json:"indent"`
	Outdent  []string `json:"outdent"`
	MoveUp   []string `json:"move_up"`
	MoveDown []string `json:"move_down"`
	Help     []string `json:"help"`
	Quit     []string `json:"quit"`
	Save     []string `json:"save"`
	Cancel   []string `json:"cancel"`
	Today    []string `json:"today"`
	Calendar []string `json:"calendar"`
	PrevView []string `json:"prev_view"`
	NextView []string `json:"next_view"`
}

type Config struct {
	General GeneralConfig `json:"general"`
	UI      UIConfig      `json:"ui"`
	Input   InputConfig   `json:"input"`
	Keys    KeysConfig    `json:"keys"`
}

func DefaultConfig() *Config {
	return &Config{
		General: GeneralConfig{
			DBName: "itodo.db",
			Theme:  "Monokai",
		},
		UI: UIConfig{
			DefaultView:     "daily",
			ShowLineNumbers: false,
			ShowFullHelp:    false,
		},
		Input: InputConfig{
			TitleCharLimit: 100,
			TitleWidth:     50,
			DescHeight:     10,
			DescWidth:      50,
		},
		Keys: KeysConfig{
			Up:       []string{"k", "up"},
			Down:     []string{"j", "down"},
			Left:     []string{"h", "left"},
			Right:    []string{"l", "right"},
			New:      []string{"n"},
			Edit:     []string{"i"},
			Delete:   []string{"d"},
			Toggle:   []string{"enter", "tab"},
			Indent:   []string{">", "."},
			Outdent:  []string{"<", ","},
			MoveUp:   []string{"K"},
			MoveDown: []string{"J"},
			Help:     []string{"?"},
			Quit:     []string{"q", "ctrl+c"},
			Save:     []string{"ctrl+s"},
			Cancel:   []string{"esc"},
			Today:    []string{" "},
			Calendar: []string{";"},
			PrevView: []string{"H"},
			NextView: []string{"L"},
		},
	}
}

func LoadConfig() (*Config, error) {
	cfg := DefaultConfig()

	var configDir string
	if runtime.GOOS == "windows" {
		ucd, err := os.UserConfigDir()
		if err == nil {
			configDir = filepath.Join(ucd, "itodo")
		}
	}

	if configDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return cfg, err // Return default if can't get home dir
		}
		configDir = filepath.Join(home, ".config", "itodo")
	}

	configPath := filepath.Join(configDir, "config.json")

	// Ensure config directory exists and create example config
	if errMkDir := os.MkdirAll(configDir, 0755); errMkDir == nil {
		defaultPath := filepath.Join(configDir, "config.default.json")
		if _, errStat := os.Stat(defaultPath); os.IsNotExist(errStat) {
			// Create default config with helpful defaults and comments
			defaultCfg := DefaultConfig()

			type CommentedGeneral struct {
				Help   string `json:"_help"`
				DBName string `json:"db_name"`
				Theme  string `json:"theme"`
			}
			type CommentedUI struct {
				Help            string `json:"_help"`
				DefaultView     string `json:"default_view"`
				ShowLineNumbers bool   `json:"show_line_numbers"`
				ShowFullHelp    bool   `json:"show_full_help"`
			}
			type CommentedInput struct {
				Help           string `json:"_help"`
				TitleCharLimit int    `json:"title_char_limit"`
				TitleWidth     int    `json:"title_width"`
				DescHeight     int    `json:"desc_height"`
				DescWidth      int    `json:"desc_width"`
			}
			type CommentedKeys struct {
				Help string `json:"_help"`
				KeysConfig
			}

			type CommentedConfig struct {
				General CommentedGeneral `json:"general"`
				UI      CommentedUI      `json:"ui"`
				Input   CommentedInput   `json:"input"`
				Keys    CommentedKeys    `json:"keys"`
			}

			cc := CommentedConfig{
				General: CommentedGeneral{
					Help:   "General settings. db_name is the database filename in $HOME.config/itodo. Themes: Monokai, OneDark, OneLight, Hacker, Catppuccin",
					DBName: defaultCfg.General.DBName,
					Theme:  defaultCfg.General.Theme,
				},
				UI: CommentedUI{
					Help:            "UI settings. default_view: daily or weekly",
					DefaultView:     defaultCfg.UI.DefaultView,
					ShowLineNumbers: defaultCfg.UI.ShowLineNumbers,
					ShowFullHelp:    defaultCfg.UI.ShowFullHelp,
				},
				Input: CommentedInput{
					Help:           "Input field dimensions and limits",
					TitleCharLimit: defaultCfg.Input.TitleCharLimit,
					TitleWidth:     defaultCfg.Input.TitleWidth,
					DescHeight:     defaultCfg.Input.DescHeight,
					DescWidth:      defaultCfg.Input.DescWidth,
				},
				Keys: CommentedKeys{
					Help:       "Key bindings. Use array of strings. Modifiers: ctrl, alt, shift",
					KeysConfig: defaultCfg.Keys,
				},
			}

			if data, errMarshal := json.MarshalIndent(cc, "", "  "); errMarshal == nil {
				_ = os.WriteFile(defaultPath, data, 0644)
			}
		}
	}

	file, err := os.Open(configPath)
	if os.IsNotExist(err) {
		// Even if config file doesn't exist, we need to set DBPath based on default DBName
		cfg.General.DBPath = filepath.Join(configDir, cfg.General.DBName)
		return cfg, nil
	}
	if err != nil {
		return cfg, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(cfg); err != nil {
		return cfg, err
	}

	// Fallback if DBName is empty
	if cfg.General.DBName == "" {
		cfg.General.DBName = "itodo.db"
	}

	// Construct full DB path
	cfg.General.DBPath = filepath.Join(configDir, cfg.General.DBName)

	return cfg, nil
}
