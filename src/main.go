package main

import (
	// "time"
	"github.com/JneiraS/AMS/src/domain/models"
	// "github.com/JneiraS/AMS/src/domain/services"
	"github.com/JneiraS/AMS/src/infrastructure/persistence"
)

func main() {
	//api.StartServer()

	db := persistence.CreateDB()

	persistence.CreateComment(db, persistence.Comment{Comment: models.Comment{
		Author: "John Doe", Text: "Test Comment",
	}})

	// services.StartTask(db, tastk)

	// fmt.Println(tastk.Title)

	// metter a jour la tache
	// var taskToUpdate persistence.Task
	// db.First(&taskToUpdate, 1)

	// taskToUpdate.Title = "Test update"

	// persistence.UpdateTask(db, taskToUpdate)

	// tastku := persistence.GetTask(db, 1)
	// fmt.Println(tastku.Title)

	// persistence.CreateTask(db, persistence.Task{Task: models.Task{
	// 	Title:         "Test",
	// 	Description:   "2y7v0@example.com",
	// 	DueDate:       time.Now(),
	// 	Status:        "En cours",
	// 	Priority:      "Moyenne",
	// 	Assignee:      "Jean",
	// 	Creator:       "Jean",
	// 	Project:       "Go",
	// 	Progress:      50,
	// 	EstimatedTime: 2,
	// 	// TimeSpent:     1,
	// }})

	// persistence.CreateSubtask(db, persistence.Subtask{Subtask: models.Subtask{
	// 	Title:  "Test",
	// 	Status: "En cours",
	// }})

}
