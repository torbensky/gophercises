package cmd

import "github.com/spf13/cobra"

var doCmd = &cobra.Command{
	Use:   "do",
	Short: "Complete a task",
	Long:  "Removes a task from the TODO list",
}
