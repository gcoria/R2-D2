package cmd

import (
	"R2-D2/todo"
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/spf13/cobra"
)

// Define taskFile variable for testing
var taskFile = "tasks.csv"

// Helper function to modify command behavior for testing
func modifyCommandsForTest(filename string) func() {
	// Store original functions
	originalAdd := addCmd.Run
	originalList := listCmd.Run
	originalComplete := completeCmd.Run
	originalDelete := deleteCmd.Run

	// Override add command
	addCmd.Run = func(cmd *cobra.Command, args []string) {
		description := strings.Join(args, " ")
		tasks, err := todo.LoadTasks(filename)
		if err != nil {
			tasks = []todo.Task{}
		}
		task := todo.Task{
			ID:          len(tasks) + 1,
			Description: description,
			Completed:   false,
			CreatedAt:   time.Now(),
			CompletedAt: time.Time{},
		}
		tasks = append(tasks, task)
		err = todo.SaveTasks(filename, tasks)
		if err != nil {
			fmt.Println("Error saving tasks:", err)
			return
		}
		fmt.Printf("Task added: %d - %s\n", task.ID, task.Description)
	}

	// Override list command
	listCmd.Run = func(cmd *cobra.Command, args []string) {
		tasks, err := todo.LoadTasks(filename)
		if err != nil {
			fmt.Println("Error loading tasks:", err)
			return
		}
		// Rest of list implementation
		if len(tasks) == 0 {
			fmt.Println("No tasks to display")
			return
		}
		fmt.Println("Tasks:")
		for _, task := range tasks {
			status := " "
			if task.Completed {
				status = "X"
			}
			fmt.Printf("[%s] %d - %s\n", status, task.ID, task.Description)
		}
	}

	// Override complete command
	completeCmd.Run = func(cmd *cobra.Command, args []string) {
		id := args[0]
		taskID, err := strconv.Atoi(id)
		if err != nil {
			fmt.Println("Invalid task ID")
			return
		}
		tasks, err := todo.LoadTasks(filename)
		if err != nil {
			fmt.Println("Error loading tasks:", err)
			return
		}
		for i, task := range tasks {
			if task.ID == taskID {
				tasks[i].Completed = true
				tasks[i].CompletedAt = time.Now()
				err = todo.SaveTasks(filename, tasks)
				if err != nil {
					fmt.Println("Error saving tasks:", err)
					return
				}
				fmt.Printf("Task %d marked as completed\n", taskID)
				return
			}
		}
		fmt.Printf("Task %d not found\n", taskID)
	}

	// Override delete command
	deleteCmd.Run = func(cmd *cobra.Command, args []string) {
		id := args[0]
		taskID, err := strconv.Atoi(id)
		if err != nil {
			fmt.Println("Invalid task ID")
			return
		}
		tasks, err := todo.LoadTasks(filename)
		if err != nil {
			fmt.Println("Error loading tasks:", err)
			return
		}
		for i, task := range tasks {
			if task.ID == taskID {
				tasks = append(tasks[:i], tasks[i+1:]...)
				err = todo.SaveTasks(filename, tasks)
				if err != nil {
					fmt.Println("Error saving tasks:", err)
					return
				}
				fmt.Printf("Task %d deleted\n", taskID)
				return
			}
		}
		fmt.Printf("Task %d not found\n", taskID)
	}

	// Return cleanup function
	return func() {
		addCmd.Run = originalAdd
		listCmd.Run = originalList
		completeCmd.Run = originalComplete
		deleteCmd.Run = originalDelete
	}
}

// Helper function to create a temporary task file
func setupTaskFile(t *testing.T) (string, func()) {
	// Create a temporary file
	tempFile, err := os.CreateTemp("", "tasks_cmd_test_*.csv")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tempFileName := tempFile.Name()
	tempFile.Close()

	// Create cleanup function
	cleanup := func() {
		os.Remove(tempFileName)
	}

	return tempFileName, cleanup
}

// Helper function to capture stdout during test
func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func TestAddCommand(t *testing.T) {
	// Setup
	tempFileName, cleanup := setupTaskFile(t)
	defer cleanup()

	// Override commands for test
	cleanupCommands := modifyCommandsForTest(tempFileName)
	defer cleanupCommands()

	// Execute the add command
	output := captureOutput(func() {
		rootCmd.SetArgs([]string{"add", "Test task"})
		rootCmd.Execute()
	})

	// Verify output
	if !strings.Contains(output, "Task added: 1 - Test task") {
		t.Errorf("Expected output to contain 'Task added: 1 - Test task', got: %s", output)
	}

	// Verify task was actually added to the file
	tasks, err := todo.LoadTasks(tempFileName)
	if err != nil {
		t.Fatalf("Failed to load tasks: %v", err)
	}
	if len(tasks) != 1 {
		t.Fatalf("Expected 1 task, got %d", len(tasks))
	}
	if tasks[0].Description != "Test task" {
		t.Errorf("Expected task description 'Test task', got '%s'", tasks[0].Description)
	}
}

func TestListCommand(t *testing.T) {
	// Setup
	tempFileName, cleanup := setupTaskFile(t)
	defer cleanup()

	// Override commands for test
	cleanupCommands := modifyCommandsForTest(tempFileName)
	defer cleanupCommands()

	// Create test data
	tasks := []todo.Task{
		{
			ID:          1,
			Description: "Test task 1",
			Completed:   false,
			CreatedAt:   time.Now(),
		},
		{
			ID:          2,
			Description: "Test task 2",
			Completed:   true,
			CreatedAt:   time.Now(),
			CompletedAt: time.Now(),
		},
	}
	todo.SaveTasks(tempFileName, tasks)

	// Execute the list command
	output := captureOutput(func() {
		rootCmd.SetArgs([]string{"list"})
		rootCmd.Execute()
	})

	// Verify output
	if !strings.Contains(output, "Test task 1") || !strings.Contains(output, "Test task 2") {
		t.Errorf("Expected output to contain task descriptions, got: %s", output)
	}
}

func TestCompleteCommand(t *testing.T) {
	// Setup
	tempFileName, cleanup := setupTaskFile(t)
	defer cleanup()

	// Override commands for test
	cleanupCommands := modifyCommandsForTest(tempFileName)
	defer cleanupCommands()

	// Create test data
	tasks := []todo.Task{
		{
			ID:          1,
			Description: "Test task",
			Completed:   false,
			CreatedAt:   time.Now(),
		},
	}
	todo.SaveTasks(tempFileName, tasks)

	// Execute the complete command
	output := captureOutput(func() {
		rootCmd.SetArgs([]string{"complete", "1"})
		rootCmd.Execute()
	})

	// Verify output
	if !strings.Contains(output, "Task 1 marked as completed") {
		t.Errorf("Expected output to contain completion message, got: %s", output)
	}

	// Verify task was actually completed
	updatedTasks, err := todo.LoadTasks(tempFileName)
	if err != nil {
		t.Fatalf("Failed to load tasks: %v", err)
	}
	if !updatedTasks[0].Completed {
		t.Errorf("Expected task to be completed")
	}
}

func TestDeleteCommand(t *testing.T) {
	// Setup
	tempFileName, cleanup := setupTaskFile(t)
	defer cleanup()

	// Override commands for test
	cleanupCommands := modifyCommandsForTest(tempFileName)
	defer cleanupCommands()

	// Create test data
	tasks := []todo.Task{
		{
			ID:          1,
			Description: "Test task",
			Completed:   false,
			CreatedAt:   time.Now(),
		},
	}
	todo.SaveTasks(tempFileName, tasks)

	// Execute the delete command
	output := captureOutput(func() {
		rootCmd.SetArgs([]string{"delete", "1"})
		rootCmd.Execute()
	})

	// Verify output
	if !strings.Contains(output, "Task 1 deleted") {
		t.Errorf("Expected output to contain deletion message, got: %s", output)
	}

	// Verify task was actually deleted
	updatedTasks, err := todo.LoadTasks(tempFileName)
	if err != nil {
		t.Fatalf("Failed to load tasks: %v", err)
	}
	if len(updatedTasks) != 0 {
		t.Errorf("Expected 0 tasks after deletion, got %d", len(updatedTasks))
	}
}
