package main

import (
	"log"

	"github.com/defioptimization/shared/database"
)

func main() {
	if err := database.InitDatabase(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	log.Println("Migrations completed successfully")
}

