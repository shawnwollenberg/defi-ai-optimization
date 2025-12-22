package handlers

import (
	"net/http"
	"strconv"

	"github.com/defioptimization/shared/database"
	"github.com/defioptimization/shared/models"
	"github.com/gin-gonic/gin"
)

// GetTransactions returns all transactions for the current user
func GetTransactions(c *gin.Context) {
	userID, _ := c.Get("user_id")
	limit := c.DefaultQuery("limit", "50")
	offset := c.DefaultQuery("offset", "0")

	limitInt, _ := strconv.Atoi(limit)
	offsetInt, _ := strconv.Atoi(offset)

	var transactions []models.Transaction
	if err := database.DB.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limitInt).
		Offset(offsetInt).
		Find(&transactions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch transactions"})
		return
	}

	c.JSON(http.StatusOK, transactions)
}

