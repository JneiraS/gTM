package domain

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/JneiraS/AMS/src/domain/models"
)

func TestTaskCreation(t *testing.T) {
	now := time.Now()
	task := models.Task{
		Title:         "Test Task",
		Description:   "Test Description",
		DueDate:       now,
		Status:        "En cours",
		Priority:      "Haute",
		Assignee:      "John Doe",
		Creator:       "Jane Doe",
		Project:       "Test Project",
		Progress:      50,
		EstimatedTime: 8,
		TimeSpent:     4,
	}

	if task.Title != "Test Task" {
		t.Errorf("Expected title 'Test Task', got %s", task.Title)
	}

	if task.Progress < 0 {
		t.Errorf("Progress should be between 0 and 100, got %d", task.Progress)
	}

	if task.Progress > 100 {
		t.Errorf("Progress should be between 0 and 100, got %d", task.Progress)
	}
}

func TestTaskJSONSerialization(t *testing.T) {
	now := time.Now()
	task := models.Task{
		Title:         "JSON Test",
		Description:   "Test JSON Serialization",
		DueDate:       now,
		Status:        "En cours",
		Priority:      "Moyenne",
		Assignee:      "John Doe",
		Creator:       "Jane Doe",
		Project:       "Test Project",
		Progress:      75,
		EstimatedTime: 16,
		TimeSpent:     8,
	}

	jsonData, err := json.Marshal(task)
	if err != nil {
		t.Fatalf("Failed to marshal task: %v", err)
	}

	var unmarshaledTask models.Task
	err = json.Unmarshal(jsonData, &unmarshaledTask)
	if err != nil {
		t.Fatalf("Failed to unmarshal task: %v", err)
	}

	if unmarshaledTask.Title != task.Title {
		t.Errorf("Expected title %s, got %s", task.Title, unmarshaledTask.Title)
	}
}

func TestSubtaskValidation(t *testing.T) {
	subtask := models.Subtask{
		Title:  "Sub-task 1",
		Status: "En cours",
	}

	if subtask.Title == "" {
		t.Error("Subtask title should not be empty")
	}

	if subtask.Status == "" {
		t.Error("Subtask status should not be empty")
	}
}

func TestCommentValidation(t *testing.T) {
	comment := models.Comment{
		Author: "John Doe",
		Text:   "Test comment",
	}

	if comment.Author == "" {
		t.Error("Comment author should not be empty")
	}

	if comment.Text == "" {
		t.Error("Comment text should not be empty")
	}
}

func TestTagValidation(t *testing.T) {
	tag := models.Tag{
		Name: "important",
	}

	if tag.Name == "" {
		t.Error("Tag name should not be empty")
	}
}
