package views

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/JneiraS/AMS/src/application/useCases"
	"github.com/JneiraS/AMS/src/domain/services"
	"github.com/JneiraS/AMS/src/infrastructure/persistence"
	"github.com/JneiraS/AMS/src/infrastructure/web/components"
	"github.com/JneiraS/AMS/src/infrastructure/web/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const (
	ErrorTemplate = "error.tmpl"
	IndexTemplate = "index.tmpl"
)

// const (
// 	Pending           = "Pending"
// 	InProgress        = "In progress"
// 	Done              = "Done"
// 	invalidRequestMsg = "Invalid request"
// 	invalidTaskIDMsg  = "Invalid task ID"
// 	taskNotFoundMsg   = "Task not found"
// 	FIND_BY_PROJECT   = "project = ?"
// )

// SetupRouter creates a web server that serves the web interface for the application.
//
// The web interface shows four lists of tasks: tasks with a due date, tasks without
// a due date, tasks that are done, and tasks that are late.
//
// The functions to create a new task, mark a task as done and update a task are
// routed to the corresponding functions in the useCases package.
func SetupRouter(db *gorm.DB) {
	router := gin.Default()

	router.SetFuncMap(template.FuncMap{
		"safe": func(s interface{}) template.HTML {
			switch v := s.(type) {
			case string:
				return template.HTML(v)
			case fmt.Stringer:
				return template.HTML(v.String())
			default:
				return template.HTML(fmt.Sprint(v))
			}
		},
	})

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

		var userID int

		priorityTasks, taskWithoutDueDate, tasksDone, lateTasks, project, err := services.GetTasksByCategory(c, db)
		if err != nil {
			c.HTML(http.StatusInternalServerError, ErrorTemplate, gin.H{"error": err.Error()})
			return
		}
		username, _ := c.Cookie("username")
		userID = services.GetUserIDFromContext(c)

		// Affichage de la vue avec toutes les données
		c.HTML(http.StatusOK, IndexTemplate, gin.H{
			"projects":        persistence.GetAllProjects(db),
			"prioritytasks":   priorityTasks,
			"taskswithoutdue": taskWithoutDueDate,
			"tasksdone":       tasksDone,
			"latetasks":       lateTasks,
			"project":         project,
			"statistics":      services.GetStatistics(priorityTasks, taskWithoutDueDate, lateTasks),
			"navbar":          components.Navbar(userID, username),
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

	var userID int

	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}
		task, _ := persistence.GetTask(db, uint(id))

		username, _ := c.Cookie("username")
		userID = services.GetUserIDFromContext(c)

		c.HTML(http.StatusOK, "details.tmpl", gin.H{
			"title":  "Détails de la tâche",
			"task":   task,
			"navbar": components.Navbar(userID, username),
		})
	}
}
