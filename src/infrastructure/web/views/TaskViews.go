package views

import (
	// "time"
	"net/http"

	"github.com/JneiraS/AMS/src/infrastructure/persistence"
	"github.com/gin-gonic/gin"
)

func Views() {
	router := gin.Default()
	router.LoadHTMLGlob("src/infrastructure/web/templates/*")
	router.Static("/static", "src/infrastructure/web/static")

	router.POST("/", func(c *gin.Context) {

		db := persistence.CreateDB()
		var allTasks []persistence.Task
		db.Find(&allTasks)

		// name := c.PostForm("name")
		// email := c.PostForm("email")
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Liste des taches",
			"tasks": allTasks,
		})
	})

	// router.GET("/", func(c *gin.Context) {
	// 	c.HTML(http.StatusOK, "index.tmpl", gin.H{
	// 		"title": "Main websites",
	// 	})
	// })

	router.GET("/indexx", func(c *gin.Context) {
		c.HTML(http.StatusOK, "indexx.tmpl", gin.H{
			"title": "Main TEST",
		})
	})
	router.Run(":8080")
}
