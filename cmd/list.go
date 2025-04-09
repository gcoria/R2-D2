package cmd

import (
	"R2-D2/todo"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var showSecretsFlag bool

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

		fmt.Fprintln(w, "ID\tSTATUS\tDESCRIPTION\tCREATED AT\tSECRET")

		for _, task := range tasks {
			status := "Pending"
			if task.Completed {
				status = "Complete"
			}

			createdTime := task.CreatedAt.Format("2006-01-02 15:04:05")

			description := task.Description
			secretStatus := "No"

			if task.Encrypted {
				secretStatus = "Yes"
				if showSecretsFlag {
					// Decrypt the task description
					decrypted, err := todo.DecryptText(task.Description)
					if err != nil {
						description = "[DECRYPT ERROR]"
					} else {
						description = decrypted
					}
				} else {
					description = "[ENCRYPTED]"
				}
			}

			fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\n",
				task.ID,
				status,
				description,
				createdTime,
				secretStatus,
			)
		}

		w.Flush()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().BoolVarP(&showSecretsFlag, "show-secrets", "d", false, "Decrypt and display secret tasks")
}
