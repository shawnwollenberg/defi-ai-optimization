package main

import (
	"log"
	"os"

	"github.com/defioptimization/defi-service/protocols"
	"github.com/defioptimization/defi-service/server"
	"github.com/defioptimization/shared/database"
)

func main() {
	// Initialize database
	if err := database.InitDatabase(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.CloseDatabase()

	// Initialize protocol clients
	ethRPC := os.Getenv("ETH_RPC_URL")
	baseRPC := os.Getenv("BASE_RPC_URL")

	if ethRPC == "" {
		log.Fatal("ETH_RPC_URL environment variable is required")
	}
	if baseRPC == "" {
		log.Fatal("BASE_RPC_URL environment variable is required")
	}

	// Initialize protocol managers
	protocolManager := protocols.NewManager(ethRPC, baseRPC)

	// Start HTTP server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	srv := server.NewServer(protocolManager)
	log.Printf("DeFi Service starting on port %s", port)
	if err := srv.Start(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

