package cmd

import (
	"reflect"
	"testing"
)

func TestParseTaskList(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "single task",
			input:    "Write documentation",
			expected: []string{"Write documentation"},
		},
		{
			name:     "multiple tasks",
			input:    "Write docs,Review PRs,Deploy code",
			expected: []string{"Write docs", "Review PRs", "Deploy code"},
		},
		{
			name:     "tasks with extra whitespace",
			input:    "  Write docs  ,  Review PRs  ,  Deploy code  ",
			expected: []string{"Write docs", "Review PRs", "Deploy code"},
		},
		{
			name:     "empty string",
			input:    "",
			expected: []string{},
		},
		{
			name:     "only commas",
			input:    ",,,",
			expected: []string{},
		},
		{
			name:     "mixed empty and valid",
			input:    "Write docs,,Review PRs,",
			expected: []string{"Write docs", "Review PRs"},
		},
		{
			name:     "task with newlines",
			input:    "Write docs,\nReview PRs",
			expected: []string{"Write docs", "Review PRs"},
		},
		{
			name:     "task with tabs",
			input:    "Write docs,\tReview PRs",
			expected: []string{"Write docs", "Review PRs"},
		},
		{
			name:     "single task no comma",
			input:    "Single task here",
			expected: []string{"Single task here"},
		},
		{
			name:     "tasks with special characters",
			input:    "Fix bug #123,Update README.md,Send email to john@example.com",
			expected: []string{"Fix bug #123", "Update README.md", "Send email to john@example.com"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseTaskList(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("parseTaskList(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestParseTaskListLength(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedCount int
	}{
		{"zero tasks", "", 0},
		{"one task", "Task", 1},
		{"three tasks", "Task1,Task2,Task3", 3},
		{"five tasks with spaces", "T1, T2, T3, T4, T5", 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseTaskList(tt.input)
			if len(result) != tt.expectedCount {
				t.Errorf("parseTaskList(%q) returned %d tasks, want %d", tt.input, len(result), tt.expectedCount)
			}
		})
	}
}

func TestParseTaskListPreservesContent(t *testing.T) {
	input := "Write code with special chars: !@#$%,Test with numbers 12345"
	result := parseTaskList(input)

	if len(result) != 2 {
		t.Fatalf("Expected 2 tasks, got %d", len(result))
	}

	if result[0] != "Write code with special chars: !@#$%" {
		t.Errorf("First task not preserved correctly: got %q", result[0])
	}

	if result[1] != "Test with numbers 12345" {
		t.Errorf("Second task not preserved correctly: got %q", result[1])
	}
}

