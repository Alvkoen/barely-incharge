package planner

import "time"

const (
	BlockTypeLunch   = "lunch"
	BlockTypeMeeting = "meeting"
	BlockTypeBreak   = "break"
	BlockTypeFocus   = "focus"

	// TimeFormat for HH:MM (24-hour)
	TimeFormat = "15:04"
)

type TimeBlock struct {
	Type  string
	Title string
	Start time.Time
	End   time.Time
}

func (b TimeBlock) GetCalendarTitle() string {
	switch b.Type {
	case BlockTypeFocus:
		return "Focus time"
	case BlockTypeBreak:
		return "Break"
	default:
		return b.Title
	}
}

func (b TimeBlock) GetCalendarDescription() string {
	switch b.Type {
	case BlockTypeFocus:
		return "Focus time block planned by Barely In Charge: " + b.Title
	case BlockTypeBreak:
		return "Break block planned by Barely In Charge: " + b.Title
	case BlockTypeLunch:
		return "Lunch break"
	default:
		return ""
	}
}

func ParseTimeOnDate(timeStr string, date time.Time) (time.Time, error) {
	t, err := time.Parse(TimeFormat, timeStr)
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
