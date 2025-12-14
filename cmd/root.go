package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "barely-incharge",
	Short: "AI-powered calendar block planner",
	Long: `Barely In Charge is a CLI tool that uses AI to intelligently 
schedule focus blocks and breaks in your Google Calendar based on 
your meetings and tasks for the day.`,
}

func Execute() error {
	return rootCmd.Execute()
}
