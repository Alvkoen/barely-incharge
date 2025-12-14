package cmd

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Alvkoen/barely-incharge/internal/ai"
	"github.com/Alvkoen/barely-incharge/internal/calendar"
	"github.com/Alvkoen/barely-incharge/internal/config"
	"github.com/Alvkoen/barely-incharge/internal/planner"
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
		taskList := planner.ParseTaskList(tasks)

		fmt.Println("🎯 Planning your day...")
		fmt.Printf("Mode: %s\n", selectedMode)
		fmt.Printf("Work Hours: %s - %s\n", cfg.WorkHours.Start, cfg.WorkHours.End)
		fmt.Printf("Lunch Time: %s - %s\n", cfg.LunchTime.Start, cfg.LunchTime.End)
		fmt.Printf("Tasks (%d):\n", len(taskList))
		for i, task := range taskList {
			fmt.Printf("  %d. %s (%d min)\n", i+1, task.Title, int(task.Duration.Minutes()))
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

		plan, err := generatePlan(cfg, selectedMode, taskList, meetings)
		if err != nil {
			return err
		}

		fmt.Printf("\n✨ Generated %d blocks:\n", len(plan.Blocks))
		for i, block := range plan.Blocks {
			icon := "🎯"
			if block.Type == "break" {
				icon = "☕"
			}
			fmt.Printf("  %d. %s %s (%s - %s)\n", i+1, icon, block.Title, block.Start, block.End)
		}

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

func generatePlan(cfg *config.Config, mode string, tasks []planner.Task, meetings []calendar.Event) (*ai.PlanResponse, error) {
	now := time.Now()

	workStart, err := parseTimeToday(cfg.WorkHours.Start, now)
	if err != nil {
		return nil, fmt.Errorf("invalid work start time: %w", err)
	}

	workEnd, err := parseTimeToday(cfg.WorkHours.End, now)
	if err != nil {
		return nil, fmt.Errorf("invalid work end time: %w", err)
	}

	lunchStart, err := parseTimeToday(cfg.LunchTime.Start, now)
	if err != nil {
		return nil, fmt.Errorf("invalid lunch start time: %w", err)
	}

	lunchEnd, err := parseTimeToday(cfg.LunchTime.End, now)
	if err != nil {
		return nil, fmt.Errorf("invalid lunch end time: %w", err)
	}

	aiTasks := make([]ai.Task, len(tasks))
	for i, task := range tasks {
		aiTasks[i] = ai.Task{
			Title:    task.Title,
			Duration: task.Duration,
		}
	}

	busyBlocks := make([]ai.TimeBlock, 0, len(meetings)+1)

	busyBlocks = append(busyBlocks, ai.TimeBlock{
		Title: "Lunch",
		Start: lunchStart,
		End:   lunchEnd,
	})

	for _, meeting := range meetings {
		busyBlocks = append(busyBlocks, ai.TimeBlock{
			Title: meeting.Title,
			Start: meeting.Start,
			End:   meeting.End,
		})
	}

	fmt.Println("\n🤖 Generating plan with AI...")

	client := ai.NewClient(cfg.OpenAIAPIKey)
	req := ai.PlanRequest{
		WorkStart:  workStart,
		WorkEnd:    workEnd,
		BusyBlocks: busyBlocks,
		Tasks:      aiTasks,
		Mode:       mode,
	}

	return client.GeneratePlan(context.Background(), req)
}

func parseTimeToday(timeStr string, baseTime time.Time) (time.Time, error) {
	t, err := time.Parse("15:04", timeStr)
	if err != nil {
		return time.Time{}, err
	}

	return time.Date(
		baseTime.Year(),
		baseTime.Month(),
		baseTime.Day(),
		t.Hour(),
		t.Minute(),
		0, 0,
		baseTime.Location(),
	), nil
}
