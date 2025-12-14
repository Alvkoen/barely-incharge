package calendar

import "time"

// Event represents any calendar event (meeting, focus block, or break)
// Type should be one of: planner.BlockTypeMeeting, BlockTypeFocus, BlockTypeBreak, or BlockTypeLunch
type Event struct {
	Type        string
	Title       string
	Description string
	Start       time.Time
	End         time.Time
}
