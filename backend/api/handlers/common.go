package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthCheck returns the health status of the API
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "api-gateway",
	})
}

// GetProtocols returns available DeFi protocols
func GetProtocols(c *gin.Context) {
	protocols := []gin.H{
		{
			"name":        "Aave",
			"description": "Decentralized lending and borrowing",
			"chains":      []string{"ethereum", "base"},
			"supported":   true,
		},
		{
			"name":        "Compound",
			"description": "Algorithmic money markets",
			"chains":      []string{"ethereum", "base"},
			"supported":   true,
		},
		{
			"name":        "EigenLayer",
			"description": "Restaking protocol",
			"chains":      []string{"ethereum"},
			"supported":   true,
		},
	}

	c.JSON(http.StatusOK, protocols)
}

