package helpers

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func getJWTSecret() []byte {
	return []byte(os.Getenv("JWT_SECRET"))
}

/* ===================== ADMIN ===================== */

type AdminClaims struct {
	AdminID uint   `json:"admin_id"`
	Role    string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateAdminToken(adminID uint, role string) (string, error) {
	claims := AdminClaims{
		AdminID: adminID,
		Role:    role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getJWTSecret())
}

func ValidateAdminToken(tokenString string) (*AdminClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&AdminClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return getJWTSecret(), nil
		},
	)

	if err != nil || !token.Valid {
		return nil, err
	}

	claims, ok := token.Claims.(*AdminClaims)
	if !ok {
		return nil, errors.New("invalid admin claims")
	}

	return claims, nil
}

/* ===================== USER ===================== */

type UserClaims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateUserToken(userID uint) (string, error) {
	claims := UserClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getJWTSecret())
}

func ValidateUserToken(tokenString string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return getJWTSecret(), nil
		},
	)

	if err != nil || !token.Valid {
		return nil, err
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}

/* ===================== CORPORATE ===================== */

type CorporateClaims struct {
	CorporateID uint `json:"corporate_id"`
	jwt.RegisteredClaims
}

func GenerateCorporateToken(corporateID uint) (string, error) {
	claims := CorporateClaims{
		CorporateID: corporateID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getJWTSecret())
}

func ValidateCorporateToken(tokenString string) (*CorporateClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&CorporateClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return getJWTSecret(), nil
		},
	)

	if err != nil || !token.Valid {
		return nil, err
	}

	claims, ok := token.Claims.(*CorporateClaims)
	if !ok {
		return nil, errors.New("invalid corporate claims")
	}

	return claims, nil
}