package views

import (
	"net/http"
	"time"

	"github.com/JneiraS/AMS/src/application/useCases"
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

// SetupRouter creates a web server that serves the web interface for the application.
//
// The web interface shows four lists of tasks: tasks with a due date, tasks without
// a due date, tasks that are done, and tasks that are late.
//
// The functions to create a new task, mark a task as done and update a task are
// routed to the corresponding functions in the useCases package.
func SetupRouter(db *gorm.DB) {
	router := gin.Default()

	// Spécifier les adresses IP ou plages autorisées
	router.SetTrustedProxies([]string{"192.168.1.2", "10.0.0.0/8"})
	router.LoadHTMLGlob("src/infrastructure/web/templates/*")
	router.Static("/static", "src/infrastructure/web/static")

	router.GET("/", DisplayTasks(db))
	router.GET("/done/:id", useCases.MarkTaskAsCompletedHandler(db))
	router.POST("/", useCases.CreateTaskHandler(db))
	router.POST("/update-task", useCases.UpdateTaskDescriptionHandler(db))
	router.POST("/update-task-title", useCases.UpdateTitleTaskHandler(db))
	router.GET("/project/:project", DisplayTasks(db))
	router.POST("/update-task-due-date", useCases.UpdateDueDateHandler(db))
	router.POST("/update-task-status", useCases.ToggleTaskStatusHandler(db))
	router.Run(":7263")

}

// DisplayTasks returns a gin.HandlerFunc that displays the main page of the
// application. The main page displays 4 lists of tasks: tasks with a deadline,
// tasks without a deadline, tasks done in the last 24 hours and tasks that are
// late.
func DisplayTasks(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// db := persistence.CreateDB()

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
