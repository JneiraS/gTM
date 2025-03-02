package views

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/JneiraS/AMS/src/application/useCases"
	"github.com/JneiraS/AMS/src/infrastructure/persistence"
	"github.com/JneiraS/AMS/src/infrastructure/web/middleware"
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

	router.GET("/", middleware.AuthMiddleware(), DisplayTasks(db))
	router.GET("/done/:id", middleware.AuthMiddleware(), useCases.MarkTaskAsCompletedHandler(db))
	router.POST("/", middleware.AuthMiddleware(), useCases.CreateTaskHandler(db))
	router.POST("/update-task", middleware.AuthMiddleware(), useCases.UpdateTaskDescriptionHandler(db))
	router.POST("/update-task-title", middleware.AuthMiddleware(), useCases.UpdateTitleTaskHandler(db))
	router.GET("/project/:project", middleware.AuthMiddleware(), DisplayTasks(db))
	router.POST("/update-task-due-date", middleware.AuthMiddleware(), useCases.UpdateDueDateHandler(db))
	router.POST("/update-task-status", middleware.AuthMiddleware(), useCases.ToggleTaskStatusHandler(db))

	router.POST("/register", useCases.RegisterUserHandler(db))
	router.POST("/login", useCases.LoginUserHandler(db))
	router.GET("/signup", DisplaySignupPage(db))
	router.GET("/login", DisplayLoginPage(db))

	router.GET("/logout", useCases.LogoutHandler())

	router.GET("/task/:id", middleware.AuthMiddleware(), DisplayDetailsTask(db))

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

			baseQuery.Where("due_date > ? AND due_date < ? AND status != ?",
				time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Now(),
				Done).
				Find(&lateTasks)
		} else {
			db.Where("status = ? AND updated_at > ? AND updated_at < ?",
				Done,
				time.Now().AddDate(0, 0, -1),
				time.Now()).
				Order("updated_at ASC").
				Find(&tasksDone)

			db.Where("due_date > ? AND due_date < ? AND status != ?",
				time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC),
				time.Now(),
				Done).
				Find(&lateTasks)
		}
		baseQuery.Order("priority ASC, due_date ASC").Find(&priorityTasks)
		baseNoDueDateQuery.Order("priority ASC").Find(&taskWithoutDueDate)

		username, _ := c.Cookie("username")

		// Affichage de la vue avec toutes les données
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title":           "Liste des taches",
			"projects":        persistence.GetAllProjects(db),
			"prioritytasks":   priorityTasks,
			"taskswithoutdue": taskWithoutDueDate,
			"tasksdone":       tasksDone,
			"latetasks":       lateTasks,
			"project":         project,
			"user":            username,
			"statistics":      getStatistics(priorityTasks, taskWithoutDueDate, lateTasks),
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

func DisplayDetailsTask(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		task, _ := persistence.GetTask(db, uint(id))
		username, _ := c.Cookie("username")

		c.HTML(http.StatusOK, "details.tmpl", gin.H{
			"title": "Détails de la tâche",
			"task":  task,
			"user":  username,
		})
	}
}

// getStatistics returns an array of strings which represent the statistics
// of the tasks. The first element is the total number of tasks to do,
// and the next elements are the number of tasks with a specific
// characteristic.
func getStatistics(priorityTasks, taskWithoutDueDate, lateTasks []persistence.Task) []string {
	taskCounts := map[string]int{
		"Tasks whith priority": len(priorityTasks),
		"Without Due Date":     len(taskWithoutDueDate),
		"Overdue":              len(lateTasks),
	}

	var totalTasks int
	for _, count := range taskCounts {
		totalTasks += count
	}

	statistics := []string{
		fmt.Sprintf("Total number of tasks to do: %d", totalTasks),
	}

	for name, count := range taskCounts {
		statistics = append(statistics, fmt.Sprintf("%s: %d", name, count))
	}

	return statistics
}
