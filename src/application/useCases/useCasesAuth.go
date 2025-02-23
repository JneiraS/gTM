package useCases

import (
	"net/http"

	"github.com/JneiraS/AMS/src/domain/services"
	"github.com/JneiraS/AMS/src/infrastructure/persistence"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
)

var SecretKey = "super_secret_key"

func LoginUserHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Récupérer les données du formulaire HTML
		username := c.PostForm("username")
		password := c.PostForm("password")

		if username == "" || password == "" {
			c.HTML(http.StatusBadRequest, "login.html", gin.H{"error": "Données invalides"})
			return
		}

		// Rechercher l'utilisateur dans la base de données
		var user persistence.User
		if err := db.Where("username = ?", username).First(&user).Error; err != nil {
			c.HTML(http.StatusUnauthorized, "login.html", gin.H{"error": "Utilisateur non trouvé"})
			return
		}

		// Vérifier le mot de passe
		if err := services.CheckPassword(user.Password, password); err != nil {
			c.HTML(http.StatusUnauthorized, "login.html", gin.H{"error": "Mot de passe incorrect"})
			return
		}

		// Générer le token JWT
		token, err := services.GenerateJWT(user)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "login.html", gin.H{"error": "Erreur lors de la génération du token"})
			return
		}

		// ⚠️ Stocker le token dans un cookie sécurisé
		c.SetCookie("token", token, 3600, "/", "", false, true) // Expire après 1h

		// Rediriger vers une page sécurisée
		c.Redirect(http.StatusFound, "/")
	}
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Récupérer le token depuis le cookie
		tokenString, err := c.Cookie("token")
		if err != nil {
			c.Redirect(http.StatusFound, "/login") // Rediriger vers la page de connexion si pas de token
			c.Abort()
			return
		}

		// Vérifier et parser le token
		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(SecretKey), nil
		})

		if err != nil || !token.Valid {
			c.Redirect(http.StatusFound, "/login") // Rediriger vers la page de connexion si token invalide
			c.Abort()
			return
		}

		// Stocker l'ID utilisateur dans le contexte Gin pour l'utiliser dans les templates
		c.Set("userID", claims["sub"])
		c.Next()
	}
}
