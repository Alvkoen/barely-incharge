package ai

import "time"

type PlanRequest struct {
	WorkStart  time.Time
	WorkEnd    time.Time
	BusyBlocks []TimeBlock
	Tasks      []Task
	Mode       string
}

type Task struct {
	Title    string
	Duration time.Duration
}

type TimeBlock struct {
	Title string
	Start time.Time
	End   time.Time
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
