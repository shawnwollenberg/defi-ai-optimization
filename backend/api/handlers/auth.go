package handlers

import (
	"net/http"
	"os"
	"time"

	"github.com/defioptimization/shared/database"
	"github.com/defioptimization/shared/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// WalletAuthRequest represents a wallet authentication request
type WalletAuthRequest struct {
	WalletAddress string `json:"wallet_address" binding:"required"`
	Signature     string `json:"signature" binding:"required"`
	Message       string `json:"message" binding:"required"`
}

// WalletAuth handles wallet-based authentication
func WalletAuth(c *gin.Context) {
	var req WalletAuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Verify signature
	// For now, we'll create or get the user

	var user models.User
	result := database.DB.Where("wallet_address = ?", req.WalletAddress).First(&user)

	if result.Error != nil {
		// Create new user
		user = models.User{
			WalletAddress:   req.WalletAddress,
			SubscriptionTier: "free",
		}
		if err := database.DB.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"wallet_address": user.WalletAddress,
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 days
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
		"user": gin.H{
			"id":             user.ID,
			"wallet_address": user.WalletAddress,
			"subscription_tier": user.SubscriptionTier,
		},
	})
}

