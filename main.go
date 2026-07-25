package main

import (
	"inv-backend/models"
	"inv-backend/routes"
	"inv-backend/utilities"
	"log"

	"github.com/gin-gonic/gin"
)

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

	routes.InitRoutes(router, db)

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
