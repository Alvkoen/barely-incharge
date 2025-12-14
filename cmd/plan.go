package cmd

import (
	"fmt"
	"strings"

	"github.com/Alvkoen/barely-incharge/internal/config"
	"github.com/spf13/cobra"
)

var (
	tasks string
	mode  string
)

var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "Plan your day with AI-powered focus blocks",
	Long:  `Create focus and break blocks in your calendar based on your tasks, meetings, and chosen mode (crunch, normal, or saver).`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Load config
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		selectedMode := mode
		if selectedMode == "" {
			selectedMode = cfg.DefaultMode 
		} else {
			if err := config.ValidateMode(selectedMode); err != nil {
				return err
			}
		}

		if tasks == "" {
			return fmt.Errorf("--tasks flag is required")
		}

		// Parse tasks into a list
		taskList := parseTaskList(tasks)

		// Display what we're working with
		fmt.Println("🎯 Planning your day...")
		fmt.Printf("Mode: %s\n", selectedMode)
		fmt.Printf("Work Hours: %s - %s\n", cfg.WorkHours.Start, cfg.WorkHours.End)
		fmt.Printf("Lunch Time: %s - %s\n", cfg.LunchTime.Start, cfg.LunchTime.End)
		fmt.Printf("Tasks (%d):\n", len(taskList))
		for i, task := range taskList {
			fmt.Printf("  %d. %s\n", i+1, task)
		}
		fmt.Printf("\nMeetings Calendar: %s\n", cfg.MeetingsCalendar)
		fmt.Printf("Blocks Calendar: %s\n", cfg.BlocksCalendar)

		// TODO: Next steps:
		// 1. Authenticate with Google Calendar
		// 2. Fetch meetings
		// 3. Call AI to generate blocks
		// 4. Create blocks in calendar

		return nil
	},
}

func init() {
	rootCmd.AddCommand(planCmd)

	// Define flags
	planCmd.Flags().StringVarP(&tasks, "tasks", "t", "", "Comma-separated list of tasks to accomplish (required)")
	planCmd.Flags().StringVarP(&mode, "mode", "m", "", "Planning mode: crunch, normal, or saver (default from config)")

	// Mark tasks as required
	planCmd.MarkFlagRequired("tasks")
}

// parseTaskList splits the tasks string by comma and trims whitespace
func parseTaskList(tasksStr string) []string {
	parts := strings.Split(tasksStr, ",")
	tasks := make([]string, 0, len(parts))

	for _, task := range parts {
		trimmed := strings.TrimSpace(task)
		if trimmed != "" {
			tasks = append(tasks, trimmed)
		}
	}

	return tasks
}
