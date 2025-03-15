package services

import (
	"time"

	"github.com/JneiraS/AMS/src/infrastructure/persistence"
	"github.com/golang-jwt/jwt"
)

var secretKey = "super_secret_key"

func GenerateJWT(user persistence.User) (string, error) {
	// Créer le token avec une expiration de 24h
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := jwt.MapClaims{
		"sub": user.ID,
		"exp": expirationTime.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Signer le token
	return token.SignedString([]byte(secretKey))
}
