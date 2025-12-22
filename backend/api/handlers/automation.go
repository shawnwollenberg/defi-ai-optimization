package handlers

import (
	"net/http"
	"strconv"

	"github.com/defioptimization/shared/database"
	"github.com/defioptimization/shared/models"
	"github.com/gin-gonic/gin"
)

// GetAutomationRules returns all automation rules for the current user
func GetAutomationRules(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var rules []models.AutomationRule
	if err := database.DB.Where("user_id = ?", userID).Find(&rules).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch automation rules"})
		return
	}

	c.JSON(http.StatusOK, rules)
}

// CreateAutomationRule creates a new automation rule
func CreateAutomationRule(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var rule models.AutomationRule
	if err := c.ShouldBindJSON(&rule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rule.UserID = userID.(uint)

	if err := database.DB.Create(&rule).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create automation rule"})
		return
	}

	c.JSON(http.StatusCreated, rule)
}

// UpdateAutomationRule updates an automation rule
func UpdateAutomationRule(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var rule models.AutomationRule
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&rule).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Automation rule not found"})
		return
	}

	if err := c.ShouldBindJSON(&rule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Save(&rule).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update automation rule"})
		return
	}

	c.JSON(http.StatusOK, rule)
}

// DeleteAutomationRule deletes an automation rule
func DeleteAutomationRule(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&models.AutomationRule{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete automation rule"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Automation rule deleted"})
}

