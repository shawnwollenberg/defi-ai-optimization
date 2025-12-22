package server

import (
	"net/http"

	"github.com/defioptimization/wallet/connector"
	"github.com/gin-gonic/gin"
)

// Server handles HTTP requests for the wallet service
type Server struct {
	router     *gin.Engine
	connector  *connector.WalletConnector
}

// NewServer creates a new server instance
func NewServer() *Server {
	r := gin.Default()
	wc := connector.NewWalletConnector()
	
	s := &Server{
		router:    r,
		connector: wc,
	}
	s.setupRoutes()
	return s
}

// setupRoutes configures API routes
func (s *Server) setupRoutes() {
	api := s.router.Group("/api/v1")
	{
		api.GET("/health", s.healthCheck)
		api.POST("/wallet/connect", s.connectWallet)
		api.POST("/wallet/disconnect", s.disconnectWallet)
		api.POST("/wallet/sign", s.signMessage)
		api.POST("/wallet/send", s.sendTransaction)
		api.POST("/wallet/build", s.buildTransaction)
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
		"service": "wallet-service",
	})
}

// connectWallet handles wallet connection requests
func (s *Server) connectWallet(c *gin.Context) {
	var req struct {
		WalletAddress string `json:"wallet_address" binding:"required"`
		Chain         string `json:"chain" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// TODO: Implement actual wallet connection logic
	// This would typically involve WalletConnect session creation
	
	c.JSON(http.StatusOK, gin.H{
		"connected": true,
		"wallet_address": req.WalletAddress,
		"chain": req.Chain,
	})
}

// disconnectWallet handles wallet disconnection
func (s *Server) disconnectWallet(c *gin.Context) {
	var req struct {
		WalletAddress string `json:"wallet_address" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"disconnected": true})
}

// signMessage handles message signing requests
func (s *Server) signMessage(c *gin.Context) {
	var req struct {
		WalletAddress string `json:"wallet_address" binding:"required"`
		Message       string `json:"message" binding:"required"`
		Chain         string `json:"chain" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// TODO: Implement actual message signing
	// This would use WalletConnect or MetaMask SDK
	
	c.JSON(http.StatusOK, gin.H{
		"signature": "0x...", // Placeholder
		"message": req.Message,
	})
}

// sendTransaction handles transaction sending
func (s *Server) sendTransaction(c *gin.Context) {
	var req struct {
		WalletAddress string                 `json:"wallet_address" binding:"required"`
		Chain         string                 `json:"chain" binding:"required"`
		To            string                 `json:"to" binding:"required"`
		Value         string                 `json:"value"`
		Data          string                 `json:"data"`
		GasLimit      string                 `json:"gas_limit"`
		GasPrice      string                 `json:"gas_price"`
		Params        map[string]interface{} `json:"params"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// TODO: Implement actual transaction building and sending
	// This would build the transaction and return it for user approval
	
	c.JSON(http.StatusOK, gin.H{
		"tx_hash": "0x...", // Placeholder
		"status": "pending",
	})
}

// buildTransaction builds a transaction without sending
func (s *Server) buildTransaction(c *gin.Context) {
	var req struct {
		Chain    string                 `json:"chain" binding:"required"`
		To       string                 `json:"to" binding:"required"`
		Value    string                 `json:"value"`
		Data     string                 `json:"data"`
		Params   map[string]interface{} `json:"params"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Build transaction
	tx, err := s.connector.BuildTransaction(req.Chain, req.To, req.Value, req.Data, req.Params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, tx)
}

