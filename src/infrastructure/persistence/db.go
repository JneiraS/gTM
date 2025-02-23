package persistence

import (
	"os"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// CreateDB creates a new GORM database connection using SQLite3 and migrates the Task
// struct to the database. This function will panic if the database connection
// cannot be established.
func CreateDB() *gorm.DB {
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

	// Schema migration
	db.AutoMigrate(&Task{}, &Subtask{}, &Comment{}, &Tag{}, &TaskTags{}, &TaskSubtasks{}, &TaskComments{}, &TaskTimeSpent{}, &User{})

	// Attach disk database and copy data
	if _, err := os.Stat("data.db"); !os.IsNotExist(err) {
		db.Exec("ATTACH DATABASE 'data.db' AS disk")
		tables := []string{"tasks", "subtasks", "comments", "tags", "task_tags", "task_subtasks", "task_comments", "task_time_spents", "users"}
		for _, table := range tables {
			db.Exec("INSERT INTO main." + table + " SELECT * FROM disk." + table)
		}
	}

	// Periodic background save
	go func() {
		for {
			time.Sleep(1 * time.Minute)
			db.Transaction(func(tx *gorm.DB) error {
				tables := []string{"tasks", "subtasks", "comments", "tags", "task_tags", "task_subtasks", "task_comments", "task_time_spents", "users"}
				for _, table := range tables {
					if err := tx.Exec("DELETE FROM disk." + table).Error; err != nil {
						return err
					}
					if err := tx.Exec("INSERT INTO disk." + table + " SELECT * FROM main." + table).Error; err != nil {
						return err
					}
				}
				return tx.Exec("PRAGMA disk.vacuum").Error
			})
		}
	}()

	return db
}
