package services

import (
	"time"

	"github.com/JneiraS/AMS/src/infrastructure/persistence"
	"gorm.io/gorm"
)

const (
	millisecondesToMinutes = 60000000000
	EndTaskStatus          = "Done"
	ProgressComplete       = 100
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

	db.Save(&taskTimeSpent)

	task, err := persistence.GetTask(db, taskID)
	if err != nil {
		return
	}
	task.Status = EndTaskStatus
	task.Progress = ProgressComplete
	task.TimeSpent = int(TimeSpent(db, taskID)) / millisecondesToMinutes
	persistence.UpdateTask(db, task)
}

// TimeSpent prend une connexion à la base de données GORM et un ID de tâche.
// Récupère la tâche depuis la base de données et renvoie la durée entre
// l'heure de début et l'heure de fin.
func TimeSpent(db *gorm.DB, taskID uint) time.Duration {
	taskTimeSpent := persistence.GetTaskTimeSpent(db, taskID)

	if taskTimeSpent.StartTime.IsZero() {
		task, err := persistence.GetTask(db, taskID)
		if err == nil {
			taskTimeSpent.StartTime = task.CreatedAt
		}
	}
	return taskTimeSpent.EndTime.Sub(taskTimeSpent.StartTime)
}
