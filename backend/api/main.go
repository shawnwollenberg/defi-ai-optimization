package main

import (
	"log"
	"os"

	"github.com/defioptimization/api/handlers"
	"github.com/defioptimization/api/middleware"
	"github.com/defioptimization/api/websocket"
	"github.com/defioptimization/shared/database"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database
	if err := database.InitDatabase(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.CloseDatabase()

	// Initialize Stripe
	handlers.InitStripe()

	// Initialize WebSocket hub
	hub := websocket.NewHub()
	go hub.Run()

	// Initialize router
	r := gin.Default()

	// CORS configuration
	config := cors.DefaultConfig()
	allowedOrigins := []string{"http://localhost:3000"}
	if frontendURL := os.Getenv("FRONTEND_URL"); frontendURL != "" {
		allowedOrigins = append(allowedOrigins, frontendURL)
	}
	config.AllowOrigins = allowedOrigins
	config.AllowCredentials = true
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	r.Use(cors.New(config))

	// Health check
	r.GET("/health", handlers.HealthCheck)

	// Public routes
	public := r.Group("/api/v1")
	{
		public.POST("/auth/wallet", handlers.WalletAuth)
		public.GET("/protocols", handlers.GetProtocols)
	}

	// Protected routes
	protected := r.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware())
	{
		// User routes
		protected.GET("/user/profile", handlers.GetUserProfile)
		protected.PUT("/user/profile", handlers.UpdateUserProfile)

		// Portfolio routes
		protected.GET("/portfolios", handlers.GetPortfolios)
		protected.POST("/portfolios", handlers.CreatePortfolio)
		protected.GET("/portfolios/:id", handlers.GetPortfolio)
		protected.PUT("/portfolios/:id", handlers.UpdatePortfolio)
		protected.DELETE("/portfolios/:id", handlers.DeletePortfolio)

		// Automation routes
		protected.GET("/automation/rules", handlers.GetAutomationRules)
		protected.POST("/automation/rules", handlers.CreateAutomationRule)
		protected.PUT("/automation/rules/:id", handlers.UpdateAutomationRule)
		protected.DELETE("/automation/rules/:id", handlers.DeleteAutomationRule)

		// Transaction routes
		protected.GET("/transactions", handlers.GetTransactions)

		// Subscription routes
		protected.GET("/subscription", handlers.GetSubscription)
		protected.POST("/subscription/upgrade", handlers.UpgradeSubscription)
		protected.POST("/subscription/checkout", func(c *gin.Context) {
			handlers.CreateCheckoutSession(c)
		})

		// WebSocket route
		protected.GET("/ws", websocket.HandleWebSocket(hub))
	}

	// Stripe webhook (no auth required, uses signature verification)
	r.POST("/api/v1/webhooks/stripe", handlers.HandleStripeWebhook)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("API Gateway starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

