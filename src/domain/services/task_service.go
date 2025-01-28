package services

import (
	"time"

	"github.com/JneiraS/AMS/src/infrastructure/persistence"
	"gorm.io/gorm"
)

// StartTask démarre une tâche et enregistre le temps de départ.
//
// La fonction prend une connexion de base de données GORM et une structure de
// tâche. Elle met à jour la tâche dans la base de données en ajoutant le temps
// de départ actuel.
func StartTask(db *gorm.DB, task persistence.Task) {
	var taskSpent persistence.TaskTimeSpent
	// Récupérer la tâche à partir de la base de données
	db.First(&task)
	// Ajouter le temps de début de tâche
	taskSpent.StartTime = time.Now()
	taskSpent.TaskID = uint(task.ID)
	// Mettre à jour la tâche dans la base de données
	db.Create(&taskSpent)
}
