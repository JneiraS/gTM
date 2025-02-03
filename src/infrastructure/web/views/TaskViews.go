package views

import (
	// "time"
	"net/http"
	"time"

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

	router.GET("/", func(c *gin.Context) {
		db := persistence.CreateDB()
		var priorityTass []persistence.Task
		var taskWithoutDueDate []persistence.Task
		var tasksDone []persistence.Task
		var lateTasks []persistence.Task

		db.Where("due_date > ? AND due_date < ? AND status != ?", time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC), time.Now(), "done").
			Order("priority ASC, due_date ASC").
			Find(&priorityTass)
		db.Where("due_date < ? AND status != ?", time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC), "done").
			Order("priority ASC").Find(&taskWithoutDueDate)
		db.Where("status = ?", "done").Find(&tasksDone)
		db.Where("due_date > ? AND due_date < ? AND status != ?", time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC), time.Now(), "done").Find(&lateTasks)

		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title":           "Liste des taches",
			"prioritytasks":   priorityTass,
			"taskswithoutdue": taskWithoutDueDate,
			"tasksdone":       tasksDone,
			"latetasks":       lateTasks,
		})
	})

	router.GET("/done/:id", useCases.UpdateTaskHandler(persistence.CreateDB()))
	router.POST("/", useCases.CreateTaskHandler(persistence.CreateDB()))

	router.Run(":7263")

}
