package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
)

func GenerateToken(userID string) (string, error) {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		return "", err // Return the error to be handled by the caller
	}

	// Get environment variables
	jwtSecret := os.Getenv("SECRET_KEY")
	if jwtSecret == "" {
		return "", fmt.Errorf("SECRET_KEY not set in .env file")
	}

	// Define the token claims
	claims := jwt.MapClaims{
		"userID": userID,
		"exp":    time.Now().Add(time.Hour * 72).Unix(), // Token expiration time (72 hours)
	}

	// Create a new JWT token with the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Convert the secret key to a byte slice
	secretKey := []byte(jwtSecret)

	// Sign the token with the secret key
	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
