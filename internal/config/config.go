package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"
)

const (
	ModeCrunch = "crunch"
	ModeNormal = "normal"
	ModeSaver  = "saver"

	// DateFormat is the expected format for the date field in config (YYYY-MM-DD)
	DateFormat = "2006-01-02"
)

var ValidModes = []string{ModeCrunch, ModeNormal, ModeSaver}

type Config struct {
	WorkHours        TimeRange `json:"work_hours"`
	LunchTime        TimeRange `json:"lunch_time"`
	MeetingsCalendar string    `json:"meetings_calendar"`
	BlocksCalendar   string    `json:"blocks_calendar"`
	DefaultMode      string    `json:"default_mode"`
	OpenAIAPIKey     string    `json:"openai_api_key"`
	Date             string    `json:"date"` // YYYY-MM-DD format, empty = today
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

	if c.Date != "" {
		if _, err := time.Parse(DateFormat, c.Date); err != nil {
			return fmt.Errorf("invalid date format (expected YYYY-MM-DD): %w", err)
		}
	}

	return nil
}

// GetPlanningDate returns the date to plan for. Returns today if date is not set in config.
// The returned time is at midnight in the local timezone - the actual times will be set
// when parsing time strings like "09:00" via ParseTimeOnDate.
func (c *Config) GetPlanningDate() (time.Time, error) {
	if c.Date == "" {
		now := time.Now()
		return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()), nil
	}

	date, err := time.Parse(DateFormat, c.Date)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid date format: %w", err)
	}

	loc := time.Now().Location()
	return time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, loc), nil
}
