package calendar

import "time"

// Event type constants
const (
	EventTypeMeeting = "meeting"
	EventTypeFocus   = "focus"
	EventTypeBreak   = "break"
)

// Event represents any calendar event (meeting, focus block, or break)
type Event struct {
	Type        string // EventTypeMeeting, EventTypeFocus, or EventTypeBreak
	Title       string
	Description string
	Start       time.Time
	End         time.Time
}
