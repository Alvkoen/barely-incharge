package ai

import (
	"fmt"
	"strings"
)

func BuildPrompt(req PlanRequest) string {
	var sb strings.Builder

	sb.WriteString("You are a calendar planning assistant. Create a day schedule with focus blocks and breaks.\n\n")

	sb.WriteString(fmt.Sprintf("Work hours: %s - %s\n",
		req.WorkStart.Format("15:04"),
		req.WorkEnd.Format("15:04")))

	sb.WriteString("\n")

	sb.WriteString("Busy times (unavailable for scheduling):\n")
	if len(req.BusyBlocks) == 0 {
		sb.WriteString("- No busy times\n")
	} else {
		for _, block := range req.BusyBlocks {
			sb.WriteString(fmt.Sprintf("- %s (%s - %s)\n",
				block.Title,
				block.Start.Format("15:04"),
				block.End.Format("15:04")))
		}
	}
	sb.WriteString("\n")

	sb.WriteString("Tasks to schedule:\n")
	for _, task := range req.Tasks {
		sb.WriteString(fmt.Sprintf("- %s (%d minutes)\n",
			task.Title,
			int(task.Duration.Minutes())))
	}
	sb.WriteString("\n")

	sb.WriteString(getModeInstructions(req.Mode))
	sb.WriteString("\n")

	sb.WriteString("IMPORTANT: Return ONLY valid JSON in this exact format with no additional text:\n")
	sb.WriteString("{\n")
	sb.WriteString("  \"blocks\": [\n")
	sb.WriteString("    {\"type\": \"focus\", \"title\": \"Task name\", \"start\": \"HH:MM\", \"end\": \"HH:MM\"},\n")
	sb.WriteString("    {\"type\": \"break\", \"title\": \"Short break\", \"start\": \"HH:MM\", \"end\": \"HH:MM\"}\n")
	sb.WriteString("  ]\n")
	sb.WriteString("}\n\n")

	sb.WriteString("Rules:\n")
	sb.WriteString("- Do NOT overlap with busy times\n")
	sb.WriteString("- Stay within work hours\n")
	sb.WriteString("- Use 24-hour format (HH:MM)\n")
	sb.WriteString("- Types: \"focus\" for tasks, \"break\" for breaks\n")
	sb.WriteString("- Return ONLY the JSON, no explanation or markdown\n")

	return sb.String()
}

func getModeInstructions(mode string) string {
	switch mode {
	case "crunch":
		return "Mode: CRUNCH - Pack as many tasks as possible with minimal breaks (5-10 min). Maximize productivity."
	case "saver":
		return "Mode: ENERGY SAVER - User is tired. Add longer breaks (15-20 min), extra padding between tasks, and consider ending early if possible."
	default:
		return "Mode: NORMAL - Balanced approach with regular breaks (10-15 min) following standard productivity practices."
	}
}
