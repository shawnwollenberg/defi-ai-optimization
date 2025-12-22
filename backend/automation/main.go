package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/defioptimization/automation/engine"
	"github.com/defioptimization/shared/database"
)

func main() {
	// Initialize database
	if err := database.InitDatabase(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.CloseDatabase()

	// Initialize automation engine
	defiServiceURL := os.Getenv("DEFI_SERVICE_URL")
	if defiServiceURL == "" {
		defiServiceURL = "http://localhost:8081"
	}

	walletServiceURL := os.Getenv("WALLET_SERVICE_URL")
	if walletServiceURL == "" {
		walletServiceURL = "http://localhost:8082"
	}

	mlServiceURL := os.Getenv("ML_SERVICE_URL")
	if mlServiceURL == "" {
		mlServiceURL = "http://localhost:8001"
	}

	automationEngine := engine.NewEngine(defiServiceURL, walletServiceURL, mlServiceURL)

	// Start the engine
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Shutting down automation engine...")
		cancel()
	}()

	// Start monitoring loop
	interval := 30 * time.Second // Check every 30 seconds
	if err := automationEngine.Start(ctx, interval); err != nil {
		log.Fatalf("Failed to start automation engine: %v", err)
	}
}

