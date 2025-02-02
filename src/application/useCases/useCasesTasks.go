package useCases

import (
	"net/http"
	"strconv"
	"time"

	"github.com/JneiraS/AMS/src/domain/models"
	"github.com/JneiraS/AMS/src/infrastructure/persistence"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// createTaskHandler crée une nouvelle tâche à partir des données du formulaire
func CreateTaskHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		dueDate, _ := time.Parse("2006-01-02T15:04", c.PostForm("due_date"))
		task := persistence.Task{
			Task: models.Task{
				Title:       c.PostForm("title"),
				Description: c.PostForm("description"),
				DueDate:     dueDate,
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
	}
}

// updateTaskHandler met à jour une tâche dans la base de données
func UpdateTaskHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		db := persistence.CreateDB()
		task := persistence.Task{}
		db.First(&task, id)
		task.Status = "done"
		persistence.UpdateTask(db, task)
		c.Redirect(http.StatusFound, "/")
	}
}
