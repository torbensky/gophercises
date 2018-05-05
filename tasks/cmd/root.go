package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

/*

- TODO: accept a flag for a config path

- TODO: dbpath can be configured multiple ways
  - accept a flag for the db path
  - read config for the db path

*/

var rootCmd = &cobra.Command{
	Use:   "task",
	Short: "task is a task CLI",
	Long:  `The best darn task CLI, trust me, I know`,
}

func Execute() {
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(doCmd)
	rootCmd.AddCommand(listCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
