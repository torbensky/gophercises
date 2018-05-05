package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var completedCmd = &cobra.Command{
	Use:   "list",
	Short: "List TODO tasks",
	Long:  "Prints a list of the remaining, unfinished TODO tasks",
	Run: func(cmd *cobra.Command, args []string) {
		tasks := todos.ListUncompleteTasks()
		if len(tasks) == 0 {
			fmt.Println("You have no tasks to do.")
			return
		}
		fmt.Println("You have the following tasks:")
		for i, todo := range tasks {
			fmt.Printf("%d. %s\n", i+1, todo.Description)
		}
	},
}
