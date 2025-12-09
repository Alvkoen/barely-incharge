package cmd

import (
	// "fmt"
	// "os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "barelyincharge",
	Short: "AI-powered calendar block planner",
	Long: `Barely In Charge is a CLI tool that uses AI to intelligently 
schedule focus blocks and breaks in your Google Calendar based on 
your meetings and tasks for the day.`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Persistent flags and configuration will go here
}
