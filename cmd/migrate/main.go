package main

import (
	"log"
	"to-do-list/internal/models"
	"to-do-list/pkg/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	dsn := conf.GetDSN()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	err = db.AutoMigrate(&models.Task{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	log.Println("Database migrations completed successfully")
}
