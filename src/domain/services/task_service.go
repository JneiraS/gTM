package services

import (
	"fmt"
	"time"

	"github.com/JneiraS/AMS/src/infrastructure/persistence"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const (
	millisecondesToMinutes = 60000000000
	EndTaskStatus          = "Done"
	ProgressComplete       = 100
	Pending                = "Pending"
	InProgress             = "In progress"
	Done                   = "Done"
	invalidRequestMsg      = "Invalid request"
	invalidTaskIDMsg       = "Invalid task ID"
	taskNotFoundMsg        = "Task not found"
	FIND_BY_PROJECT        = "project = ?"
)

// EndTask met à jour une tâche en la définissant comme terminée.
//
// La fonction prend une connexion de base de données GORM et un identifiant de
// tâche. Elle mettra à jour la tâche dans la base de données en fixant son
// statut à "Done", son pourcentage de progression à 100, et sa durée passée
// en la calculant en prenant la durée entre l'heure de début et l'heure de fin.
func EndTask(db *gorm.DB, taskID uint) {
	taskTimeSpent := persistence.GetTaskTimeSpent(db, taskID)
	taskTimeSpent.EndTime = time.Now()
	persistence.UpdateTaskTimeSpent(db, taskTimeSpent)

	task, err := persistence.GetTask(db, taskID)
	if err != nil {
		return
	}

	if task.Status == "Pending" && task.TimeSpent != 0 {
		task.Status = EndTaskStatus
		task.Progress = ProgressComplete
	}
	task.Status = EndTaskStatus
	task.Progress = ProgressComplete
	task.TimeSpent = task.TimeSpent + int(TimeSpent(db, taskTimeSpent))/millisecondesToMinutes
	persistence.UpdateTask(db, task)
}

// TimeSpent prend une connexion à la base de données GORM et un ID de tâche.
// Récupère la tâche depuis la base de données et renvoie la durée entre
// l'heure de début et l'heure de fin.
func TimeSpent(db *gorm.DB, taskTimeSpent persistence.TaskTimeSpent) time.Duration {

	return taskTimeSpent.EndTime.Sub(taskTimeSpent.StartTime)
}

func UpdateStartTime(db *gorm.DB, taskID uint) {
	taskTimeSpent := persistence.GetTaskTimeSpent(db, taskID)
	taskTimeSpent.StartTime = time.Now()
	persistence.UpdateTaskTimeSpent(db, taskTimeSpent)
}

func IncrementTimeSpent(db *gorm.DB, taskID uint) int {
	taskTimeSpent := persistence.GetTaskTimeSpent(db, taskID)

	return int(time.Now().Sub(taskTimeSpent.StartTime)) / millisecondesToMinutes
}

// getStatistics returns an array of strings which represent the statistics
// of the tasks. The first element is the total number of tasks to do,
// and the next elements are the number of tasks with a specific
// characteristic.
func GetStatistics(priorityTasks, taskWithoutDueDate, lateTasks []persistence.Task) []string {
	taskCounts := map[string]int{
		"Tasks with priority": len(priorityTasks),
		"Without Due Date":    len(taskWithoutDueDate),
		"Overdue":             len(lateTasks),
	}

	totalTasks := 0
	for _, count := range taskCounts {
		totalTasks += count
	}

	statistics := make([]string, 0, len(taskCounts)+1)
	statistics = append(statistics, fmt.Sprintf("Total number of tasks to do: %d", totalTasks))

	for name, count := range taskCounts {
		statistics = append(statistics, fmt.Sprintf("%s: %d", name, count))
	}

	return statistics
}

func GetUserIDFromContext(c *gin.Context) int {
	claims, _ := c.Get("userID")
	return int(claims.(float64))
}

// getTasksByCategory retourne quatre groupe de tâches et un nom de projet.
func GetTasksByCategory(c *gin.Context, db *gorm.DB) ([]persistence.Task, []persistence.Task, []persistence.Task, []persistence.Task, string, error) {
	var priorityTasks []persistence.Task
	var taskWithoutDueDate []persistence.Task
	var tasksDone []persistence.Task
	var lateTasks []persistence.Task

	project := c.Param("project")

	baseQuery := db.Where("due_date > ? AND due_date > ? AND status != ?",
		time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Now(),
		Done)

	baseNoDueDateQuery := db.Where("due_date < ? AND status != ?",
		time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC),
		Done)

	if project != "" {
		baseQuery = baseQuery.Where(FIND_BY_PROJECT, project)
		baseNoDueDateQuery = baseNoDueDateQuery.Where(FIND_BY_PROJECT, project)
		if err := db.Where("status = ? AND project = ?", Done, project).
			Order("updated_at ASC").
			Find(&tasksDone).Error; err != nil {
			return nil, nil, nil, nil, "", err
		}

		if err := baseQuery.Where("due_date > ? AND due_date < ? AND status != ?",
			time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Now(),
			Done).
			Find(&lateTasks).Error; err != nil {
			return nil, nil, nil, nil, "", err
		}
	} else {
		if err := db.Where("status = ? AND updated_at > ? AND updated_at < ?",
			Done,
			time.Now().AddDate(0, 0, -1),
			time.Now()).
			Order("updated_at ASC").
			Find(&tasksDone).Error; err != nil {
			return nil, nil, nil, nil, "", err
		}

		if err := db.Where("due_date > ? AND due_date < ? AND status != ?",
			time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Now(),
			Done).
			Find(&lateTasks).Error; err != nil {
			return nil, nil, nil, nil, "", err
		}
	}

	if err := baseQuery.Order("priority ASC, due_date ASC").Find(&priorityTasks).Error; err != nil {
		return nil, nil, nil, nil, "", err
	}
	if err := baseNoDueDateQuery.Order("priority ASC").Find(&taskWithoutDueDate).Error; err != nil {
		return nil, nil, nil, nil, "", err
	}
	return priorityTasks, taskWithoutDueDate, tasksDone, lateTasks, project, nil
}
