package database

import (
	"fmt"
	"log"
	"os"

	"github.com/defioptimization/shared/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDatabase initializes the database connection and runs migrations
func InitDatabase() error {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			getEnv("DB_HOST", "localhost"),
			getEnv("DB_USER", "defi_user"),
			getEnv("DB_PASSWORD", "defi_password"),
			getEnv("DB_NAME", "defi_optimization"),
			getEnv("DB_PORT", "5432"),
		)
	}

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("Database connection established")

	// Run migrations
	if err := runMigrations(); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

// runMigrations runs database migrations
func runMigrations() error {
	return DB.AutoMigrate(
		&models.User{},
		&models.Portfolio{},
		&models.Position{},
		&models.PortfolioSnapshot{},
		&models.AutomationRule{},
		&models.Transaction{},
		&models.Subscription{},
	)
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// CloseDatabase closes the database connection
func CloseDatabase() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

