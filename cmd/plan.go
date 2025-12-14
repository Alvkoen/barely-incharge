package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/Alvkoen/barely-incharge/internal/calendar"
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

		// Authenticate with Google Calendar
		ctx := context.Background()
		if err := authenticateCalendar(ctx); err != nil {
			return err
		}

		// TODO: Next steps:
		// 1. Fetch meetings from meetings_calendar
		// 2. Call AI to generate blocks
		// 3. Create blocks in blocks_calendar

		return nil
	},
}

// authenticateCalendar authenticates with Google Calendar and returns the service
func authenticateCalendar(ctx context.Context) error {
	fmt.Println("\n🔐 Authenticating with Google Calendar...")

	calService, err := calendar.GetClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to authenticate with Google Calendar: %w", err)
	}

	fmt.Println("✅ Successfully authenticated!")

	// Test: List available calendars
	calendarList, err := calService.CalendarList.List().Do()
	if err != nil {
		return fmt.Errorf("failed to list calendars: %w", err)
	}

	fmt.Printf("\n📅 Available calendars (%d):\n", len(calendarList.Items))
	for _, cal := range calendarList.Items {
		fmt.Printf("  - %s (ID: %s)\n", cal.Summary, cal.Id)
	}

	return nil
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
