package todo

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/csv"
	"errors"
	"io"
	"os"
	"strconv"
	"time"
)

type Task struct {
	ID          int
	Description string
	Completed   bool
	CreatedAt   time.Time
	CompletedAt time.Time
	Encrypted   bool
}

// EncryptText encrypts plaintext string with AES-GCM and returns base64 encoded result
func EncryptText(plaintext string) (string, error) {
	// Use a static key for simplicity - in a real app should use a better key management system
	key := sha256.Sum256([]byte("R2D2SecretKey"))

	block, err := aes.NewCipher(key[:])
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptText decrypts base64 encoded ciphertext with AES-GCM
func DecryptText(encryptedText string) (string, error) {
	// Use a static key for simplicity - in a real app should use a better key management system
	key := sha256.Sum256([]byte("R2D2SecretKey"))

	ciphertext, err := base64.StdEncoding.DecodeString(encryptedText)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key[:])
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < gcm.NonceSize() {
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:gcm.NonceSize()], ciphertext[gcm.NonceSize():]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

func SaveTasks(filename string, tasks []Task) error {
	file, err := os.Create(filename)
	if err != nil {
		return errors.New("failed to create file")
	}
	defer file.Close()

	write := csv.NewWriter(file)
	defer write.Flush()

	for _, task := range tasks {
		record := []string{
			strconv.Itoa(task.ID),
			task.Description,
			strconv.FormatBool(task.Completed),
			task.CreatedAt.Format(time.RFC3339),
			task.CompletedAt.Format(time.RFC3339),
			strconv.FormatBool(task.Encrypted),
		}
		if !task.CompletedAt.IsZero() {
			record[4] = task.CompletedAt.Format(time.RFC3339)
		}
		if err := write.Write(record); err != nil {
			return errors.New("failed to write record")
		}
	}
	return nil
}

func LoadTasks(filename string) ([]Task, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, errors.New("failed to open file")
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll() //just because it's a small file
	if err != nil {
		return nil, errors.New("failed to read records")
	}
	tasks := []Task{}
	for _, record := range records {
		id, err := strconv.Atoi(record[0])
		if err != nil {
			return nil, errors.New("failed to convert ID to int")
		}
		completed, err := strconv.ParseBool(record[2])
		if err != nil {
			return nil, errors.New("failed to convert completed to bool")
		}
		createdAt, err := time.Parse(time.RFC3339, record[3])
		if err != nil {
			return nil, errors.New("failed to parse createdAt")
		}
		completedAt := time.Time{}
		if record[4] != "" {
			completedAt, err = time.Parse(time.RFC3339, record[4])
			if err != nil {
				return nil, errors.New("failed to parse completedAt")
			}
		}

		encrypted := false
		if len(record) > 5 {
			encrypted, _ = strconv.ParseBool(record[5])
		}

		tasks = append(tasks, Task{
			ID:          id,
			Description: record[1],
			Completed:   completed,
			CreatedAt:   createdAt,
			CompletedAt: completedAt,
			Encrypted:   encrypted,
		})
	}
	return tasks, nil
}
