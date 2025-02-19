package views

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

	sqliteDB := persistence.CreateDB()

	// Toutes les 5 minutes

	router.GET("/", DisplayTasks(sqliteDB))
	router.GET("/done/:id", MarkTaskAsDoneHandler(sqliteDB))
	router.POST("/", CreateTaskHandler(sqliteDB))
	router.POST("/update-task", updateTaskHandler(sqliteDB))
	router.POST("/update-task-title", updateTitleTaskHandler(sqliteDB))
	router.GET("/project/:project", DisplayTasks(sqliteDB))
	router.POST("/update-task-due-date", updateDueDateHandler(sqliteDB))
	router.POST("/update-task-status", updateStatusHandler(sqliteDB))
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
func MarkTaskAsDoneHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		db := persistence.CreateDB()
		services.EndTask(db, uint(id))
		c.Redirect(http.StatusFound, "/")
	}
}

// DisplayTasks returns a gin.HandlerFunc that displays the main page of the
// application. The main page displays 4 lists of tasks: tasks with a deadline,
// tasks without a deadline, tasks done in the last 24 hours and tasks that are
// late.
//
// The tasks are ordered by priority and deadline.
//
// The function takes a gorm.DB as a parameter, which is used to query the
// database.
//
// The function returns a gin.HandlerFunc that will be called when a request is
// made to the main page. The gin.HandlerFunc will render the main page with the
// tasks in the 4 lists.
func DisplayTasks(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		db := persistence.CreateDB()

		var priorityTasks []persistence.Task
		var taskWithoutDueDate []persistence.Task
		var tasksDone []persistence.Task
		var lateTasks []persistence.Task

		project := c.Param("project")

		baseQuery := db.Where("due_date > ? AND due_date > ? AND status != ?",
			time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Now(),
			Done)

		baseNoDueDateQuery := db.Where("due_date < ? AND status != ?",
			time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC),
			Done)

		if project != "" {
			baseQuery = baseQuery.Where("project = ?", project)
			baseNoDueDateQuery = baseNoDueDateQuery.Where("project = ?", project)
			db.Where("status = ? AND project = ?", Done, project).
				Order("updated_at ASC").
				Find(&tasksDone)
		} else {
			db.Where("status = ? AND updated_at > ? AND updated_at < ?",
				Done,
				time.Now().AddDate(0, 0, -1),
				time.Now()).
				Order("updated_at ASC").
				Find(&tasksDone)
		}

		baseQuery.Order("priority ASC, due_date ASC").Find(&priorityTasks)
		baseNoDueDateQuery.Order("priority ASC").Find(&taskWithoutDueDate)

		db.Where("due_date > ? AND due_date < ? AND status != ?",
			time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Now(),
			Done).
			Find(&lateTasks)

		// Affichage de la vue avec toutes les données
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title":           "Liste des taches",
			"projects":        persistence.GetAllProjects(db),
			"prioritytasks":   priorityTasks,
			"taskswithoutdue": taskWithoutDueDate,
			"tasksdone":       tasksDone,
			"latetasks":       lateTasks,
			"project":         project,
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

// updateDueDateHandler updates a task's due date and redirects to the main page.
//
// The task ID to update must be given in the request body as a JSON object with the key
// "id". The request body must also contain the new due date in the format
// "YYYY-MM-DD HH:mm".
func updateDueDateHandler(db *gorm.DB) gin.HandlerFunc {
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

// updateStatusHandler updates a task's status to "In progress" if the task is currently
// "Pending", or updates a task's status to "Pending" if the task is currently "In progress".
//
// The task ID to update must be given in the request body as a JSON object with the key
// "id". The request body must also contain the current status of the task.
func updateStatusHandler(db *gorm.DB) gin.HandlerFunc {
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
			task.Status = "In progress"
			services.UpdateStartTime(db, task.ID)
		} else if task.Status == "In progress" {
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
