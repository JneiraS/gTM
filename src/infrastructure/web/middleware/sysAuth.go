package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"net/http"
)

var SecretKey = "super_secret_key"

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
