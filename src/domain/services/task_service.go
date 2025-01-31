package services

import (
	"time"

	"github.com/JneiraS/AMS/src/infrastructure/persistence"
	"gorm.io/gorm"
)

// Prend une connexion à la base de données GORM et un ID de tâche.
// Récupère la tâche depuis la base de données et la met à jour en
// ajoutant l'heure de début actuelle.
func StartTask(db *gorm.DB, id uint) {
	var taskSpent persistence.TaskTimeSpent

	// Récupère la tâche depuis la base de données
	task, _ := persistence.GetTask(db, id)
	// Crée une entrée TaskTimeSpent avec l'heure de début actuelle et l'ID de la tâche
	taskSpent.StartTime = time.Now()
	taskSpent.TaskID = uint(task.ID)
	// Insère l'enregistrement TaskTimeSpent dans la base de données
	db.Create(&taskSpent)
}
