package handlers

import (
	"net/http"

	"github.com/defioptimization/shared/database"
	"github.com/defioptimization/shared/models"
	"github.com/gin-gonic/gin"
)

// GetUserProfile returns the current user's profile
func GetUserProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var user models.User
	if err := database.DB.Preload("Portfolios").Preload("AutomationRules").First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateUserProfile updates the current user's profile
func UpdateUserProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var updateData struct {
		Email       *string                `json:"email"`
		Preferences map[string]interface{} `json:"preferences"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if updateData.Email != nil {
		user.Email = *updateData.Email
	}
	if updateData.Preferences != nil {
		user.Preferences = updateData.Preferences
	}

	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, user)
}

