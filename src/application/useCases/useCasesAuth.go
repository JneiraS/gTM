package useCases

import (
	"net/http"

	as "github.com/JneiraS/AMS/src/application/services"
	"github.com/JneiraS/AMS/src/domain/services"
	"github.com/JneiraS/AMS/src/infrastructure/persistence"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const (
	LOGIN_PAGE = "login.html"
)

func LoginUserHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Récupérer les données du formulaire HTML
		username := c.PostForm("username")
		password := c.PostForm("password")

		if username == "" || password == "" {
			c.HTML(http.StatusBadRequest, LOGIN_PAGE, gin.H{"error": "Données invalides"})
			return
		}

		// Rechercher l'utilisateur dans la base de données
		var user persistence.User
		if err := db.Where("username = ?", username).First(&user).Error; err != nil {
			c.HTML(http.StatusUnauthorized, LOGIN_PAGE, gin.H{"error": "Utilisateur non trouvé"})
			return
		}

		// Vérifier le mot de passe
		if err := services.CheckPassword(user.Password, password); err != nil {
			c.HTML(http.StatusUnauthorized, LOGIN_PAGE, gin.H{"error": "Mot de passe incorrect"})
			return
		}

		// Générer le token JWT
		token, err := as.GenerateJWT(user)
		if err != nil {
			c.HTML(http.StatusInternalServerError, LOGIN_PAGE, gin.H{"error": "Erreur lors de la génération du token"})
			return
		}

		// ⚠️ Stocker le token dans un cookie sécurisé
		c.SetCookie("token", token, 3600, "/", "", false, true)
		c.SetCookie("username", user.Username, 3600, "/", "", false, true) // Expire après 1h

		// Rediriger vers une page sécurisée
		c.Redirect(http.StatusFound, "/")
	}
}

// LogoutHandler supprime le cookie de session et redirige vers la page d'accueil.
func LogoutHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.SetCookie("token", "", -1, "/", "", false, true)
		c.Redirect(http.StatusFound, "/login")
	}
}
