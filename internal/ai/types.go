package ai

import (
	"fmt"
	"time"

	"github.com/Alvkoen/barely-incharge/internal/planner"
)

type PlanRequest struct {
	WorkStart  time.Time
	WorkEnd    time.Time
	BusyBlocks []planner.TimeBlock
	Tasks      []planner.Task
	Mode       string
}

type PlanResponse struct {
	Blocks []Block `json:"blocks"`
}

type Block struct {
	Type  string `json:"type"`
	Title string `json:"title"`
	Start string `json:"start"`
	End   string `json:"end"`
}

func (b Block) ToTimeBlock(date time.Time) (planner.TimeBlock, error) {
	startTime, err := planner.ParseTimeOnDate(b.Start, date)
	if err != nil {
		return planner.TimeBlock{}, fmt.Errorf("invalid start time: %w", err)
	}

	endTime, err := planner.ParseTimeOnDate(b.End, date)
	if err != nil {
		return planner.TimeBlock{}, fmt.Errorf("invalid end time: %w", err)
	}

	return planner.TimeBlock{
		Type:  b.Type,
		Title: b.Title,
		Start: startTime,
		End:   endTime,
	}, nil
}
