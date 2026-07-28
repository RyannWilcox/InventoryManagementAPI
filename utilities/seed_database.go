package utilities

import (
	"inv-backend/models"
	"log"

	"gorm.io/gorm"
)

func toPtr[T any](v T) *T {
	return &v
}

func SeedDatabase(db *gorm.DB) {
	var count int64
	db.Model(&models.Item{}).Count(&count)
	if count == 0 {
		items := []models.Item{
			{Name: "Laptop", Stock: toPtr(10), Price: toPtr(999.99)},
			{Name: "Smartphone", Stock: toPtr(20), Price: toPtr(699.99)},
			{Name: "Headphones", Stock: toPtr(15), Price: toPtr(199.99)},
			{Name: "Keyboard", Stock: toPtr(25), Price: toPtr(89.99)},
			{Name: "Mouse", Stock: toPtr(30), Price: toPtr(49.99)},
			{Name: "Monitor", Stock: toPtr(12), Price: toPtr(299.99)},
			{Name: "Webcam", Stock: toPtr(18), Price: toPtr(79.99)},
			{Name: "Printer", Stock: toPtr(7), Price: toPtr(149.99)},
			{Name: "Tablet", Stock: toPtr(5), Price: toPtr(399.99)},
			{Name: "Smartwatch", Stock: toPtr(14), Price: toPtr(249.99)},
			{Name: "External Hard Drive", Stock: toPtr(8), Price: toPtr(119.99)},
			{Name: "USB Flash Drive", Stock: toPtr(50), Price: toPtr(19.99)},
			{Name: "Router", Stock: toPtr(6), Price: toPtr(89.99)},
			{Name: "Projector", Stock: toPtr(3), Price: toPtr(499.99)},
			{Name: "Bluetooth Speaker", Stock: toPtr(22), Price: toPtr(129.99)},
			{Name: "Gaming Console", Stock: toPtr(11), Price: toPtr(499.99)},
			{Name: "Camera", Stock: toPtr(4), Price: toPtr(599.99)},
			{Name: "Fitness Tracker", Stock: toPtr(16), Price: toPtr(99.99)},
			{Name: "Drone", Stock: toPtr(2), Price: toPtr(899.99)},
			{Name: "VR Headset", Stock: toPtr(9), Price: toPtr(399.99)},
		}

		db.Create(&items)
		log.Println("Database seeded with 20 sample items.")
	} else {
		log.Println("Database already contains data, skipping seeding.")
	}
}
