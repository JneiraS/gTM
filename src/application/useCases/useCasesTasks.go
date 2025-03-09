package useCases

import (
	"net/http"
	"strconv"
	"time"

	"github.com/JneiraS/AMS/src/domain/models"
	"github.com/JneiraS/AMS/src/domain/services"
	"github.com/JneiraS/AMS/src/infrastructure/persistence"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const (
	Pending           = "Pending"
	InProgress        = "In progress"
	Done              = "Done"
	invalidRequestMsg = "Invalid request"
	invalidTaskIDMsg  = "Invalid task ID"
	taskNotFoundMsg   = "Task not found"
)

// CreateTaskHandler handles the creation of a new task and redirects to the main page.
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
		if task.Project != "" {
			c.Redirect(http.StatusFound, "/project/"+task.Project)
		} else {
			c.Redirect(http.StatusFound, "/")

		}
	}
}

// ToggleTaskStatusHandler toggles the status of a task and updates the time spent.
func ToggleTaskStatusHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var requestBody struct {
			ID string `json:"id"`
		}

		if err := c.BindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": invalidRequestMsg, "status": 400})
			return
		}

		id, err := strconv.Atoi(requestBody.ID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": invalidTaskIDMsg, "status": 400})
			return
		}

		task := persistence.Task{}
		if err := db.First(&task, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": taskNotFoundMsg, "status": 404})
			return
		}

		if task.Status == "Pending" || task.Status == "Stopped" {
			task.Status = InProgress
			services.UpdateStartTime(db, task.ID)
		} else if task.Status == InProgress {
			task.Status = "Stopped"
			task.TimeSpent = task.TimeSpent + services.IncrementTimeSpent(db, task.ID)
		}
		if err := db.Save(&task).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task status", "status": 500})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Task status updated successfully", "status": 200})
	}
}

// MarkTaskAsCompletedHandler updates a task to done and redirects to the main page.
func MarkTaskAsCompletedHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		services.EndTask(db, uint(id))
		c.Redirect(http.StatusFound, "/")
	}
}

// updateTaskHandler updates a task's description and returns a HTTP 200 OK status.
func UpdateTaskDescriptionHandler(db *gorm.DB) gin.HandlerFunc {
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
func UpdateTitleTaskHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var update struct {
			ID    string `json:"id"`
			Title string `json:"title"`
		}

		if err := c.BindJSON(&update); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": invalidRequestMsg})
			return
		}

		id, err := strconv.Atoi(update.ID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": invalidTaskIDMsg})
			return
		}

		task := persistence.Task{}
		if err := db.First(&task, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": taskNotFoundMsg})
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

// UpdateDueDateHandler updates a task's due date and redirects to the main page.
func UpdateDueDateHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var requestBody struct {
			ID   string `json:"id"`
			Date string `json:"date"`
		}

		if err := c.BindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": invalidRequestMsg})
			return
		}

		id, err := strconv.Atoi(requestBody.ID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": invalidTaskIDMsg})
			return
		}
		dueDate, err := time.Parse("2006-01-02 15:04", requestBody.Date)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid due date format. Expected format: YYYY-MM-DD HH:mm"})
			return
		}

		// Add validation for future dates
		if dueDate.Before(time.Now().Add(-time.Minute)) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":  "Due date must be in the future",
				"status": 400,
			})
			return

		}
		task := persistence.Task{}
		if err := db.First(&task, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": taskNotFoundMsg})
			return
		}

		task.DueDate = dueDate
		if err := db.Save(&task).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "Task due date updated successfully"})
	}
}

func CreateCommentHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		username, _ := c.Cookie("username")
		id, _ := strconv.Atoi(c.Param("id"))

		comment := persistence.Comment{
			Comment: models.Comment{
				Author: username,
				Text:   c.PostForm("comment")},
		}

		persistence.CreateComment(db, comment, uint(id))

		c.Status(http.StatusCreated)
		c.Redirect(http.StatusFound, "/task/"+c.Param("id"))
	}
}
