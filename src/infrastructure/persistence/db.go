package persistence

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// CreateDB creates a new GORM database connection using SQLite3 and migrates the Task
// struct to the database. This function will panic if the database connection
// cannot be established.
func CreateDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("data.db"), &gorm.Config{})
	db.Set("gorm:table_options", "ENGINE=InnoDB")
	// Migrations
	db.AutoMigrate(&Task{}, &Subtask{}, &Comment{}, &Tag{}, &TaskTags{}, &TaskSubtasks{}, &TaskComments{}, &TaskTimeSpent{})
	if err != nil {
		panic("failed to connect database")
	}
	return db
}
