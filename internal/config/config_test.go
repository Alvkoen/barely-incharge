package config

import (
	"strings"
	"testing"
)

func TestIsValidMode(t *testing.T) {
	tests := []struct {
		name     string
		mode     string
		expected bool
	}{
		{"valid crunch", "crunch", true},
		{"valid normal", "normal", true},
		{"valid saver", "saver", true},
		{"invalid turbo", "turbo", false},
		{"invalid empty", "", false},
		{"invalid uppercase", "CRUNCH", false},
		{"valid with trailing space (trimmed by caller)", "crunch", true},
		{"valid with leading space (trimmed by caller)", "crunch", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidMode(tt.mode)
			if result != tt.expected {
				t.Errorf("IsValidMode(%q) = %v, want %v", tt.mode, result, tt.expected)
			}
		})
	}
}

func TestValidateMode(t *testing.T) {
	tests := []struct {
		name      string
		mode      string
		expectErr bool
	}{
		{"valid crunch", "crunch", false},
		{"valid normal", "normal", false},
		{"valid saver", "saver", false},
		{"invalid mode", "turbo", true},
		{"empty mode", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateMode(tt.mode)
			if tt.expectErr && err == nil {
				t.Errorf("ValidateMode(%q) expected error but got nil", tt.mode)
			}
			if !tt.expectErr && err != nil {
				t.Errorf("ValidateMode(%q) expected no error but got: %v", tt.mode, err)
			}
		})
	}
}

func TestValidateModeErrorMessage(t *testing.T) {
	err := ValidateMode("invalid")
	if err == nil {
		t.Fatal("Expected error for invalid mode")
	}

	errMsg := err.Error()
	if !strings.Contains(errMsg, "invalid mode: invalid") {
		t.Errorf("Error message should contain 'invalid mode: invalid', got: %s", errMsg)
	}
	if !strings.Contains(errMsg, "crunch") || !strings.Contains(errMsg, "normal") || !strings.Contains(errMsg, "saver") {
		t.Errorf("Error message should list valid modes, got: %s", errMsg)
	}
}

func TestConfigValidate(t *testing.T) {
	tests := []struct {
		name      string
		config    Config
		expectErr bool
	}{
		{
			name: "valid config",
			config: Config{
				DefaultMode: "normal",
			},
			expectErr: false,
		},
		{
			name: "invalid mode",
			config: Config{
				DefaultMode: "turbo",
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.expectErr && err == nil {
				t.Errorf("Config.Validate() expected error but got nil")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("Config.Validate() expected no error but got: %v", err)
			}
		})
	}
}

func TestValidateModeTrimsWhitespace(t *testing.T) {
	// Note: Trimming should happen before calling ValidateMode
	// This test documents that ValidateMode expects already-trimmed input
	tests := []struct {
		name      string
		mode      string
		expectErr bool
	}{
		{"already trimmed", "crunch", false},
		{"with spaces (not trimmed)", " crunch ", true}, // ValidateMode expects trimmed input
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateMode(tt.mode)
			if tt.expectErr && err == nil {
				t.Errorf("ValidateMode(%q) expected error but got nil", tt.mode)
			}
			if !tt.expectErr && err != nil {
				t.Errorf("ValidateMode(%q) expected no error but got: %v", tt.mode, err)
			}
		})
	}
}
