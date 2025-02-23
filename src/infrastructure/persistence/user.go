package persistence

import (
	"github.com/JneiraS/AMS/src/domain/models"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	models.User
}
