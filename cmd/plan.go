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

		selectedMode := strings.TrimSpace(mode)
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
		taskList := parseTaskList(tasks)

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

		ctx := context.Background()
		calClient, err := authenticateCalendar(ctx)
		if err != nil {
			return err
		}

		meetings, err := fetchMeetings(calClient, cfg.MeetingsCalendar)
		if err != nil {
			return err
		}

		// TODO: Next steps:
		// 1. Call AI to generate blocks based on meetings + tasks + mode
		// 2. Create blocks in blocks_calendar

		// Prevent unused variable error
		_ = meetings

		return nil
	},
}

func authenticateCalendar(ctx context.Context) (*calendar.GoogleClient, error) {
	fmt.Println("\n🔐 Authenticating with Google Calendar...")

	client, err := calendar.NewGoogleClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to authenticate with Google Calendar: %w", err)
	}

	fmt.Println("✅ Successfully authenticated!")

	return client, nil
}

func fetchMeetings(client *calendar.GoogleClient, calendarID string) ([]calendar.Event, error) {
	fmt.Printf("\n📆 Fetching meetings from calendar: %s\n", calendarID)

	meetings, err := client.FetchTodaysMeetings(calendarID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch meetings: %w", err)
	}

	if len(meetings) == 0 {
		fmt.Println("  No meetings found for today")
	} else {
		fmt.Printf("  Found %d meeting(s):\n", len(meetings))
		for i, meeting := range meetings {
			fmt.Printf("  %d. %s (%s - %s)\n",
				i+1,
				meeting.Title,
				meeting.Start.Format("15:04"),
				meeting.End.Format("15:04"))
		}
	}

	return meetings, nil
}

func init() {
	rootCmd.AddCommand(planCmd)
	planCmd.Flags().StringVarP(&tasks, "tasks", "t", "", "Comma-separated list of tasks to accomplish (required)")
	planCmd.Flags().StringVarP(&mode, "mode", "m", "", "Planning mode: crunch, normal, or saver (default from config)")
	if err := planCmd.MarkFlagRequired("tasks"); err != nil {
		panic(err)
	}
}

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
