package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

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

	return &cfg, nil
}
