package persistence

import (
	"fmt"
	"github.com/JneiraS/AMS/src/domain/models"
	"gorm.io/gorm"
	"strings"
	"time"
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
//
// La fonction prend une connexion de base de données GORM et une structure de
// tâche. Elle crée une nouvelle tâche dans la base de données avec les valeurs
// présentes dans la structure de tâche.
//
// La fonction ne renvoie pas de valeur, mais provoquera une panique si l'opération
// de création  choue.
func CreateTask(db *gorm.DB, task Task) {
	// Convert Windows style line endings (\r\n) to Unix style (\n)
	task.Description = strings.ReplaceAll(task.Description, "\r\n", "\n")
	// Convert single carriage returns to line breaks
	task.Description = strings.ReplaceAll(task.Description, "\r", "\n")
	db.Create(&task)
	CreateTaskTimeSpent(db, task.ID)
}

// Crée une nouvelle sous-tâche dans la base de données.
//
// La fonction prend une connexion de base de données GORM et une structure de
// sous-tâche. Elle crée une nouvelle sous-tâche dans la base de données avec les
// valeurs présentes dans la structure de sous-tâche.
//
// La fonction ne renvoie pas de valeur, mais provoquera une panique si l'opération
// de création  choue.
func CreateSubtask(db *gorm.DB, subtask Subtask) {
	db.Create(&subtask)
}

// Crée un nouveau commentaire dans la base de données.
//
// La fonction prend une connexion de base de données GORM et une structure de
// commentaire. Elle crée un nouveau commentaire dans la base de données avec les
// valeurs présentes dans la structure de commentaire.
//
// La fonction ne renvoie pas de valeur, mais provoquera une panique si l'opération
// de création  choue.
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
	if err := db.Where("task_id = ?", taskID).First(&created).Error; err != nil {
		return TaskTimeSpent{}, fmt.Errorf("failed to verify created record: %v", err)
	}

	return created, nil
}

//------READ------

// Récupère une tâche de la base de données.
//
// La fonction prend une connexion de base de données GORM et un identifiant de
// tâche, et renvoie la structure de tâche correspondante, si elle existe.
// Si la tâche n'existe pas, la structure renvoyée est vide.
func GetTask(db *gorm.DB, id uint) (Task, error) {
	var task Task
	r := db.Find(&task, id)
	if r.Error == gorm.ErrRecordNotFound {
		return Task{}, fmt.Errorf("task not found")
	}

	return task, nil
}

// Récupère une sous-tâche de la base de données.
//
// La fonction prend une connexion de base de données GORM et un identifiant de
// sous-tâche, et renvoie la structure de sous-tâche correspondante, si elle existe.
// Si la sous-tâche n'existe pas, la structure renvoyée est vide.
func GetSubtask(db *gorm.DB, id uint) Subtask {
	var subtask Subtask
	db.First(&subtask, id)
	return subtask
}

// Récupère des commentaires de la base de données.
//
// La fonction prend une connexion de base de données GORM et un identifiant de
// commentaire, et renvoie une liste de structures de commentaires correspondantes,
// si elles existent. Si aucun commentaire n'existe pour l'identifiant donné, la
// liste renvoyée est vide.

func GetComment(db *gorm.DB, id uint) []Comment {
	var comments []Comment
	db.First(&comments, id)
	return comments
}

// GetAllProjects renvoie une liste de tous les noms de projet dans la base de données.
//
// La fonction prend une connexion de base de données GORM et renvoie une liste de
// chaînes de caractères correspondant aux noms de projet stockés dans la base de
// données. Si la liste est vide, cela signifie qu'aucun projet n'a été créé.
//
// La fonction renvoie une erreur si la requête SQL échoue.
func GetAllProjects(db *gorm.DB) []string {
	var projects []string
	db.Model(&Task{}).Distinct("project").Where("project != ?", "").Pluck("project", &projects)
	return projects
}

func GetTaskTimeSpent(db *gorm.DB, taskID uint) TaskTimeSpent {
	var taskTimeSpent TaskTimeSpent
	db.Where("task_id = ?", taskID).First(&taskTimeSpent)
	return taskTimeSpent
}

//------UPDATE------

// Mettre à jour une tâche dans la base de donn es.
//
// La fonction prend une connexion de base de donn es GORM et une structure de
// tâche. Elle mettra à jour la tâche dans la base de données avec les valeurs
// présentes dans la structure de tâche.
//
// La fonction ne renvoie pas de valeur, mais provoquera une panique si l'opération
// de mise à jour  choue.
func UpdateTask(db *gorm.DB, task Task) {
	db.Model(&task).Updates(map[string]interface{}{
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
	})
}

// Mettre à jour une sous-tâche dans la base de données.
//
// La fonction prend une connexion de base de données GORM et une structure de
// sous-tâche. Elle mettra à jour la sous-tâche dans la base de données avec les
// valeurs présentes dans la structure de sous-tâche.
//
// La fonction ne renvoie pas de valeur, mais provoquera une panique si l'opération
// de mise à jour  choue.
func UpdateSubtask(db *gorm.DB, subtask Subtask) {
	db.Model(&subtask).Updates(map[string]interface{}{
		"Title":  subtask.Title,
		"Status": subtask.Status,
	})
}

// Mettre à jour un commentaire dans la base de données.
//
// La fonction prend une connexion de base de données GORM et une structure de
// commentaire. Elle mettra à jour le commentaire dans la base de données avec les
// valeurs.presentes dans la structure de commentaire.
//
// La fonction ne renvoie pas de valeur, mais provoquera une panique si l'opération
// de mise à jour  choue.
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
//
// La fonction prend une connexion de base de données GORM et un identifiant de
// tâche, et supprime la tâche correspondante de la base de données.
//
// La fonction ne renvoie pas de valeur, mais provoquera une panique si l'opération
// d'effacement  choue.
func DeleteTask(db *gorm.DB, id uint) {
	var task Task
	db.Delete(&task, id)
}

// Supprime une sous-tâche de la base de données.
//
// La fonction prend une connexion de base de données GORM et un identifiant de
// sous-tâche, et supprime la sous-tâche correspondante de la base de données.
//
// La fonction ne renvoie pas de valeur, mais provoquera une panique si l'opération
// d'effacement  choue.
func DeleteSubtask(db *gorm.DB, id uint) {
	var subtask Subtask
	db.Delete(&subtask, id)
}

// Supprime un commentaire de la base de données.
//
// La fonction prend une connexion de base de données GORM et un identifiant de
// commentaire, et supprime le commentaire correspondant de la base de données.
//
// La fonction ne renvoie pas de valeur, mais provoquera une panique si l'opération
// d'effacement  choue.
func DeleteComment(db *gorm.DB, id uint) {
	var comment Comment
	db.Delete(&comment, id)
}
