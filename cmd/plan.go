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

		planningDate, err := cfg.GetPlanningDate()
		if err != nil {
			return fmt.Errorf("failed to parse planning date: %w", err)
		}

		fmt.Println("🎯 Planning your day...")
		fmt.Printf("Date: %s\n", planningDate.Format("Monday, January 2, 2006"))
		fmt.Printf("Mode: %s\n", selectedMode)
		fmt.Printf("Work Hours: %s - %s\n", cfg.WorkHours.Start, cfg.WorkHours.End)
		fmt.Printf("Lunch Time: %s - %s\n", cfg.LunchTime.Start, cfg.LunchTime.End)
		fmt.Printf("Tasks (%d):\n", len(taskList))
		for i, task := range taskList {
			fmt.Printf("  %d. %s (%d min)\n", i+1, task.Title, int(task.Duration.Minutes()))
		}
		fmt.Printf("\nCalendar: %s\n", cfg.Calendar)

		ctx := context.Background()

		calClient, err := authenticateCalendar(ctx)
		if err != nil {
			return err
		}

		meetings, err := fetchMeetings(calClient, cfg.Calendar)
		if err != nil {
			return err
		}

		workStart, err := planner.ParseTimeOnDate(cfg.WorkHours.Start, planningDate)
		if err != nil {
			return fmt.Errorf("invalid work start time: %w", err)
		}
		workEnd, err := planner.ParseTimeOnDate(cfg.WorkHours.End, planningDate)
		if err != nil {
			return fmt.Errorf("invalid work end time: %w", err)
		}
		lunchStart, err := planner.ParseTimeOnDate(cfg.LunchTime.Start, planningDate)
		if err != nil {
			return fmt.Errorf("invalid lunch start time: %w", err)
		}
		lunchEnd, err := planner.ParseTimeOnDate(cfg.LunchTime.End, planningDate)
		if err != nil {
			return fmt.Errorf("invalid lunch end time: %w", err)
		}

		busyBlocks := make([]planner.TimeBlock, 0, len(meetings)+1)
		busyBlocks = append(busyBlocks, planner.TimeBlock{
			Type:  planner.BlockTypeLunch,
			Title: "Lunch",
			Start: lunchStart,
			End:   lunchEnd,
		})
		for _, meeting := range meetings {
			busyBlocks = append(busyBlocks, meeting.ToTimeBlock())
		}

		plan, err := generatePlan(cfg, selectedMode, workStart, workEnd, taskList, busyBlocks)
		if err != nil {
			return err
		}

		fmt.Printf("\n✨ Generated %d blocks:\n", len(plan.Blocks))
		for i, block := range plan.Blocks {
			icon := "🎯"
			if block.Type == planner.BlockTypeBreak {
				icon = "☕"
			}
			fmt.Printf("  %d. %s %s (%s - %s)\n", i+1, icon, block.Title, block.Start, block.End)
		}

		fmt.Println("\n🔍 Validating schedule...")
		parsedBlocks, err := parseAIBlocks(plan.Blocks, planningDate)
		if err != nil {
			return fmt.Errorf("failed to parse AI blocks: %w", err)
		}

		err = planner.ValidateBlocks(parsedBlocks, busyBlocks)
		if err != nil {
			return fmt.Errorf("schedule validation failed: %w", err)
		}
		fmt.Println("✓ No conflicts detected")

		err = createBlocks(calClient, cfg.Calendar, parsedBlocks)
		if err != nil {
			return err
		}

		fmt.Println("\n✅ Successfully created all blocks in calendar!")

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
				meeting.Start.Format(planner.TimeFormat),
				meeting.End.Format(planner.TimeFormat))
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

func generatePlan(cfg *config.Config, mode string, workStart, workEnd time.Time, tasks []planner.Task, busyBlocks []planner.TimeBlock) (*ai.PlanResponse, error) {
	fmt.Println("\n🤖 Generating plan with AI...")

	client := ai.NewClient(cfg.OpenAIAPIKey)
	req := ai.PlanRequest{
		WorkStart:  workStart,
		WorkEnd:    workEnd,
		BusyBlocks: busyBlocks,
		Tasks:      tasks,
		Mode:       mode,
	}

	return client.GeneratePlan(context.Background(), req)
}

func parseAIBlocks(aiBlocks []ai.Block, date time.Time) ([]planner.TimeBlock, error) {
	blocks := make([]planner.TimeBlock, len(aiBlocks))

	for i, block := range aiBlocks {
		timeBlock, err := block.ToTimeBlock(date)
		if err != nil {
			return nil, fmt.Errorf("invalid block %d: %w", i+1, err)
		}
		blocks[i] = timeBlock
	}

	return blocks, nil
}

func createBlocks(client *calendar.GoogleClient, calendarID string, parsedBlocks []planner.TimeBlock) error {
	fmt.Println("\n📝 Creating blocks in calendar...")

	for _, block := range parsedBlocks {
		event := calendar.Event{
			Type:        block.Type,
			Title:       block.GetCalendarTitle(),
			Description: block.GetCalendarDescription(),
			Start:       block.Start,
			End:         block.End,
		}

		if err := client.CreateEvent(calendarID, event); err != nil {
			return fmt.Errorf("failed to create block '%s': %w", event.Title, err)
		}

		fmt.Printf("  ✓ Created: %s (%s)\n", event.Title, block.Title)
	}

	return nil
}
