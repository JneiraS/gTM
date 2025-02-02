package views

import (
	// "time"
	"net/http"
	"time"

	"github.com/JneiraS/AMS/src/domain/models"
	"github.com/JneiraS/AMS/src/infrastructure/persistence"
	"github.com/gin-gonic/gin"
)

func Views() {
	router := gin.Default()
	router.LoadHTMLGlob("src/infrastructure/web/templates/*")
	router.Static("/static", "src/infrastructure/web/static")

	router.GET("/", func(c *gin.Context) {
		db := persistence.CreateDB()
		var allTasks []persistence.Task
		db.Find(&allTasks)

		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Liste des taches",
			"tasks": allTasks,
		})
	})

	router.POST("/", func(c *gin.Context) {
		db := persistence.CreateDB()
		task := persistence.Task{
			Task: models.Task{
				Title:       c.PostForm("title"),
				Description: c.PostForm("description"),
				DueDate:     time.Now(),
				Status:      c.PostForm("status"),
				Priority:    c.PostForm("priority"),
				Assignee:    c.PostForm("assignee"),
				Creator:     c.PostForm("creator"),
				Project:     c.PostForm("project"),
				Progress:    0,
			},
		}
		persistence.CreateTask(db, task)

		c.Redirect(http.StatusFound, "/")
	})

	router.Run(":8080")
}
