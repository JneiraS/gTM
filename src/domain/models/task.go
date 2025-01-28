package models

import "time"

// Task représente une tâche.
type Task struct {
	Title       string    `json:"title"`       // Titre ou nom de la tâche
	Description string    `json:"description"` // Description détaillée
	DueDate     time.Time `json:"due_date"`    // Date limite ou échéance
	Status      string    `json:"status"`      // Statut de la tâche (par ex., "En cours", "Terminée")
	Priority    string    `json:"priority"`    // Priorité (par ex., "Basse", "Moyenne", "Haute")
	Assignee    string    `json:"assignee"`    // Personne à qui la tâche est assignée
	Creator     string    `json:"creator"`     // Personne qui a créé la tâche
	Project     string    `json:"project"`     // Nom du projet ou catégorie
	// Attachments   []string  `json:"attachments"`    // Liste des fichiers joints
	Progress      int `json:"progress"`       // Progression en pourcentage (0-100)
	EstimatedTime int `json:"estimated_time"` // Temps estimé pour terminer (en heures)
	TimeSpent     int `json:"time_spent"`     // Temps réellement passé (en heures)
}

// Subtask représente une sous-tâche.
type Subtask struct {
	Title  string `json:"title"`  // Nom de la sous-tâche
	Status string `json:"status"` // Statut de la sous-tâche
}

// Comment représente un commentaire sur une tâche.
type Comment struct {
	Author string `json:"author"` // Auteur du commentaire
	Text   string `json:"text"`   // Contenu du commentaire
}

type Tag struct {
	Name string `gorm:"uniqueIndex"` // Nom unique du tag
}

type TimeSpent struct {
	StartTime time.Time `json:"start_time"` // Début de l'intervalle de temps
	EndTime   time.Time `json:"end_time"`   // Fin de l'intervalle de temps
}
