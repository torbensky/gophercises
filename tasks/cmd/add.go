package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/torbensky/gophercises/tasks/store"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new task",
	Long:  `Adds a new task to the TODO list`,
	Run: func(cmd *cobra.Command, args []string) {
		desc := strings.Join(args, " ")
		todos.AddTask(&store.Task{Description: desc})
		fmt.Printf("Added \"%s\" to your task list.\n", desc)
	},
}
