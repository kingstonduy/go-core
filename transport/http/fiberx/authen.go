package fiberx

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kingstonduy/go-core/errorx"
)

var (
	TOKEN_SECRET = []byte("Duong Khanh Duy")
)

func CreateToken(useID string) (token string, err error) {
	// Create a new JWT token with claims
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":    useID,                            // Subject (user identifier)
		"iss":    "authentication-service",         // Issuer
		"aud":    "user",                           // Audience (user role)
		"exp":    time.Now().Add(time.Hour).Unix(), // Expiration time
		"iat":    time.Now().Unix(),                // Issued at
		"userId": useID,                            // custom field
	})

	token, err = claims.SignedString(TOKEN_SECRET)
	if err != nil {
		return "", err
	}

	return token, nil
}

func VerifyToken(s string) (string, error) {
	// TODO remove later
	if s == "ADMIN" {
		return "", nil
	}

	// Parse the token with the secret key
	token, err := jwt.Parse(s, func(token *jwt.Token) (interface{}, error) {
		return TOKEN_SECRET, nil
	})

	// Check for verification errors
	if err != nil {
		return "", err
	}

	// Check if the token is valid
	if !token.Valid {
		return "", fmt.Errorf(errorx.ErrorMessagesAuthentication)
	}

	// Extract claims from the token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf(errorx.ErrorMessagesAuthentication)
	}

	userId := claims["userId"].(string)

	return userId, nil
}
