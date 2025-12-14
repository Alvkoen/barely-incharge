package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

const (
	ModeCrunch = "crunch"
	ModeNormal = "normal"
	ModeSaver  = "saver"
)

var ValidModes = []string{ModeCrunch, ModeNormal, ModeSaver}

type Config struct {
	WorkHours        TimeRange `json:"work_hours"`
	LunchTime        TimeRange `json:"lunch_time"`
	MeetingsCalendar string    `json:"meetings_calendar"`
	BlocksCalendar   string    `json:"blocks_calendar"`
	DefaultMode      string    `json:"default_mode"`
}

type TimeRange struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

func GetConfigPath() (string, error) {
	execPath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("failed to get executable path: %w", err)
	}

	execDir := filepath.Dir(execPath)
	return filepath.Join(execDir, "config.json"), nil
}

func Load() (*Config, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// Trim whitespace from string fields
	cfg.DefaultMode = strings.TrimSpace(cfg.DefaultMode)
	cfg.MeetingsCalendar = strings.TrimSpace(cfg.MeetingsCalendar)
	cfg.BlocksCalendar = strings.TrimSpace(cfg.BlocksCalendar)
	cfg.WorkHours.Start = strings.TrimSpace(cfg.WorkHours.Start)
	cfg.WorkHours.End = strings.TrimSpace(cfg.WorkHours.End)
	cfg.LunchTime.Start = strings.TrimSpace(cfg.LunchTime.Start)
	cfg.LunchTime.End = strings.TrimSpace(cfg.LunchTime.End)

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func IsValidMode(mode string) bool {
	return slices.Contains(ValidModes, mode)
}

// ValidateMode validates a mode and returns a descriptive error if invalid
func ValidateMode(mode string) error {
	if !IsValidMode(mode) {
		return fmt.Errorf("invalid mode: %s (valid modes: %s)",
			mode, strings.Join(ValidModes, ", "))
	}
	return nil
}

func (c *Config) Validate() error {
	if err := ValidateMode(c.DefaultMode); err != nil {
		return fmt.Errorf("invalid default_mode in config: %w", err)
	}
	//todo other validations here
	return nil
}
