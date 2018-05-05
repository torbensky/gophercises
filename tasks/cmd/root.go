package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/torbensky/gophercises/tasks/store"
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

var todos store.TodoService

func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath("$HOME/.tasks")
	viper.AddConfigPath(".")
	viper.SetDefault("database.path", "./tasks.db")
}

func Execute() {
	t, err := store.NewBolt(viper.GetString("database.path"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	todos = t

	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(doCmd)
	rootCmd.AddCommand(rmCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(completedCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
