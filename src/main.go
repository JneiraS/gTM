package main

import "github.com/JneiraS/AMS/src/infrastructure/web/views"

//"github.com/JneiraS/AMS/src/domain/models"
//"github.com/JneiraS/AMS/src/infrastructure/persistence"
//"time"

func main() {
	views.Views()
	// db := persistence.CreateDB()

	// persistence.CreateTask(db, persistence.Task{
	// 	Task: models.Task{
	// 		Title:         "Test GORM",
	// 		Description:   "Test Description",
	// 		DueDate:       time.Now(),
	// 		Status:        "pending",
	// 		Priority:      "low",
	// 		Assignee:      "test-user",
	// 		Creator:       "admin",
	// 		Project:       "test-project",
	// 		Progress:      0,
	// 		EstimatedTime: 2,
	// 	},
	// })
}
