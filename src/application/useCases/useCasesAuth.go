package useCases

import (
	"net/http"

	"github.com/JneiraS/AMS/src/domain/models"
	"github.com/JneiraS/AMS/src/domain/services"
	"github.com/JneiraS/AMS/src/infrastructure/persistence"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func LoginUserHandler(db *gorm.DB) gin.HandlerFunc {

	return func(c *gin.Context) {
		var input models.LoginInput

		// Vérifier que l'entrée est correcte
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Données invalides"})
			return
		}

		// Rechercher l'utilisateur par email
		var user persistence.User
		if err := db.Where("username = ?", c.PostForm("username")).First(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Utilisateur non trouvé"})
			return
		}

		// Vérifier le mot de passe
		if err := services.CheckPassword(user.Password, c.PostForm("password")); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Mot de passe incorrect"})
			return
		}

		// Générer un token JWT
		token, err := services.GenerateJWT(user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la génération du token"})
			return
		}

		// Retourner le token
		c.JSON(http.StatusOK, gin.H{"token": token})
	}
}
