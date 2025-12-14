package planner

import (
	"reflect"
	"testing"
	"time"
)

func TestParseTaskList(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Task
	}{
		{
			name:  "single task",
			input: "Write documentation",
			expected: []Task{
				{Title: "Write documentation", Duration: SizeM},
			},
		},
		{
			name:  "multiple tasks",
			input: "Write docs,Review PRs,Deploy code",
			expected: []Task{
				{Title: "Write docs", Duration: SizeM},
				{Title: "Review PRs", Duration: SizeM},
				{Title: "Deploy code", Duration: SizeM},
			},
		},
		{
			name:  "tasks with sizes",
			input: "Write docs:L,Review PRs:S,Quick fix:XS",
			expected: []Task{
				{Title: "Write docs", Duration: SizeL},
				{Title: "Review PRs", Duration: SizeS},
				{Title: "Quick fix", Duration: SizeXS},
			},
		},
		{
			name:  "tasks with lowercase sizes",
			input: "Write docs:l,Review PRs:s",
			expected: []Task{
				{Title: "Write docs", Duration: SizeL},
				{Title: "Review PRs", Duration: SizeS},
			},
		},
		{
			name:  "tasks with mixed sizes and no sizes",
			input: "Write docs:XL,Review PRs,Quick fix:XS",
			expected: []Task{
				{Title: "Write docs", Duration: SizeXL},
				{Title: "Review PRs", Duration: SizeM},
				{Title: "Quick fix", Duration: SizeXS},
			},
		},
		{
			name:  "task with invalid size defaults to M",
			input: "Write docs:INVALID",
			expected: []Task{
				{Title: "Write docs:INVALID", Duration: SizeM},
			},
		},
		{
			name:     "empty string",
			input:    "",
			expected: []Task{},
		},
		{
			name:     "only commas",
			input:    ",,,",
			expected: []Task{},
		},
		{
			name:  "mixed empty and valid",
			input: "Write docs:L,,Review PRs:S,",
			expected: []Task{
				{Title: "Write docs", Duration: SizeL},
				{Title: "Review PRs", Duration: SizeS},
			},
		},
		{
			name:  "tasks with extra whitespace",
			input: "  Write docs:L  ,  Review PRs:S  ",
			expected: []Task{
				{Title: "Write docs", Duration: SizeL},
				{Title: "Review PRs", Duration: SizeS},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseTaskList(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ParseTaskList(%q) = %v, want %v", tt.input, result, tt.expected)
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
		{"tasks with sizes", "Task1:S,Task2:L,Task3:XL", 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseTaskList(tt.input)
			if len(result) != tt.expectedCount {
				t.Errorf("ParseTaskList(%q) returned %d tasks, want %d", tt.input, len(result), tt.expectedCount)
			}
		})
	}
}

func TestParseTaskListPreservesContent(t *testing.T) {
	input := "Write code:L,Test with numbers 12345:S"
	result := ParseTaskList(input)

	if len(result) != 2 {
		t.Fatalf("Expected 2 tasks, got %d", len(result))
	}

	if result[0].Title != "Write code" || result[0].Duration != SizeL {
		t.Errorf("First task not parsed correctly: got %v", result[0])
	}

	if result[1].Title != "Test with numbers 12345" || result[1].Duration != SizeS {
		t.Errorf("Second task not parsed correctly: got %v", result[1])
	}
}

func TestParseTaskSizes(t *testing.T) {
	tests := []struct {
		name             string
		input            string
		expectedTitle    string
		expectedDuration time.Duration
	}{
		{"XS size", "Quick task:XS", "Quick task", SizeXS},
		{"S size", "Small task:S", "Small task", SizeS},
		{"M size", "Medium task:M", "Medium task", SizeM},
		{"L size", "Large task:L", "Large task", SizeL},
		{"XL size", "Extra large task:XL", "Extra large task", SizeXL},
		{"no size defaults to M", "Default task", "Default task", SizeM},
		{"lowercase size", "Task:xl", "Task", SizeXL},
		{"mixed case size", "Task:Xl", "Task", SizeXL},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseTaskList(tt.input)
			if len(result) != 1 {
				t.Fatalf("Expected 1 task, got %d", len(result))
			}
			if result[0].Title != tt.expectedTitle {
				t.Errorf("Title = %q, want %q", result[0].Title, tt.expectedTitle)
			}
			if result[0].Duration != tt.expectedDuration {
				t.Errorf("Duration = %v, want %v", result[0].Duration, tt.expectedDuration)
			}
		})
	}
}
