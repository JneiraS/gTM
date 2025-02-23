package views

import (
	"fmt"
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
	FIND_BY_PROJECT   = "project = ?"
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

	router.GET("/", useCases.AuthMiddleware(), DisplayTasks(db))
	router.GET("/done/:id", useCases.AuthMiddleware(), useCases.MarkTaskAsCompletedHandler(db))
	router.POST("/", useCases.AuthMiddleware(), useCases.CreateTaskHandler(db))
	router.POST("/update-task", useCases.AuthMiddleware(), useCases.UpdateTaskDescriptionHandler(db))
	router.POST("/update-task-title", useCases.AuthMiddleware(), useCases.UpdateTitleTaskHandler(db))
	router.GET("/project/:project", useCases.AuthMiddleware(), DisplayTasks(db))
	router.POST("/update-task-due-date", useCases.AuthMiddleware(), useCases.UpdateDueDateHandler(db))
	router.POST("/update-task-status", useCases.AuthMiddleware(), useCases.ToggleTaskStatusHandler(db))

	router.POST("/register", useCases.RegisterUserHandler(db))
	router.POST("/login", useCases.LoginUserHandler(db))
	router.GET("/signup", DisplaySignupPage(db))
	router.GET("/login", DisplayLoginPage(db))

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
			baseQuery = baseQuery.Where(FIND_BY_PROJECT, project)
			baseNoDueDateQuery = baseNoDueDateQuery.Where(FIND_BY_PROJECT, project)
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

		if project != "" {
			baseQuery.Where("due_date > ? AND due_date < ? AND status != ?",
				time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Now(),
				Done).
				Find(&lateTasks)
		} else {
			db.Where("due_date > ? AND due_date < ? AND status != ?",
				time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Now(),
				Done).
				Find(&lateTasks)
		}

		// Affichage de la vue avec toutes les données
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title":           "Liste des taches",
			"projects":        persistence.GetAllProjects(db),
			"prioritytasks":   priorityTasks,
			"taskswithoutdue": taskWithoutDueDate,
			"tasksdone":       tasksDone,
			"latetasks":       lateTasks,
			"project":         project,
			"statistics":      statistics(priorityTasks, taskWithoutDueDate, lateTasks),
		})
	}
}

func DisplaySignupPage(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "signup.tmpl", gin.H{
			"title": "Inscription",
		})
	}
}

func DisplayLoginPage(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.tmpl", gin.H{
			"title": "Connexion",
		})
	}
}

func statistics(priorityTasks, taskWithoutDueDate, lateTasks []persistence.Task) []string {
	taskCounts := map[string]int{
		"Priority":         len(priorityTasks),
		"Without Due Date": len(taskWithoutDueDate),
		"Late":             len(lateTasks),
	}

	var totalTasks int
	for _, count := range taskCounts {
		totalTasks += count
	}

	statistics := []string{
		fmt.Sprintf("Total number of tasks: %d", totalTasks),
	}

	for name, count := range taskCounts {
		statistics = append(statistics, fmt.Sprintf("%s: %d", name, count))
	}

	return statistics
}
