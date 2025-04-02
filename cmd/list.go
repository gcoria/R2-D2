package cmd

import (
	"R2-D2/todo"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tasks",
	Run: func(cmd *cobra.Command, args []string) {
		tasks, err := todo.LoadTasks("tasks.csv")
		if err != nil {
			fmt.Println("Error loading tasks:", err)
			return
		}

		if len(tasks) == 0 {
			fmt.Println("No tasks to display")
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.Debug)

		fmt.Fprintln(w, "ID\tSTATUS\tDESCRIPTION\tCREATED AT")

		for _, task := range tasks {
			status := "Pending"
			if task.Completed {
				status = "Complete"
			}

			createdTime := task.CreatedAt.Format("2006-01-02 15:04:05")

			fmt.Fprintf(w, "%d\t%s\t%s\t%s\n",
				task.ID,
				status,
				task.Description,
				createdTime,
			)
		}

		w.Flush()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
