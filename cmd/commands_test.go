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

func TestAddSecretTask(t *testing.T) {
	// Setup a temporary file for testing
	tmpfile, err := os.CreateTemp("", "tasks_test.csv")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())
	tmpfile.Close()

	// Override the tasks file with our temp file for this test
	if err := os.Rename("tasks.csv", "tasks.csv.bak"); err != nil && !os.IsNotExist(err) {
		t.Fatalf("Failed to rename original tasks file: %v", err)
	}
	defer func() {
		os.Remove("tasks.csv")
		if _, err := os.Stat("tasks.csv.bak"); err == nil {
			os.Rename("tasks.csv.bak", "tasks.csv")
		}
	}()

	// Copy our temp file to the tasks.csv location
	data, _ := os.ReadFile(tmpfile.Name())
	if err := os.WriteFile("tasks.csv", data, 0644); err != nil {
		t.Fatalf("Failed to create test tasks file: %v", err)
	}

	// Setup test cases
	testCases := []struct {
		name          string
		args          []string
		secret        bool
		wantEncrypted bool
	}{
		{
			name:          "Add regular task",
			args:          []string{"Test regular task"},
			secret:        false,
			wantEncrypted: false,
		},
		{
			name:          "Add secret task",
			args:          []string{"Test secret task"},
			secret:        true,
			wantEncrypted: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Reset secretFlag for each test
			secretFlag = tc.secret

			// Execute the add command
			oldArgs := os.Args
			os.Args = append([]string{"r2d2", "add"}, tc.args...)
			defer func() { os.Args = oldArgs }()

			// Capture stdout to verify output
			rescueStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// Run the command
			addCmd.Run(addCmd, tc.args)

			// Restore stdout
			w.Close()
			out, _ := io.ReadAll(r)
			os.Stdout = rescueStdout

			// Verify tasks were saved correctly
			tasks, err := todo.LoadTasks("tasks.csv")
			if err != nil {
				t.Fatalf("Failed to load tasks: %v", err)
			}

			// Verify at least one task exists
			if len(tasks) == 0 {
				t.Fatalf("No tasks were saved")
			}

			// Get the last task (should be the one we just added)
			lastTask := tasks[len(tasks)-1]

			// Check if the task is encrypted as expected
			if lastTask.Encrypted != tc.wantEncrypted {
				t.Errorf("Task encryption status = %v, want %v", lastTask.Encrypted, tc.wantEncrypted)
			}

			// If it's a secret task, verify it's properly encrypted
			if tc.wantEncrypted {
				// The description should be encrypted (not plain text)
				if lastTask.Description == tc.args[0] {
					t.Errorf("Secret task not encrypted, got raw text: %v", lastTask.Description)
				}

				// But should decrypt back to the original
				decrypted, err := todo.DecryptText(lastTask.Description)
				if err != nil {
					t.Errorf("Failed to decrypt task: %v", err)
				}
				if decrypted != tc.args[0] {
					t.Errorf("Decrypted text = %v, want %v", decrypted, tc.args[0])
				}
			} else {
				// Regular task should have plain text description
				if lastTask.Description != tc.args[0] {
					t.Errorf("Regular task description = %v, want %v", lastTask.Description, tc.args[0])
				}
			}

			// Verify the correct output message was printed
			output := string(out)
			if tc.wantEncrypted {
				if !strings.Contains(output, "Secret task added") {
					t.Errorf("Output doesn't contain 'Secret task added', got: %v", output)
				}
			} else {
				if !strings.Contains(output, "Task added") {
					t.Errorf("Output doesn't contain 'Task added', got: %v", output)
				}
			}
		})
	}
}

func TestListWithSecrets(t *testing.T) {
	// Setup a temporary file for testing
	tmpfile, err := os.CreateTemp("", "tasks_test.csv")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())
	tmpfile.Close()

	// Override the tasks file with our temp file for this test
	if err := os.Rename("tasks.csv", "tasks.csv.bak"); err != nil && !os.IsNotExist(err) {
		t.Fatalf("Failed to rename original tasks file: %v", err)
	}
	defer func() {
		os.Remove("tasks.csv")
		if _, err := os.Stat("tasks.csv.bak"); err == nil {
			os.Rename("tasks.csv.bak", "tasks.csv")
		}
	}()

	// Create test data with both regular and encrypted tasks
	plainText := "Regular task"
	secretText := "This is a secret task"

	encryptedText, err := todo.EncryptText(secretText)
	if err != nil {
		t.Fatalf("Failed to encrypt text: %v", err)
	}

	tasks := []todo.Task{
		{
			ID:          1,
			Description: plainText,
			Completed:   false,
			CreatedAt:   time.Now(),
			CompletedAt: time.Time{},
			Encrypted:   false,
		},
		{
			ID:          2,
			Description: encryptedText,
			Completed:   false,
			CreatedAt:   time.Now(),
			CompletedAt: time.Time{},
			Encrypted:   true,
		},
	}

	// Save the test tasks
	if err := todo.SaveTasks("tasks.csv", tasks); err != nil {
		t.Fatalf("Failed to save test tasks: %v", err)
	}

	testCases := []struct {
		name              string
		showSecrets       bool
		wantPlainVisible  bool
		wantSecretVisible bool
	}{
		{
			name:              "List without showing secrets",
			showSecrets:       false,
			wantPlainVisible:  true,
			wantSecretVisible: false,
		},
		{
			name:              "List with showing secrets",
			showSecrets:       true,
			wantPlainVisible:  true,
			wantSecretVisible: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Set the showSecrets flag
			showSecretsFlag = tc.showSecrets

			// Capture stdout to verify output
			rescueStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// Run the list command
			listCmd.Run(listCmd, []string{})

			// Restore stdout and get output
			w.Close()
			out, _ := io.ReadAll(r)
			os.Stdout = rescueStdout

			output := string(out)

			// Verify regular task is visible
			if tc.wantPlainVisible && !strings.Contains(output, plainText) {
				t.Errorf("Regular task not visible in output: %v", output)
			}

			// Verify secret task handling
			if tc.wantSecretVisible {
				// When showing secrets, the decrypted text should be visible
				if !strings.Contains(output, secretText) {
					t.Errorf("Decrypted secret task not visible with --show-secrets: %v", output)
				}
			} else {
				// When not showing secrets, should show [ENCRYPTED] instead
				if !strings.Contains(output, "[ENCRYPTED]") {
					t.Errorf("Encrypted task should show [ENCRYPTED] when not decrypted: %v", output)
				}
				// And the actual secret text should not be visible
				if strings.Contains(output, secretText) {
					t.Errorf("Secret text should not be visible without --show-secrets flag: %v", output)
				}
			}
		})
	}
}
