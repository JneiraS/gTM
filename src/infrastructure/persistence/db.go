package persistence

import (
	"log"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// CreateDB initialise la connexion SQLite avec une meilleure gestion de concurrence
func CreateDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("data.db"), &gorm.Config{
		// Mode silencieux pour éviter trop de logs
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatal("⚠️ Impossible de se connecter à la base de données :", err)
	}

	// Configuration SQLite pour meilleure gestion des accès concurrents
	db.Exec("PRAGMA journal_mode=WAL;")             // Active WAL pour meilleure concurrence
	db.Exec("PRAGMA synchronous=NORMAL;")           // Réduit la latence des écritures
	db.Exec("PRAGMA cache_size=-10000;")            // Définit le cache en nombre de pages
	db.Exec("PRAGMA busy_timeout=5000;")            // Attente si la base est verrouillée
	db.Exec("PRAGMA journal_size_limit=100000000;") // Limite la taille du journal WAL

	// Récupération de la connexion SQL pour configurer le pool de connexions
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("⚠️ Impossible de récupérer la connexion SQL :", err)
	}

	// Pool de connexions optimisé pour SQLite
	sqlDB.SetMaxOpenConns(5)                   // Limite le nombre total de connexions ouvertes
	sqlDB.SetMaxIdleConns(2)                   // Nombre max de connexions inactives
	sqlDB.SetConnMaxLifetime(time.Minute * 30) // Durée max de vie d'une connexion

	// Migration du schéma
	db.AutoMigrate(&Task{},
		&Subtask{}, &Comment{},
		&Tag{}, &TaskTags{},
		&TaskSubtasks{},
		&TaskComments{},
		&TaskTimeSpent{},
		&User{})

	return db
}
