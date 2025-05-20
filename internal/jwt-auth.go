package main

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

// JWT_SECRET should be set in the environment variables
var secretKey = []byte(os.Getenv("JWT_SECRET"))

type MyCustomClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// CreateJWT creates a JWT token string
func CreateJWT(username string) (string, error) {
	claims := MyCustomClaims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)), // token expires in 15 minutes
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "bcr-auth",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

// ValidateJWT parses and validates the JWT token string
func ValidateJWT(tokenStr string) (*MyCustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &MyCustomClaims{}, func(token *jwt.Token) (any, error) {
		// Make sure that the token method conforms to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
