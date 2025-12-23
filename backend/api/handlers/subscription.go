package handlers

import (
	"net/http"

	"github.com/defioptimization/shared/database"
	"github.com/defioptimization/shared/models"
	"github.com/gin-gonic/gin"
)

// GetSubscription returns the current user's subscription
func GetSubscription(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var subscription models.Subscription
	if err := database.DB.Where("user_id = ?", userID).First(&subscription).Error; err != nil {
		// Return default free tier if no subscription exists
		c.JSON(http.StatusOK, gin.H{
			"tier":   "free",
			"status": "active",
		})
		return
	}

	c.JSON(http.StatusOK, subscription)
}

// UpgradeSubscriptionRequest represents a subscription upgrade request
type UpgradeSubscriptionRequest struct {
	Tier string `json:"tier" binding:"required"` // basic, premium
}

// UpgradeSubscription upgrades the user's subscription
func UpgradeSubscription(c *gin.Context) {
	var req UpgradeSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create Stripe checkout session
	CreateCheckoutSession(c)
}

