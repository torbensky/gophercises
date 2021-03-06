package cmd

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var doCmd = &cobra.Command{
	Use:   "do",
	Short: "Complete a task",
	Long:  "Marks a task in the TODO list as completed at the current time.",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("one argument (task number) is required")
		}

		_, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid task number %s - use a number corresponding to one output by the task list command", args[0])
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		num, _ := strconv.Atoi(args[0])
		task, err := todos.CompleteTaskNum(num)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf(`You have completed the "%s" task.`, task.Description)
		fmt.Println()
	},
}
