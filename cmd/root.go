package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "R2-D2",
	Short: "A simple CLI for managing your todo list",
}

func Execute() error {
	return rootCmd.Execute()
}
