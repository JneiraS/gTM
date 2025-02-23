package persistence

import (
	"log"
	"sync"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	db   *gorm.DB
	once sync.Once
)

// CreateDB initialise une seule instance de la base de données
func CreateDB() *gorm.DB {
	once.Do(func() {
		// Configuration de la base de données en mémoire
		dsn := "file::memory:?cache=shared" +
			"&_pragma=foreign_keys(1)" +
			"&_pragma=busy_timeout(10000)" +
			"&_pragma=synchronous(NORMAL)" +
			"&_pragma=journal_mode=WAL" +
			"&_pragma=cache_size(-2000000)" +
			"&_pragma=mmap_size(2147483648)" +
			"&_pragma=temp_store(MEMORY)" +
			"&_pragma=threads(4)" +
			"&_pragma=page_size(4096)"

		var err error
		db, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{
			SkipDefaultTransaction: true,
			PrepareStmt:            true,
		})
		if err != nil {
			log.Fatal("Erreur de connexion à la base de données:", err)
		}

		sqlDB, err := db.DB()
		if err != nil {
			log.Fatal("Erreur lors de l'obtention de sqlDB:", err)
		}

		sqlDB.SetMaxOpenConns(25)
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetConnMaxLifetime(time.Hour)

		// Attacher la base persistante
		db.Exec("ATTACH DATABASE 'data.db' AS disk")

		// Migration des tables
		err = db.AutoMigrate(&Task{}, &Subtask{}, &Comment{}, &Tag{}, &TaskTags{}, &TaskSubtasks{}, &TaskComments{}, &TaskTimeSpent{})
		if err != nil {
			log.Fatal("Erreur lors de la migration:", err)
		}

		// Synchronisation initiale des données
		err = db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Exec("INSERT INTO main.tasks SELECT * FROM disk.tasks WHERE id NOT IN (SELECT id FROM main.tasks)").Error; err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			log.Fatal("Erreur lors de la synchronisation des données:", err)
		}

		// Sauvegarde périodique en arrière-plan
		go func() {
			for {
				time.Sleep(1 * time.Minute)
				err := db.Transaction(func(tx *gorm.DB) error {
					tables := []string{"tasks", "subtasks", "comments", "tags", "task_tags", "task_subtasks", "task_comments", "task_time_spents"}
					for _, table := range tables {
						if err := tx.Exec("DELETE FROM disk." + table).Error; err != nil {
							return err
						}
						if err := tx.Exec("INSERT INTO disk." + table + " SELECT * FROM main." + table).Error; err != nil {
							return err
						}
					}
					// Nettoyage et optimisation de la base sur disque
					return tx.Exec("VACUUM disk").Error
				})
				if err != nil {
					log.Println("Erreur lors de la sauvegarde automatique:", err)
				}
			}
		}()
	})

	return db
}

// GetDB retourne l'instance unique de la base de données
func GetDB() *gorm.DB {
	if db == nil {
		return CreateDB()
	}
	return db
}
