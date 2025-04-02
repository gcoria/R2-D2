package todo

import (
	"os"
	"testing"
	"time"
)

func TestSaveAndLoadTasks(t *testing.T) {
	// Create a temporary file for testing
	tempFile, err := os.CreateTemp("", "tasks_test_*.csv")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tempFileName := tempFile.Name()
	defer os.Remove(tempFileName) // Clean up after test
	tempFile.Close()

	// Create test data
	createdTime := time.Now().Round(time.Second) // Round to avoid precision issues
	completedTime := createdTime.Add(1 * time.Hour).Round(time.Second)

	originalTasks := []Task{
		{
			ID:          1,
			Description: "Test task 1",
			Completed:   false,
			CreatedAt:   createdTime,
			CompletedAt: time.Time{},
		},
		{
			ID:          2,
			Description: "Test task 2",
			Completed:   true,
			CreatedAt:   createdTime,
			CompletedAt: completedTime,
		},
	}

	// Test SaveTasks
	if err := SaveTasks(tempFileName, originalTasks); err != nil {
		t.Fatalf("SaveTasks failed: %v", err)
	}

	// Test LoadTasks
	loadedTasks, err := LoadTasks(tempFileName)
	if err != nil {
		t.Fatalf("LoadTasks failed: %v", err)
	}

	// Verify that loaded tasks match the original tasks
	if len(loadedTasks) != len(originalTasks) {
		t.Errorf("Expected %d tasks, got %d", len(originalTasks), len(loadedTasks))
	}

	for i, original := range originalTasks {
		loaded := loadedTasks[i]
		if original.ID != loaded.ID {
			t.Errorf("Task %d: Expected ID %d, got %d", i, original.ID, loaded.ID)
		}
		if original.Description != loaded.Description {
			t.Errorf("Task %d: Expected Description %s, got %s", i, original.Description, loaded.Description)
		}
		if original.Completed != loaded.Completed {
			t.Errorf("Task %d: Expected Completed %v, got %v", i, original.Completed, loaded.Completed)
		}
		if !original.CreatedAt.Equal(loaded.CreatedAt) {
			t.Errorf("Task %d: Expected CreatedAt %v, got %v", i, original.CreatedAt, loaded.CreatedAt)
		}
		if original.Completed && !original.CompletedAt.Equal(loaded.CompletedAt) {
			t.Errorf("Task %d: Expected CompletedAt %v, got %v", i, original.CompletedAt, loaded.CompletedAt)
		}
	}
}

func TestLoadTasksError(t *testing.T) {
	// Test loading from a non-existent file
	_, err := LoadTasks("non_existent_file.csv")
	if err == nil {
		t.Error("Expected error when loading from non-existent file, got nil")
	}
}

func TestEmptyTasksList(t *testing.T) {
	// Create a temporary file for testing
	tempFile, err := os.CreateTemp("", "tasks_test_empty_*.csv")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tempFileName := tempFile.Name()
	defer os.Remove(tempFileName) // Clean up after test
	tempFile.Close()

	// Save empty task list
	emptyTasks := []Task{}
	err = SaveTasks(tempFileName, emptyTasks)
	if err != nil {
		t.Fatalf("SaveTasks failed: %v", err)
	}

	// Load the empty task list
	loadedTasks, err := LoadTasks(tempFileName)
	if err != nil {
		t.Fatalf("LoadTasks failed: %v", err)
	}

	if len(loadedTasks) != 0 {
		t.Errorf("Expected 0 tasks, got %d", len(loadedTasks))
	}
}

func TestSaveTasksError(t *testing.T) {
	// Test saving to an invalid path
	err := SaveTasks("/nonexistent/directory/file.csv", []Task{
		{
			ID:          1,
			Description: "Test task",
			Completed:   false,
			CreatedAt:   time.Now(),
		},
	})
	if err == nil {
		t.Error("Expected error when saving to invalid path, got nil")
	}
}
