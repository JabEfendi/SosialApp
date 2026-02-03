package helpers

import (
	"time"
	"os"

  "github.com/golang-jwt/jwt/v5"
)

func GenerateTemp2FAToken(corporateID uint) (string, error) {
	claims := jwt.MapClaims{
		"corporate_id": corporateID,
		"type":         "2fa",
		"exp":          time.Now().Add(5 * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}