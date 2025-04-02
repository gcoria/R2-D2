package cmd

import (
	"R2-D2/todo"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new task",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Join all arguments to form the complete task description
		description := strings.Join(args, " ")

		tasks, err := todo.LoadTasks("tasks.csv")
		if err != nil {
			fmt.Println("Error loading tasks:", err)
			return
		}
		task := todo.Task{
			ID:          len(tasks) + 1,
			Description: description,
			Completed:   false,
			CreatedAt:   time.Now(),
			CompletedAt: time.Time{},
		}
		tasks = append(tasks, task)
		err = todo.SaveTasks("tasks.csv", tasks)
		if err != nil {
			fmt.Println("Error saving tasks:", err)
			return
		}
		fmt.Printf("Task added: %d - %s\n", task.ID, task.Description)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
