package persistence

import (
	"fmt"
	"strings"
	"time"

	"github.com/JneiraS/AMS/src/domain/models"
	"gorm.io/gorm"
)

const (
	REQUEST_TASK_ID = "task_id = ?"
)

type Task struct {
	gorm.Model
	models.Task
}
type Subtask struct {
	gorm.Model
	models.Subtask
}

type TaskSubtasks struct {
	SubtaskID uint `gorm:"primaryKey"`
	TaskID    uint `gorm:"primaryKey"`
}

type Comment struct {
	gorm.Model
	models.Comment
}

type TaskComments struct {
	CommentID uint `gorm:"primaryKey"`
	TaskID    uint `gorm:"primaryKey"`
}

type Tag struct {
	gorm.Model
	models.Tag
}
type TaskTags struct {
	TaskID uint `gorm:"primaryKey"`
	TagID  uint `gorm:"primaryKey"`
}

type TaskTimeSpent struct {
	gorm.Model
	models.TimeSpent
	TaskID uint `gorm:"primaryKey"`
}

//------CREATE------

// Crée une nouvelle tâche dans la base de données.
func CreateTask(db *gorm.DB, task Task) {
	// Convert Windows style line endings (\r\n) to Unix style (\n)
	task.Description = strings.ReplaceAll(task.Description, "\r\n", "\n")
	// Convert single carriage returns to line breaks
	task.Description = strings.ReplaceAll(task.Description, "\r", "\n")

	transactionQueue <- Transaction{
		Func: func(tx *gorm.DB) error {
			if err := tx.Create(&task).Error; err != nil {
				return err
			}
			CreateTaskTimeSpent(tx, task.ID)
			return nil
		},
	}
}

// Crée une nouvelle sous-tâche dans la base de données.
func CreateSubtask(db *gorm.DB, subtask Subtask, TaskID uint) {
	db.Create(&subtask)

	taskSubtask := TaskSubtasks{
		SubtaskID: subtask.ID,
		TaskID:    TaskID,
	}
	db.Create(&taskSubtask)
}

// Crée un nouveau commentaire dans la base de données.
func CreateComment(db *gorm.DB, comment Comment, TaskID uint) {
	db.Create(&comment)

	TaskComments := TaskComments{
		CommentID: comment.ID,
		TaskID:    TaskID,
	}
	db.Create(&TaskComments)
}

func CreateTaskTimeSpent(db *gorm.DB, taskID uint) (TaskTimeSpent, error) {
	taskTimeSpent := TaskTimeSpent{
		TimeSpent: models.TimeSpent{
			StartTime: time.Now(),
			EndTime:   time.Time{},
		},
		TaskID: taskID,
	}

	result := db.Create(&taskTimeSpent)
	if result.RowsAffected == 0 {
		return TaskTimeSpent{}, fmt.Errorf("no record created for task ID: %d", taskID)
	}

	if result.Error != nil {
		return TaskTimeSpent{}, result.Error
	}

	// Verify the record was created by retrieving it
	var created TaskTimeSpent
	if err := db.Where(REQUEST_TASK_ID, taskID).First(&created).Error; err != nil {
		return TaskTimeSpent{}, fmt.Errorf("failed to verify created record: %v", err)
	}

	return created, nil
}

//------READ------

// Récupère une tâche de la base de données.
func GetTask(db *gorm.DB, id uint) (Task, error) {
	var task Task
	r := db.Find(&task, id)
	if r.Error == gorm.ErrRecordNotFound {
		return Task{}, fmt.Errorf("task not found")
	}

	return task, nil
}

// Récupère une sous-tâche de la base de données.
func GetSubtask(db *gorm.DB, id uint) Subtask {
	var subtask Subtask
	db.First(&subtask, id)
	return subtask
}

// Récupère des commentaires de la base de données.
func GetComment(db *gorm.DB, id uint) []Comment {
	var comments []Comment
	db.First(&comments, id)
	return comments
}

func GetTaskComments(db *gorm.DB, id uint) ([]TaskComments, error) {
	var taskComments []TaskComments
	result := db.Where(REQUEST_TASK_ID, id).Order("comment_id desc").Find(&taskComments)
	if result.Error != nil {
		return nil, fmt.Errorf("error while retrieving task comments: %w", result.Error)
	}
	return taskComments, nil
}

func GetAllCommentsOfTask(db *gorm.DB, taskID uint) ([]Comment, error) {
	var comments []Comment
	taskComments, err := GetTaskComments(db, taskID)
	if err != nil {
		return nil, err
	}
	for _, taskComment := range taskComments {
		comment := GetComment(db, taskComment.CommentID)
		comments = append(comments, comment...)
	}
	return comments, nil
}

func GetTaskSubtasks(db *gorm.DB, id uint) ([]TaskSubtasks, error) {
	var taskSubtasks []TaskSubtasks
	result := db.Where(REQUEST_TASK_ID, id).Order("subtask_id desc").Find(&taskSubtasks)
	if result.Error != nil {
		return nil, fmt.Errorf("error while retrieving task subtasks: %w", result.Error)
	}
	return taskSubtasks, nil
}

func GetAllSubtasksOfTask(db *gorm.DB, taskID uint) ([]Subtask, error) {
	var subtasks []Subtask
	taskSubtasks, err := GetTaskSubtasks(db, taskID)
	if err != nil {
		return nil, err
	}
	for _, taskSubtask := range taskSubtasks {
		subtask := GetSubtask(db, taskSubtask.SubtaskID)
		subtasks = append(subtasks, subtask)
	}
	return subtasks, nil
}

// GetAllProjects renvoie une liste de tous les noms de projet dans la base de données.
func GetAllProjects(db *gorm.DB) []string {
	var projects []string
	db.Model(&Task{}).Distinct("project").Where("project != ?", "").Pluck("project", &projects)
	return projects
}

func GetTaskTimeSpent(db *gorm.DB, taskID uint) TaskTimeSpent {
	var taskTimeSpent TaskTimeSpent
	db.Where(REQUEST_TASK_ID, taskID).First(&taskTimeSpent)
	return taskTimeSpent
}

//------UPDATE------

// Mettre à jour une tâche dans la base de données.
func UpdateTask(db *gorm.DB, task Task) {
	transactionQueue <- Transaction{
		Func: func(tx *gorm.DB) error {
			return tx.Model(&task).Updates(map[string]interface{}{
				"Title":         task.Title,
				"Description":   task.Description,
				"DueDate":       task.DueDate,
				"Status":        task.Status,
				"Priority":      task.Priority,
				"Assignee":      task.Assignee,
				"Creator":       task.Creator,
				"Project":       task.Project,
				"Progress":      task.Progress,
				"EstimatedTime": task.EstimatedTime,
				"TimeSpent":     task.TimeSpent,
			}).Error
		},
	}
}

// Mettre à jour une sous-tâche dans la base de données.
func UpdateSubtask(db *gorm.DB, subtask Subtask) {
	db.Model(&subtask).Updates(map[string]interface{}{
		"Title":  subtask.Title,
		"Status": subtask.Status,
	})
}

// Mettre à jour un commentaire dans la base de données.
func UpdateComment(db *gorm.DB, comment Comment) {
	db.Model(&comment).Updates(map[string]interface{}{
		"Author": comment.Author,
		"Text":   comment.Text,
	})
}

func UpdateTaskTimeSpent(db *gorm.DB, taskTimeSpent TaskTimeSpent) {
	db.Model(&taskTimeSpent).Updates(map[string]interface{}{
		"StartTime": taskTimeSpent.StartTime,
		"EndTime":   taskTimeSpent.EndTime,
	})
}

//------DELETE------

// Supprime une tâche de la base de données.
func DeleteTask(db *gorm.DB, id uint) {
	var task Task
	db.Delete(&task, id)
}

// Supprime une sous-tâche de la base de données.
func DeleteSubtask(db *gorm.DB, id uint) {
	var subtask Subtask
	db.Delete(&subtask, id)
}

// Supprime un commentaire de la base de données.
func DeleteComment(db *gorm.DB, id uint) {
	var comment Comment
	db.Delete(&comment, id)

}
