package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "completed",
	Short: "List completed TODO tasks from today",
	Long:  "Prints a list of TODO tasks that have been completed today",
	Run: func(cmd *cobra.Command, args []string) {
		t := time.Now()
		startOfDay := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
		tasks := todos.ListTasksCompletedAfter(startOfDay)
		if len(tasks) == 0 {
			fmt.Println("You have no completed tasks today.")
			return
		}
		fmt.Println("You have finished the following tasks today:")
		for i, todo := range tasks {
			fmt.Printf("%d. %s\n", i+1, todo.Description)
		}
	},
}
