package persistence

import (
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// CreateDB creates a new GORM database connection using SQLite3 and migrates the Task
// struct to the database. This function will panic if the database connection
// cannot be established.
func CreateDB() *gorm.DB {

	dbb, err := gorm.Open(sqlite.Open("data.db"), &gorm.Config{})
	dbb.Set("gorm:table_options", "ENGINE=InnoDB")
	// Migrations
	dbb.AutoMigrate(&Task{}, &Subtask{}, &Comment{}, &Tag{}, &TaskTags{}, &TaskSubtasks{}, &TaskComments{}, &TaskTimeSpent{})
	if err != nil {
		panic("failed to connect database")
	}

	dsn := "file::memory:?cache=shared" +
		"&_pragma=foreign_keys(1)" +
		"&_pragma=busy_timeout(10000)" +
		"&_pragma=synchronous(NORMAL)" +
		"&_pragma=journal_mode(WAL)" +
		"&_pragma=cache_size(-2000000)" +
		"&_pragma=mmap_size(2147483648)" +
		"&_pragma=temp_store(MEMORY)" +
		"&_pragma=threads(4)" +
		"&_pragma=page_size(4096)"

	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	})
	if err != nil {
		panic("failed to connect database")
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic("failed to get database")
	}
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Attach disk database
	db.Exec("ATTACH DATABASE 'data.db' AS disk")

	// Schema and initial sync
	db.AutoMigrate(&Task{}, &Subtask{}, &Comment{}, &Tag{}, &TaskTags{}, &TaskSubtasks{}, &TaskComments{}, &TaskTimeSpent{})

	db.Transaction(func(tx *gorm.DB) error {
		tx.Exec("INSERT INTO main.tasks SELECT * FROM disk.tasks WHERE id NOT IN (SELECT id FROM main.tasks)")
		return nil
	})
	// Après l'AutoMigrate et la création des tables, ajouter :

	// Sauvegarde périodique
	go func() {
		for {
			time.Sleep(1 * time.Minute)
			db.Transaction(func(tx *gorm.DB) error {
				tx.Exec("DELETE FROM disk.tasks")
				tx.Exec("INSERT INTO disk.tasks SELECT * FROM main.tasks")
				tx.Exec("DELETE FROM disk.subtasks")
				tx.Exec("INSERT INTO disk.subtasks SELECT * FROM main.subtasks")
				tx.Exec("DELETE FROM disk.comments")
				tx.Exec("INSERT INTO disk.comments SELECT * FROM main.comments")
				tx.Exec("DELETE FROM disk.tags")
				tx.Exec("INSERT INTO disk.tags SELECT * FROM main.tags")
				tx.Exec("DELETE FROM disk.task_tags")
				tx.Exec("INSERT INTO disk.task_tags SELECT * FROM main.task_tags")
				tx.Exec("DELETE FROM disk.task_subtasks")
				tx.Exec("INSERT INTO disk.task_subtasks SELECT * FROM main.task_subtasks")
				tx.Exec("DELETE FROM disk.task_time_spents")
				tx.Exec("INSERT INTO disk.task_time_spents SELECT * FROM main.task_time_spents")
				tx.Exec("DELETE FROM disk.task_comments")
				tx.Exec("INSERT INTO disk.task_comments SELECT * FROM main.task_comments")
				return nil
			})
		}
	}()

	return db
}
