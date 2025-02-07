package views

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/JneiraS/AMS/src/domain/models"
	"github.com/JneiraS/AMS/src/infrastructure/persistence"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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

	// Spécifier les adresses IP ou plages autorisées
	router.SetTrustedProxies([]string{"192.168.1.2", "10.0.0.0/8"})
	router.LoadHTMLGlob("src/infrastructure/web/templates/*")
	router.Static("/static", "src/infrastructure/web/static")

	router.GET("/", MainRender(persistence.CreateDB()))
	router.GET("/done/:id", UpdateStatusHandler(persistence.CreateDB()))
	router.POST("/", CreateTaskHandler(persistence.CreateDB()))
	router.POST("/update-task", updateTaskHandler(persistence.CreateDB()))
	router.POST("/update-task-title", updateTitleTaskHandler(persistence.CreateDB()))
	//updateTitleTaskHandler

	router.Run(":7263")

}

// CreateTaskHandler creates a new task and redirects to the main page.
//
// The due date must be in the format "2006-01-02T15:04". The fields title,
// description, status, priority, assignee, creator and project must be filled.
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

// UpdateStatusHandler updates a task to done and redirects to the main page.
//
// The task ID to update must be given as a parameter in the URL path.
// The task is updated to have the status "done".
func UpdateStatusHandler(db *gorm.DB) gin.HandlerFunc {
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

// MainRender renders the main page of the application.
//
// The handler renders the main page of the application. It queries the database
// to get the list of tasks and passes it to the template.
//
// The tasks are ordered by priority and then by due date.
// The tasks are divided into four categories:
// - priorityTass: tasks that are not due yet and are not done.
// - taskWithoutDueDate: tasks that are not due yet and do not have a due date.
// - tasksDone: tasks that are done and were updated in the last 24 hours.
// - lateTasks: tasks that are not done and are due.
//
// The handler returns a HTML response with the given template and data.
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

// updateTaskHandler updates a task's description and returns a HTTP 200 OK status.
//
// The task ID to update must be given as a parameter in the URL path.
// The task is updated to have the given description.
//
// The handler returns a HTTP 200 OK response if the update is successful.
// The handler returns a HTTP 400 Bad Request response if the request is invalid.
// The handler returns a HTTP 500 Internal Server Error response if the update fails.
func updateTaskHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var update struct {
			ID          string `json:"id"`
			Description string `json:"description"`
		}

		if err := c.BindJSON(&update); err != nil {
			c.Status(http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(update.ID)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}

		task := persistence.Task{}
		db.First(&task, id)

		task.Description = update.Description
		if err := db.Save(&task).Error; err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Status(http.StatusOK)
	}
}

// UpdateTitleTaskHandler updates a task's title and redirects to the main page.
//
// The task ID to update must be given as a parameter in the URL path.
// The task is updated to have the given title.
func updateTitleTaskHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var update struct {
			ID    string `json:"id"`
			Title string `json:"title"`
		}

		if err := c.BindJSON(&update); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		id, err := strconv.Atoi(update.ID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
			return
		}

		task := persistence.Task{}
		if err := db.First(&task, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			return
		}

		task.Title = update.Title
		if err := db.Save(&task).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "Task title updated successfully"})
	}
}
