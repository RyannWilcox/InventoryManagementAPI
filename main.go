package main

import (
	"inv-backend/middleware"
	"inv-backend/models"
	"inv-backend/routes"
	"inv-backend/utilities"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// @title Inventory Management API
// @version 1.0
// @description This a backend API for an inventory management system
// @host localhost:8080
// @BasePath /api/v1

func main() {
	db, err := utilities.Connect()

	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}

	sqlDB, err := db.DB()

	if err != nil {
		log.Fatalf("Failed to get database instance: %v", err)
	}

	defer sqlDB.Close()

	if err := db.AutoMigrate(&models.Item{}); err != nil {
		log.Fatalf("Failed to automigrate database: %v", err)
	}

	router := gin.Default()

	// Setup middleware
	router.Use(middleware.InjectDB(db))
	router.Use(middleware.RateLimit(utilities.NewTokenBucket(5, time.Second)))
	router.Use(middleware.ErrorHandler())

	// create endpoint routes
	routes.InitRoutes(router, db)

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
