package todo

import (
	"encoding/csv"
	"errors"
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
		tasks = append(tasks, Task{
			ID:          id,
			Description: record[1],
			Completed:   completed,
			CreatedAt:   createdAt,
			CompletedAt: completedAt,
		})
	}
	return tasks, nil
}
