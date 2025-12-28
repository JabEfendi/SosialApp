package helpers

import (
	"time"
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("super-secret-key-123")

func GenerateAdminToken(adminID uint, role string) (string, error) {
	claims := jwt.MapClaims{
		"admin_id": adminID,
		"role":     role,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

type AdminClaims struct {
	AdminID uint   `json:"admin_id"`
	Role    string `json:"role"`
	jwt.RegisteredClaims
}

func ValidateAdminToken(tokenString string) (*AdminClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&AdminClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		},
	)

	if err != nil || !token.Valid {
		return nil, err
	}

	claims, ok := token.Claims.(*AdminClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}

	return claims, nil
}