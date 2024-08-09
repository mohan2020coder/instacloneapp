package controller

import (
	"instacloneapp/server/pkg/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UserController handles user-related API requests
type UserController struct {
	DB db.Database
}

// GetUsers handles GET requests to retrieve users
func (uc *UserController) GetUsers(c *gin.Context) {
	users, err := uc.DB.GetUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

// CreateUser handles POST requests to create a new user
func (uc *UserController) CreateUser(c *gin.Context) {
	var input struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := db.User{
		Name: input.Name,
	}

	user, err := uc.DB.CreateUser(user)
	// user, err := uc.DB.CreateUser(input.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, user)
}
