package main

import (
	"log"
	"os"

	"github.com/defioptimization/wallet/server"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	srv := server.NewServer()
	log.Printf("Wallet Service starting on port %s", port)
	if err := srv.Start(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

