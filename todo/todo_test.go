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

func TestEncryptDecryptText(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		wantSame bool
	}{
		{
			name:     "Simple text",
			input:    "This is a secret task",
			wantSame: true,
		},
		{
			name:     "Empty string",
			input:    "",
			wantSame: true,
		},
		{
			name:     "With special characters",
			input:    "Secret task with !@#$%^&*()_+",
			wantSame: true,
		},
		{
			name:     "Long text",
			input:    "This is a very long secret task description that needs to be encrypted and then decrypted back to the original text without any data loss or corruption in the process of encryption and decryption.",
			wantSame: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Encrypt the input
			encrypted, err := EncryptText(tc.input)
			if err != nil {
				t.Fatalf("EncryptText() error = %v", err)
			}

			// Verify encryption produced different output
			if encrypted == tc.input && tc.input != "" {
				t.Errorf("EncryptText() didn't change the input, got = %v", encrypted)
			}

			// Decrypt back to original
			decrypted, err := DecryptText(encrypted)
			if err != nil {
				t.Fatalf("DecryptText() error = %v", err)
			}

			// Check if decryption resulted in the original input
			if tc.wantSame && decrypted != tc.input {
				t.Errorf("DecryptText() = %v, want %v", decrypted, tc.input)
			}
		})
	}
}

func TestInvalidDecryption(t *testing.T) {
	testCases := []struct {
		name        string
		invalidText string
		wantErr     bool
	}{
		{
			name:        "Invalid base64",
			invalidText: "This is not valid base64!",
			wantErr:     true,
		},
		{
			name:        "Too short after base64 decode",
			invalidText: "aGVsbG8=", // "hello" base64 encoded
			wantErr:     true,
		},
		{
			name:        "Empty string",
			invalidText: "",
			wantErr:     true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := DecryptText(tc.invalidText)
			if (err != nil) != tc.wantErr {
				t.Errorf("DecryptText() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

// Test that the encryption is consistent across different runs
func TestEncryptionDifferentRuns(t *testing.T) {
	input := "This is a secret message"

	// First encryption
	encrypted1, err := EncryptText(input)
	if err != nil {
		t.Fatalf("First EncryptText() error = %v", err)
	}

	// Second encryption
	encrypted2, err := EncryptText(input)
	if err != nil {
		t.Fatalf("Second EncryptText() error = %v", err)
	}

	// The encrypted values should be different due to random nonce
	if encrypted1 == encrypted2 {
		t.Errorf("Multiple encryptions of the same text should produce different results")
	}

	// But both should decrypt to the original text
	decrypted1, err := DecryptText(encrypted1)
	if err != nil {
		t.Fatalf("DecryptText() first error = %v", err)
	}

	decrypted2, err := DecryptText(encrypted2)
	if err != nil {
		t.Fatalf("DecryptText() second error = %v", err)
	}

	if decrypted1 != input || decrypted2 != input {
		t.Errorf("Decrypted texts do not match original input")
	}
}
