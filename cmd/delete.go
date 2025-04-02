package cmd

import (
	"R2-D2/todo"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete [task ID]",
	Short: "Delete a task",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Error converting task ID to int:", err)
			return
		}
		tasks, err := todo.LoadTasks("tasks.csv")
		if err != nil {
			fmt.Println("Error loading tasks:", err)
			return
		}
		for i, task := range tasks {
			if task.ID == id {
				tasks = append(tasks[:i], tasks[i+1:]...)
				err = todo.SaveTasks("tasks.csv", tasks)
				if err != nil {
					fmt.Println("Error saving tasks:", err)
					return
				}
				fmt.Printf("Task %d deleted\n", id)
				return
			}
		}
		fmt.Printf("Task %d not found\n", id)
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
