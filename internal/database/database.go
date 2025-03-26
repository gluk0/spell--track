package database

import (
	"event-tracking-service/internal/models"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		getEnvOrDefault("DB_HOST", "localhost"),
		getEnvOrDefault("DB_USER", "postgres"),
		getEnvOrDefault("DB_PASSWORD", "postgres"),
		getEnvOrDefault("DB_NAME", "event_tracking"),
		getEnvOrDefault("DB_PORT", "5432"),
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Drop existing tables to ensure clean state
	err = DB.Migrator().DropTable(&models.Event{}, &models.Case{})
	if err != nil {
		log.Fatal("Failed to drop tables:", err)
	}

	// Create tables with new schema
	err = DB.AutoMigrate(&models.Event{}, &models.Case{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
