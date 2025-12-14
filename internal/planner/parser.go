package planner

import (
	"strings"
	"time"
)

const (
	SizeXS = 10 * time.Minute
	SizeS  = 15 * time.Minute
	SizeM  = 30 * time.Minute
	SizeL  = 60 * time.Minute
	SizeXL = 90 * time.Minute
)

var taskSizes = map[string]time.Duration{
	"XS": SizeXS,
	"S":  SizeS,
	"M":  SizeM,
	"L":  SizeL,
	"XL": SizeXL,
}

func ParseTaskList(tasksStr string) []Task {
	parts := strings.Split(tasksStr, ",")
	tasks := make([]Task, 0, len(parts))

	for _, taskStr := range parts {
		trimmed := strings.TrimSpace(taskStr)
		if trimmed == "" {
			continue
		}

		title, duration := parseTask(trimmed)
		tasks = append(tasks, Task{
			Title:    title,
			Duration: duration,
		})
	}

	return tasks
}

func parseTask(taskStr string) (string, time.Duration) {
	parts := strings.Split(taskStr, ":")
	if len(parts) == 1 {
		return strings.TrimSpace(parts[0]), SizeM
	}

	title := strings.TrimSpace(parts[0])
	sizeStr := strings.TrimSpace(strings.ToUpper(parts[1]))

	if duration, ok := taskSizes[sizeStr]; ok {
		return title, duration
	}

	return taskStr, SizeM
}
