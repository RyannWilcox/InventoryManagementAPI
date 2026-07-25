package utilities

import (
	"inv-backend/models"
	"log"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func SeedDatabase(db *gorm.DB) {
	var count int64
	db.Model(&models.Item{}).Count(&count)
	if count == 0 {
		items := []models.Item{
			{ID: uuid.New(), Name: "Laptop", Stock: 10, Price: 999.99},
			{ID: uuid.New(), Name: "Smartphone", Stock: 20, Price: 699.99},
			{ID: uuid.New(), Name: "Headphones", Stock: 15, Price: 199.99},
			{ID: uuid.New(), Name: "Keyboard", Stock: 25, Price: 89.99},
			{ID: uuid.New(), Name: "Mouse", Stock: 30, Price: 49.99},
			{ID: uuid.New(), Name: "Monitor", Stock: 12, Price: 299.99},
			{ID: uuid.New(), Name: "Webcam", Stock: 18, Price: 79.99},
			{ID: uuid.New(), Name: "Printer", Stock: 7, Price: 149.99},
			{ID: uuid.New(), Name: "Tablet", Stock: 5, Price: 399.99},
			{ID: uuid.New(), Name: "Smartwatch", Stock: 14, Price: 249.99},
			{ID: uuid.New(), Name: "External Hard Drive", Stock: 8, Price: 119.99},
			{ID: uuid.New(), Name: "USB Flash Drive", Stock: 50, Price: 19.99},
			{ID: uuid.New(), Name: "Router", Stock: 6, Price: 89.99},
			{ID: uuid.New(), Name: "Projector", Stock: 3, Price: 499.99},
			{ID: uuid.New(), Name: "Bluetooth Speaker", Stock: 22, Price: 129.99},
			{ID: uuid.New(), Name: "Gaming Console", Stock: 11, Price: 499.99},
			{ID: uuid.New(), Name: "Camera", Stock: 4, Price: 599.99},
			{ID: uuid.New(), Name: "Fitness Tracker", Stock: 16, Price: 99.99},
			{ID: uuid.New(), Name: "Drone", Stock: 2, Price: 899.99},
			{ID: uuid.New(), Name: "VR Headset", Stock: 9, Price: 399.99},
		}

		db.Create(&items)
		log.Println("Database seeded with 20 sample items.")
	} else {
		log.Println("Database already contains data, skipping seeding.")
	}
}
