package useCases

import (
	"net/http"

	"github.com/JneiraS/AMS/src/domain/models"
	"github.com/JneiraS/AMS/src/domain/services"
	"github.com/JneiraS/AMS/src/infrastructure/persistence"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

// Route pour l'inscription d'un utilisateur
func RegisterUserHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")

		if username == "" || password == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Username and password are required"})
			return
		}

		hashedPassword, err := services.HashPassword(password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}

		user := persistence.User{
			User: models.User{
				Username: username,
				Password: hashedPassword,
			},
		}

		if err := db.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
		c.Redirect(http.StatusFound, "/")
	}
}
