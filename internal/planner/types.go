package planner

import "time"

const (
	BlockTypeLunch   = "lunch"
	BlockTypeMeeting = "meeting"
	BlockTypeBreak   = "break"
	BlockTypeFocus   = "focus"
)

type Task struct {
	Title    string
	Duration time.Duration
}

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
