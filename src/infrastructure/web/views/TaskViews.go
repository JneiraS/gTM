package views

import (
	"github.com/JneiraS/AMS/src/application/useCases"
	"github.com/JneiraS/AMS/src/infrastructure/persistence"
	"github.com/gin-gonic/gin"
)

// Views creates a web server that serves the web interface for the application.
//
// The web interface shows four lists of tasks: tasks with a due date, tasks without
// a due date, tasks that are done, and tasks that are late.
//
// The functions to create a new task, mark a task as done and update a task are
// routed to the corresponding functions in the useCases package.
func Views() {
	router := gin.Default()
	router.LoadHTMLGlob("src/infrastructure/web/templates/*")
	router.Static("/static", "src/infrastructure/web/static")

	router.GET("/", useCases.MainRender(persistence.CreateDB()))
	router.GET("/done/:id", useCases.UpdateTaskHandler(persistence.CreateDB()))
	router.POST("/", useCases.CreateTaskHandler(persistence.CreateDB()))

	router.Run(":7263")

}
