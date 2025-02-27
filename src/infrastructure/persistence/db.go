package persistence

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"time"
)

func CreateDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("data.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	db.Exec("PRAGMA journal_mode=WAL;")
	db.Exec("PRAGMA cache_size=10000;")
	db.Exec("PRAGMA synchronous=NORMAL;")

	if err != nil {
		log.Fatal("⚠️ Unable to connect to the database:", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("⚠️ Unable to retrieve SQL connection:", err)
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Schema migration
	db.AutoMigrate(&Task{}, &Subtask{}, &Comment{}, &Tag{}, &TaskTags{}, &TaskSubtasks{}, &TaskComments{}, &TaskTimeSpent{}, &User{})

	return db
}
