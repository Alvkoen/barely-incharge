package planner

import (
	"strings"
	"testing"
	"time"
)

func TestValidateBlocks_NoOverlaps(t *testing.T) {
	blocks := []TimeBlock{
		{Title: "Task 1", Start: time9AM(), End: time10AM()},
		{Title: "Task 2", Start: time10AM(), End: time11AM()},
	}

	busyBlocks := []TimeBlock{
		{Title: "Lunch", Start: time12PM(), End: time1PM()},
	}

	err := ValidateBlocks(blocks, busyBlocks)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

func TestValidateBlocks_SelfOverlap(t *testing.T) {
	blocks := []TimeBlock{
		{Title: "Task 1", Start: time9AM(), End: time10AM()},
		{Title: "Task 2", Start: time930AM(), End: time11AM()},
	}

	err := ValidateBlocks(blocks, nil)
	if err == nil {
		t.Error("Expected overlap error, got nil")
	}
	if !strings.Contains(err.Error(), "blocks overlap") {
		t.Errorf("Expected 'blocks overlap' in error, got: %v", err)
	}
}

func TestValidateBlocks_BusyOverlap(t *testing.T) {
	blocks := []TimeBlock{
		{Title: "Task 1", Start: time9AM(), End: time10AM()},
		{Title: "Task 2", Start: time12PM(), End: time1PM()},
	}

	busyBlocks := []TimeBlock{
		{Title: "Lunch", Start: time1230PM(), End: time1PM()},
	}

	err := ValidateBlocks(blocks, busyBlocks)
	if err == nil {
		t.Error("Expected overlap error, got nil")
	}
	if !strings.Contains(err.Error(), "overlaps with busy time") {
		t.Errorf("Expected 'overlaps with busy time' in error, got: %v", err)
	}
}

func TestValidateBlocks_AdjacentBlocks(t *testing.T) {
	blocks := []TimeBlock{
		{Title: "Task 1", Start: time9AM(), End: time10AM()},
		{Title: "Task 2", Start: time10AM(), End: time11AM()},
	}

	err := ValidateBlocks(blocks, nil)
	if err != nil {
		t.Errorf("Adjacent blocks should not overlap, got: %v", err)
	}
}

func TestBlocksOverlap(t *testing.T) {
	tests := []struct {
		name     string
		a        TimeBlock
		b        TimeBlock
		expected bool
	}{
		{
			name:     "complete overlap",
			a:        TimeBlock{Start: time9AM(), End: time11AM()},
			b:        TimeBlock{Start: time10AM(), End: time12PM()},
			expected: true,
		},
		{
			name:     "a contains b",
			a:        TimeBlock{Start: time9AM(), End: time12PM()},
			b:        TimeBlock{Start: time10AM(), End: time11AM()},
			expected: true,
		},
		{
			name:     "b contains a",
			a:        TimeBlock{Start: time10AM(), End: time11AM()},
			b:        TimeBlock{Start: time9AM(), End: time12PM()},
			expected: true,
		},
		{
			name:     "no overlap - before",
			a:        TimeBlock{Start: time9AM(), End: time10AM()},
			b:        TimeBlock{Start: time11AM(), End: time12PM()},
			expected: false,
		},
		{
			name:     "no overlap - after",
			a:        TimeBlock{Start: time11AM(), End: time12PM()},
			b:        TimeBlock{Start: time9AM(), End: time10AM()},
			expected: false,
		},
		{
			name:     "adjacent - a before b",
			a:        TimeBlock{Start: time9AM(), End: time10AM()},
			b:        TimeBlock{Start: time10AM(), End: time11AM()},
			expected: false,
		},
		{
			name:     "adjacent - b before a",
			a:        TimeBlock{Start: time10AM(), End: time11AM()},
			b:        TimeBlock{Start: time9AM(), End: time10AM()},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := blocksOverlap(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("blocksOverlap() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func time9AM() time.Time {
	return time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC)
}

func time930AM() time.Time {
	return time.Date(2024, 1, 1, 9, 30, 0, 0, time.UTC)
}

func time10AM() time.Time {
	return time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
}

func time11AM() time.Time {
	return time.Date(2024, 1, 1, 11, 0, 0, 0, time.UTC)
}

func time12PM() time.Time {
	return time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
}

func time1230PM() time.Time {
	return time.Date(2024, 1, 1, 12, 30, 0, 0, time.UTC)
}

func time1PM() time.Time {
	return time.Date(2024, 1, 1, 13, 0, 0, 0, time.UTC)
}
