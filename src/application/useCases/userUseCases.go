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
		var user models.User

		// Récupérer les données JSON envoyées par le client
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Hasher le mot de passe
		hashedPassword, err := services.HashPassword(user.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors du hashage du mot de passe"})
			return
		}

		// Créer l'utilisateur
		persistenceUser := persistence.User{
			User: models.User{
				Username: c.PostForm("username"),
				Password: hashedPassword,
			},
		}

		// Enregistrer dans la base de données
		result := db.Create(&persistenceUser)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}

		c.Redirect(http.StatusFound, "/")
		c.JSON(http.StatusCreated, gin.H{"message": "Utilisateur créé avec succès"})
	}
}
