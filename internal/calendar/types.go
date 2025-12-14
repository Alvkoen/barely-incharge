package calendar

import (
	"time"

	"github.com/Alvkoen/barely-incharge/internal/planner"
)

// Event represents any calendar event (meeting, focus block, or break)
// Type should be one of: planner.BlockTypeMeeting, BlockTypeFocus, BlockTypeBreak, or BlockTypeLunch
type Event struct {
	Type        string
	Title       string
	Description string
	Start       time.Time
	End         time.Time
}

// ToTimeBlock converts a calendar event to a planner TimeBlock
func (e Event) ToTimeBlock() planner.TimeBlock {
	return planner.TimeBlock{
		Type:  e.Type,
		Title: e.Title,
		Start: e.Start,
		End:   e.End,
	}
}
