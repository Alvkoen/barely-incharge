package calendar

import (
	"time"

	"github.com/Alvkoen/barely-incharge/internal/planner"
)

type Event struct {
	Type        string
	Title       string
	Description string
	Start       time.Time
	End         time.Time
}

func (e Event) ToTimeBlock() planner.TimeBlock {
	return planner.TimeBlock{
		Type:  e.Type,
		Title: e.Title,
		Start: e.Start,
		End:   e.End,
	}
}
