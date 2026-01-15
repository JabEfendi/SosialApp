package helpers

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Ambil secret dari environment variable
var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

// ===================== ADMIN =====================

// Claims untuk admin
type AdminClaims struct {
	AdminID uint   `json:"admin_id"`
	Role    string `json:"role"`
	jwt.RegisteredClaims
}

// Generate token untuk admin
func GenerateAdminToken(adminID uint, role string) (string, error) {
	claims := AdminClaims{
		AdminID: adminID,
		Role:    role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // token 24 jam
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// Validasi token admin
func ValidateAdminToken(tokenString string) (*AdminClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AdminClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	claims, ok := token.Claims.(*AdminClaims)
	if !ok {
		return nil, errors.New("invalid admin claims")
	}

	return claims, nil
}

// ===================== USER =====================

// Claims untuk user
type UserClaims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

// Generate token untuk user
func GenerateUserToken(userID uint) (string, error) {
	claims := UserClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// Validasi token user
func ValidateUserToken(tokenString string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}