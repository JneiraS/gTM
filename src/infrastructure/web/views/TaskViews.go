package views

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	as "github.com/JneiraS/AMS/src/application/services"
	uc "github.com/JneiraS/AMS/src/application/useCases"
	p "github.com/JneiraS/AMS/src/infrastructure/persistence"
	comp "github.com/JneiraS/AMS/src/infrastructure/web/components"
	mw "github.com/JneiraS/AMS/src/infrastructure/web/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const (
	ErrorTemplate = "error.tmpl"
	IndexTemplate = "index.tmpl"
)

// SetupRouter creates a web server that serves the web interface for the application.
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

	// Requests GET
	router.GET("/", mw.AuthMiddleware(), DisplayTasks(db))
	router.GET("/done/:id", mw.AuthMiddleware(), uc.MarkTaskAsCompletedHandler(db))
	router.GET("/project/:project", mw.AuthMiddleware(), DisplayTasks(db))
	router.GET("/signup", DisplaySignupPage(db))
	router.GET("/login", DisplayLoginPage(db))
	router.GET("/logout", uc.LogoutHandler())
	router.GET("/task/:id", mw.AuthMiddleware(), DisplayDetailsTask(db))

	// Requests POST
	router.POST("/", mw.AuthMiddleware(), uc.CreateTaskHandler(db))
	router.POST("/update-task", mw.AuthMiddleware(), uc.UpdateTaskDescriptionHandler(db))
	router.POST("/update-task-title", mw.AuthMiddleware(), uc.UpdateTitleTaskHandler(db))
	router.POST("/update-task-due-date", mw.AuthMiddleware(), uc.UpdateDueDateHandler(db))
	router.POST("/update-task-status", mw.AuthMiddleware(), uc.ToggleTaskStatusHandler(db))
	router.POST("/register", uc.RegisterUserHandler(db))
	router.POST("/login", uc.LoginUserHandler(db))
	router.POST("/comment/:id", mw.AuthMiddleware(), uc.CreateCommentHandler(db))
	router.POST("/subtask/:id", mw.AuthMiddleware(), uc.CreateSubtaskHandler(db))
	router.POST("/subtask/:id/status-change/:id_sub", mw.AuthMiddleware(), uc.UpdateSubtaskStatusHandler(db))

	router.Run(":7263")

}

// DisplayTasks returns a gin.HandlerFunc that displays the main page of the
// application. The main page displays 4 lists of tasks: tasks with a deadline,
// tasks without a deadline, tasks done in the last 24 hours and tasks that are
// late.
func DisplayTasks(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		var userID int

		priorityTasks, taskWithoutDueDate, tasksDone, lateTasks, project, err := as.GetTasksByCategory(c, db)
		if err != nil {
			c.HTML(http.StatusInternalServerError, ErrorTemplate, gin.H{"error": err.Error()})
			return
		}
		username, _ := c.Cookie("username")
		userID = as.GetUserIDFromContext(c)

		// Affichage de la vue avec toutes les données
		c.HTML(http.StatusOK, IndexTemplate, gin.H{
			"projects":        p.GetAllProjects(db),
			"prioritytasks":   priorityTasks,
			"taskswithoutdue": taskWithoutDueDate,
			"tasksdone":       tasksDone,
			"latetasks":       lateTasks,
			"project":         project,
			"statistics":      as.GetStatistics(priorityTasks, taskWithoutDueDate, lateTasks),
			"navbar":          comp.Navbar(userID, username),
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
		task, _ := p.GetTask(db, uint(id))

		username, _ := c.Cookie("username")
		userID = as.GetUserIDFromContext(c)

		listOfComment, err := p.GetAllCommentsOfTask(db, uint(id))
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}

		listOfSubtask, _ := p.GetAllSubtasksOfTask(db, uint(id))

		c.HTML(http.StatusOK, "details.tmpl", gin.H{
			"title":    "Détails de la tâche",
			"navbar":   comp.Navbar(userID, username),
			"detail":   comp.CardDetails(task),
			"comments": comp.CardComments(task, listOfComment),
			"subtasks": comp.CardSubtasks(task, listOfSubtask),
		})
	}
}
