package planner

import (
	"fmt"
	"time"
)

type Context struct {
	Mode       string
	WorkStart  time.Time
	WorkEnd    time.Time
	LunchStart time.Time
	LunchEnd   time.Time
	Tasks      []Task
	BusyBlocks []TimeBlock
}

func NewContext(
	mode string,
	workStartStr, workEndStr string,
	lunchStartStr, lunchEndStr string,
	tasks []Task,
	busyBlocks []TimeBlock,
) (*Context, error) {
	now := time.Now()

	workStart, err := ParseTimeOnDate(workStartStr, now)
	if err != nil {
		return nil, fmt.Errorf("invalid work start time: %w", err)
	}

	workEnd, err := ParseTimeOnDate(workEndStr, now)
	if err != nil {
		return nil, fmt.Errorf("invalid work end time: %w", err)
	}

	lunchStart, err := ParseTimeOnDate(lunchStartStr, now)
	if err != nil {
		return nil, fmt.Errorf("invalid lunch start time: %w", err)
	}

	lunchEnd, err := ParseTimeOnDate(lunchEndStr, now)
	if err != nil {
		return nil, fmt.Errorf("invalid lunch end time: %w", err)
	}

	allBusyBlocks := make([]TimeBlock, 0, len(busyBlocks)+1)
	allBusyBlocks = append(allBusyBlocks, TimeBlock{
		Type:  BlockTypeLunch,
		Title: "Lunch",
		Start: lunchStart,
		End:   lunchEnd,
	})
	allBusyBlocks = append(allBusyBlocks, busyBlocks...)

	return &Context{
		Mode:       mode,
		WorkStart:  workStart,
		WorkEnd:    workEnd,
		LunchStart: lunchStart,
		LunchEnd:   lunchEnd,
		Tasks:      tasks,
		BusyBlocks: allBusyBlocks,
	}, nil
}

// ParseTimeOnDate converts a "HH:MM" time string to a full timestamp on the given date.
// The date parameter ensures all times in a single operation use the same reference date,
// which is important for consistency and testability.
func ParseTimeOnDate(timeStr string, date time.Time) (time.Time, error) {
	t, err := time.Parse("15:04", timeStr)
	if err != nil {
		return time.Time{}, err
	}

	return time.Date(
		date.Year(),
		date.Month(),
		date.Day(),
		t.Hour(),
		t.Minute(),
		0, 0,
		date.Location(),
	), nil
}
