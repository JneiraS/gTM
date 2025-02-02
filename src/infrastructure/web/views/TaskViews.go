package views

import (
	// "time"
	"net/http"
	"strconv"
	"time"

	"github.com/JneiraS/AMS/src/application/useCases"
	"github.com/JneiraS/AMS/src/infrastructure/persistence"
	"github.com/gin-gonic/gin"
)

func Views() {
	router := gin.Default()
	router.LoadHTMLGlob("src/infrastructure/web/templates/*")
	router.Static("/static", "src/infrastructure/web/static")

	router.GET("/", func(c *gin.Context) {
		db := persistence.CreateDB()
		var priorityTass []persistence.Task
		var taskWithoutDueDate []persistence.Task

		db.Where("due_date > ? AND due_date < ?", time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC), time.Now()).
			Order("priority ASC, due_date ASC").
			Find(&priorityTass)

		db.Where("due_date < ?", time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)).
			Order("priority ASC").Find(&taskWithoutDueDate)

		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title":           "Liste des taches",
			"prioritytasks":   priorityTass,
			"taskswithoutdue": taskWithoutDueDate,
		})
	})
	router.GET("/done/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		db := persistence.CreateDB()
		task := persistence.Task{}
		db.Model(&task).Where("id = ?", id).Update("status", "done")
		c.Redirect(http.StatusFound, "/")
	})

	router.POST("/", useCases.CreateTaskHandler(persistence.CreateDB()))

	router.Run(":8080")
}
