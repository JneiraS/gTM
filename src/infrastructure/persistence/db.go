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
	// Create the file-based database
	db, err := gorm.Open(sqlite.Open("data.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		panic("failed to get database")
	}
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Schema migration
	db.AutoMigrate(&Task{}, &Subtask{}, &Comment{}, &Tag{}, &TaskTags{}, &TaskSubtasks{}, &TaskComments{}, &TaskTimeSpent{}, &User{})

	return db
}
