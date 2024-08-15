package middleware

import (
	"fmt"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// Middleware function to check if the user is authenticated
func IsAuthenticated() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := extractClaims(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "User not authenticated or invalid token",
				"success": false,
			})
			c.Abort()
			return
		}
		//fmt.Println("User ID:", claims)

		// Set user ID in the request context
		c.Set("userID", claims["userID"])

		// Proceed to the next handler
		c.Next()
	}
}

// Extract claims from the token
func extractClaims(c *gin.Context) (map[string]interface{}, error) {
	// Load environment variables from .env file (if using godotenv)
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error loading .env file")
	}

	// Get the token from the cookie
	tokenString, err := c.Cookie("token")
	if err != nil {
		return nil, fmt.Errorf("error retrieving token from cookie")
	}

	//fmt.Println("Token from cookie:", tokenString)

	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the token's signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		// Return the key for validating the token
		return []byte(os.Getenv("SECRET_KEY")), nil
	})

	if err != nil {
		return nil, err
	}

	// Extract claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token claims")
}
