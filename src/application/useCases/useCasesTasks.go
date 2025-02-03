package useCases

import (
	"github.com/JneiraS/AMS/src/domain/models"
	"github.com/JneiraS/AMS/src/infrastructure/persistence"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"time"
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

func MainRender(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		db := persistence.CreateDB()
		var priorityTass []persistence.Task
		var taskWithoutDueDate []persistence.Task
		var tasksDone []persistence.Task
		var lateTasks []persistence.Task

		db.Where("due_date > ? AND due_date > ? AND status != ?", time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC), time.Now(), "done").
			Order("priority ASC, due_date ASC").
			Find(&priorityTass)
		db.Where("due_date < ? AND status != ?", time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC), "done").
			Order("priority ASC").Find(&taskWithoutDueDate)
		db.Where("status = ? AND updated_at > ? AND updated_at < ?", "done", time.Now().AddDate(0, 0, -1), time.Now()).Order("updated_at ASC").Find(&tasksDone)
		db.Where("due_date > ? AND due_date < ? AND status != ?", time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC), time.Now(), "done").Find(&lateTasks)

		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title":           "Liste des taches",
			"prioritytasks":   priorityTass,
			"taskswithoutdue": taskWithoutDueDate,
			"tasksdone":       tasksDone,
			"latetasks":       lateTasks,
		})
	}
}
