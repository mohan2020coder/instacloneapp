package middleware

import (
	
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	
)

// Middleware function to check if the user is authenticated
func IsAuthenticated() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Load environment variables from .env file (if using godotenv)
		// err := godotenv.Load()
		// if err != nil {
		// 	log.Fatalf("Error loading .env file")
		// }

		tokenString, err := c.Cookie("token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "User not authenticated",
				"success": false,
			})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate the token's signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, &jwt.ValidationError{Errors: jwt.ValidationErrorSignatureInvalid}
			}
			// Return the key for validating the token
			return []byte(os.Getenv("SECRET_KEY")), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Invalid token",
				"success": false,
			})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Invalid token claims",
				"success": false,
			})
			c.Abort()
			return
		}

		// Set user ID in the request context
		c.Set("userId", claims["userId"])

		// Proceed to the next handler
		c.Next()
	}
}
