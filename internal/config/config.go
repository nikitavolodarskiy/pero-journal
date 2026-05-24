package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

const (
	configDirName  = ".config/pero"
	configFileName = "config.toml"
	defaultEditor  = "nvim"
)

// Config holds all user-configurable settings for pero.
type Config struct {
	JournalDir string `toml:"journal_dir"`
	Editor     string `toml:"editor"`
}

// ConfigPath returns the absolute path to the config file.
func ConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("cannot find home directory: %w", err)
	}
	return filepath.Join(home, configDirName, configFileName), nil
}

// Load reads the config file and returns the resulting Config.
// If the file does not exist, sensible defaults are returned without error.
func Load() (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("cannot find home directory: %w", err)
	}

	cfg := &Config{
		JournalDir: filepath.Join(home, "pero"),
		Editor:     defaultEditor,
	}

	cfgPath := filepath.Join(home, configDirName, configFileName)
	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		return cfg, nil
	}

	if _, err := toml.DecodeFile(cfgPath, cfg); err != nil {
		return nil, fmt.Errorf("cannot parse config file %s: %w", cfgPath, err)
	}

	// Expand leading ~ in journal_dir.
	cfg.JournalDir = expandHome(cfg.JournalDir, home)

	// Environment variable overrides the config file — useful for running a
	// dev/test instance alongside the installed production app.
	//   PERO_JOURNAL_DIR=~/pero-test wails dev
	if dir := os.Getenv("PERO_JOURNAL_DIR"); dir != "" {
		cfg.JournalDir = expandHome(dir, home)
	}

	return cfg, nil
}

// Init creates the config directory and a default config file if one does not
// already exist. It returns the path of the config file.
func Init() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("cannot find home directory: %w", err)
	}

	cfgDir := filepath.Join(home, configDirName)
	if err := os.MkdirAll(cfgDir, 0o755); err != nil {
		return "", fmt.Errorf("cannot create config directory: %w", err)
	}

	cfgPath := filepath.Join(cfgDir, configFileName)

	// Do not overwrite an existing config.
	if _, err := os.Stat(cfgPath); err == nil {
		return cfgPath, nil
	}

	f, err := os.Create(cfgPath)
	if err != nil {
		return "", fmt.Errorf("cannot create config file: %w", err)
	}
	defer f.Close()

	defaults := Config{
		JournalDir: filepath.Join(home, "pero"),
		Editor:     defaultEditor,
	}

	if err := toml.NewEncoder(f).Encode(defaults); err != nil {
		return "", fmt.Errorf("cannot write default config: %w", err)
	}

	return cfgPath, nil
}

// expandHome replaces a leading ~ with the provided home directory.
func expandHome(path, home string) string {
	if path == "~" {
		return home
	}
	if strings.HasPrefix(path, "~/") {
		return filepath.Join(home, path[2:])
	}
	return path
}
