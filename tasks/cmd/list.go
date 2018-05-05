package cmd

import "github.com/spf13/cobra"

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List TODO tasks",
	Long:  "Prints a list of the remaining, unfinished TODO tasks",
}
