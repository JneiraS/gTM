package views

import (
	// "time"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Views() {
	router := gin.Default()
	router.LoadHTMLGlob("src/infrastructure/web/templates/*")
	router.Static("/static", "src/infrastructure/web/static")

	router.POST("/", func(c *gin.Context) {
		name := c.PostForm("name")
		email := c.PostForm("email")
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Main website",
			"name":  name,
			"email": email,
		})
	})

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Main website",
		})
	})

	router.GET("/indexx", func(c *gin.Context) {
		c.HTML(http.StatusOK, "indexx.tmpl", gin.H{
			"title": "Main TEST",
		})
	})
	router.Run(":8080")
}
