package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/defioptimization/defi-service/protocols"
	"github.com/gin-gonic/gin"
)

// Server handles HTTP requests for the DeFi service
type Server struct {
	protocolManager *protocols.Manager
	router          *gin.Engine
}

// NewServer creates a new server instance
func NewServer(pm *protocols.Manager) *Server {
	r := gin.Default()
	s := &Server{
		protocolManager: pm,
		router:          r,
	}
	s.setupRoutes()
	return s
}

// setupRoutes configures API routes
func (s *Server) setupRoutes() {
	api := s.router.Group("/api/v1")
	{
		api.GET("/health", s.healthCheck)
		api.GET("/protocols", s.getProtocols)
		api.GET("/protocols/:name/apy", s.getAPY)
		api.GET("/protocols/:name/positions", s.getUserPositions)
		api.GET("/protocols/:name/health-factor", s.getHealthFactor)
		api.GET("/protocols/:name/price", s.getAssetPrice)
	}
}

// Start starts the HTTP server
func (s *Server) Start(addr string) error {
	return s.router.Run(addr)
}

// healthCheck returns service health status
func (s *Server) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "defi-service",
	})
}

// getProtocols returns all available protocols
func (s *Server) getProtocols(c *gin.Context) {
	protocolList := s.protocolManager.GetAllProtocols()
	protocols := make([]gin.H, len(protocolList))
	
	for i, p := range protocolList {
		protocols[i] = gin.H{
			"name": p.GetName(),
		}
	}
	
	c.JSON(http.StatusOK, protocols)
}

// getAPY returns the APY for a specific protocol and asset
func (s *Server) getAPY(c *gin.Context) {
	protocolName := c.Param("name")
	asset := c.Query("asset")
	chain := c.DefaultQuery("chain", "ethereum")
	
	protocol, ok := s.protocolManager.GetProtocol(protocolName)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Protocol not found"})
		return
	}
	
	if asset == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "asset parameter is required"})
		return
	}
	
	apy, err := protocol.GetAPY(c.Request.Context(), asset, chain)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"protocol": protocolName,
		"asset":    asset,
		"chain":    chain,
		"apy":      apy,
	})
}

// getUserPositions returns user positions for a protocol
func (s *Server) getUserPositions(c *gin.Context) {
	protocolName := c.Param("name")
	userAddress := c.Query("user_address")
	chain := c.DefaultQuery("chain", "ethereum")
	
	if userAddress == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_address parameter is required"})
		return
	}
	
	protocol, ok := s.protocolManager.GetProtocol(protocolName)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Protocol not found"})
		return
	}
	
	positions, err := protocol.GetUserPositions(c.Request.Context(), userAddress, chain)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, positions)
}

// getHealthFactor returns the health factor for a user
func (s *Server) getHealthFactor(c *gin.Context) {
	protocolName := c.Param("name")
	userAddress := c.Query("user_address")
	chain := c.DefaultQuery("chain", "ethereum")
	
	if userAddress == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_address parameter is required"})
		return
	}
	
	protocol, ok := s.protocolManager.GetProtocol(protocolName)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Protocol not found"})
		return
	}
	
	healthFactor, err := protocol.GetHealthFactor(c.Request.Context(), userAddress, chain)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"protocol":     protocolName,
		"user_address": userAddress,
		"chain":        chain,
		"health_factor": healthFactor,
	})
}

// getAssetPrice returns the price of an asset
func (s *Server) getAssetPrice(c *gin.Context) {
	protocolName := c.Param("name")
	asset := c.Query("asset")
	chain := c.DefaultQuery("chain", "ethereum")
	
	if asset == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "asset parameter is required"})
		return
	}
	
	protocol, ok := s.protocolManager.GetProtocol(protocolName)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Protocol not found"})
		return
	}
	
	price, err := protocol.GetAssetPrice(c.Request.Context(), asset, chain)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"protocol": protocolName,
		"asset":    asset,
		"chain":    chain,
		"price":    price,
	})
}

// Helper function (unused but available)
func _() {
	_ = json.Marshal
	_ = strconv.Atoi
}

