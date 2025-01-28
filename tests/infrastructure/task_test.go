package persistence

import (
	"testing"
	"time"

	"github.com/JneiraS/AMS/src/domain/models"
	"github.com/JneiraS/AMS/src/infrastructure/persistence"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Auto migrate the schemas
	db.AutoMigrate(&persistence.Task{}, &persistence.Subtask{}, &persistence.Comment{}, &persistence.Tag{}, &persistence.TaskSubtasks{}, &persistence.TaskComments{}, &persistence.TaskTags{})

	return db
}

func TestCreateTask(t *testing.T) {
	db := setupTestDB(t)

	originalTask := persistence.Task{
		Task: models.Task{
			Title:         "Test Task",
			Description:   "Test Description",
			DueDate:       time.Now(),
			Status:        "pending",
			Priority:      "low",
			Assignee:      "test-user",
			Creator:       "admin",
			Project:       "test-project",
			Progress:      0,
			EstimatedTime: 2,
		},
	}

	persistence.CreateTask(db, originalTask)

	var savedTask persistence.Task
	result := db.First(&savedTask, 1)
	if result.Error != nil {
		t.Fatalf("Failed to retrieve created task: %v", result.Error)
	}

	if savedTask.Task.Title != originalTask.Task.Title {
		t.Errorf("Expected task title %s, got %s", originalTask.Task.Title, savedTask.Task.Title)
	}
}

func TestGetTask(t *testing.T) {
	db := setupTestDB(t)

	// Create a test task
	originalTask := persistence.Task{
		Task: models.Task{
			Title: "Test Get Task",
		},
	}
	db.Create(&originalTask)

	// Test retrieving existing task
	retrievedTask, _ := persistence.GetTask(db, originalTask.ID)
	if retrievedTask.Title != originalTask.Title {
		t.Errorf("Expected task title %s, got %s", originalTask.Title, retrievedTask.Title)
	}

	// Test retrieving non-existent task
	nonExistentTask, _ := persistence.GetTask(db, 99999)
	if nonExistentTask.ID != 0 {
		t.Error("Expected empty task for non-existent ID")
	}
}

func TestUpdateTask(t *testing.T) {
	db := setupTestDB(t)

	task := persistence.Task{
		Task: models.Task{
			Title:       "Original Title",
			Description: "Original Description",
		},
	}
	db.Create(&task)

	// Update task
	task.Title = "Updated Title"
	task.Description = "Updated Description"
	persistence.UpdateTask(db, task)

	// Verify update
	var updatedTask persistence.Task
	db.First(&updatedTask, task.ID)
	if updatedTask.Title != "Updated Title" {
		t.Errorf("Expected updated title 'Updated Title', got %s", updatedTask.Title)
	}
}

func TestDeleteTask(t *testing.T) {
	db := setupTestDB(t)

	task := persistence.Task{
		Task: models.Task{
			Title: "Task to Delete",
		},
	}
	db.Create(&task)

	// Delete task
	persistence.DeleteTask(db, task.ID)

	// Verify deletion
	var deletedTask persistence.Task
	result := db.First(&deletedTask, task.ID)
	if result.Error == nil {
		t.Error("Expected task to be deleted")
	}
}

func TestCreateSubtask(t *testing.T) {
	db := setupTestDB(t)

	subtask := persistence.Subtask{
		Subtask: models.Subtask{
			Title: "Test Subtask",
		},
	}

	persistence.CreateSubtask(db, subtask)

	var savedSubtask persistence.Subtask
	result := db.First(&savedSubtask, 1)
	if result.Error != nil {
		t.Errorf("Failed to retrieve created subtask: %v", result.Error)
	}
}

func TestCreateComment(t *testing.T) {
	db := setupTestDB(t)

	comment := persistence.Comment{
		Comment: models.Comment{
			Text: "Test Comment",
		},
	}

	persistence.CreateComment(db, comment)

	var savedComment persistence.Comment
	result := db.First(&savedComment, 1)
	if result.Error != nil {
		t.Errorf("Failed to retrieve created comment: %v", result.Error)
	}
}

func TestGetSubtask(t *testing.T) {
	db := setupTestDB(t)
	// Create a test subtask
	subtask := persistence.Subtask{
		Subtask: models.Subtask{
			Title:  "Test Subtask",
			Status: "En cours",
		},
	}
	db.Create(&subtask)
	// Test retrieving existing subtask
	retrievedSubtask := persistence.GetSubtask(db, subtask.ID)
	if retrievedSubtask.ID != subtask.ID {
		t.Errorf("Expected subtask ID %d, got %d", subtask.ID, retrievedSubtask.ID)
	}
	// Test retrieving non-existent subtask
	nonExistentSubtask := persistence.GetSubtask(db, 99999)
	if nonExistentSubtask.ID != 0 {
		t.Errorf("Expected empty subtask for non-existent ID, got ID %d", nonExistentSubtask.ID)
	}
	// Test retrieving subtask with invalid ID (zero)
	invalidSubtask := persistence.GetSubtask(db, 0)
	if invalidSubtask.ID != 0 {
		t.Errorf("Expected empty subtask for invalid ID, got ID %d", invalidSubtask.ID)
	}
}

func TestGetComment(t *testing.T) {
	db := setupTestDB(t)
	// Create a test comment
	comment := persistence.Comment{
		Comment: models.Comment{
			Text: "Test Comment",
		},
	}
	db.Create(&comment)
	// Test retrieving existing comments
	comments := persistence.GetComment(db, comment.ID)
	if len(comments) != 1 {
		t.Errorf("Expected 1 comment, got %d", len(comments))
	}
	// Test retrieving non-existent comments
	nonExistentComments := persistence.GetComment(db, 99999)
	if len(nonExistentComments) != 0 {
		t.Errorf("Expected 0 comments, got %d", len(nonExistentComments))
	}
	// Test retrieving comments with invalid ID (zero)
	invalidComments := persistence.GetComment(db, 0)
	if len(invalidComments) != 0 {
		t.Errorf("Expected 0 comments, got %d", len(invalidComments))
	}
}
