package handlers

import (
	"net/http"
	"strconv"

	"github.com/defioptimization/shared/database"
	"github.com/defioptimization/shared/models"
	"github.com/gin-gonic/gin"
)

// GetPortfolios returns all portfolios for the current user
func GetPortfolios(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var portfolios []models.Portfolio
	if err := database.DB.Where("user_id = ?", userID).Preload("Positions").Find(&portfolios).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch portfolios"})
		return
	}

	c.JSON(http.StatusOK, portfolios)
}

// CreatePortfolio creates a new portfolio
func CreatePortfolio(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var portfolio models.Portfolio
	if err := c.ShouldBindJSON(&portfolio); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	portfolio.UserID = userID.(uint)

	if err := database.DB.Create(&portfolio).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create portfolio"})
		return
	}

	c.JSON(http.StatusCreated, portfolio)
}

// GetPortfolio returns a specific portfolio
func GetPortfolio(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var portfolio models.Portfolio
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).Preload("Positions").Preload("Snapshots").First(&portfolio).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Portfolio not found"})
		return
	}

	c.JSON(http.StatusOK, portfolio)
}

// UpdatePortfolio updates a portfolio
func UpdatePortfolio(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var portfolio models.Portfolio
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&portfolio).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Portfolio not found"})
		return
	}

	if err := c.ShouldBindJSON(&portfolio); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Save(&portfolio).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update portfolio"})
		return
	}

	c.JSON(http.StatusOK, portfolio)
}

// DeletePortfolio deletes a portfolio
func DeletePortfolio(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Portfolio{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete portfolio"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Portfolio deleted"})
}

