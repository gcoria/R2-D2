package cmd

import (
	"R2-D2/todo"
	"fmt"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

var completeCmd = &cobra.Command{
	Use:   "complete [task ID]",
	Short: "Complete a task",
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
				tasks[i].Completed = true
				tasks[i].CompletedAt = time.Now()
				err = todo.SaveTasks("tasks.csv", tasks)
				if err != nil {
					fmt.Println("Error saving tasks:", err)
				}
			}
		}
		fmt.Printf("Task %d completed\n", id)
	},
}

func init() {
	rootCmd.AddCommand(completeCmd)
}
