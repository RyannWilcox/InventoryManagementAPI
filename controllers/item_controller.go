package controllers

import (
	"inv-backend/models"
	"inv-backend/utilities"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Retieve all inventory items
func GetItems(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var items []models.Item

	//TODO: Add pagination

	//TODO: Add sorting

	//TODO: Add filtering

	// Fetch all items from the database.
	if err := db.Preload("Item").Find(&items).Error; err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, items)
}

// Will retrieve a single inventory item by its id.
func GetItem(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id := c.Param("id")

	var item models.Item
	// Look up the item by its primary key
	if result := db.First(&item, "id = ?", id); result.Error != nil {
		c.Error(result.Error)
		return
	}

	c.JSON(http.StatusOK, item)
}

// Will create a new item based on the data passed in by the request body
func CreateItem(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var item models.Item
	// Bind and validate the incoming JSON payload into an Item
	if err := c.ShouldBindJSON(&item); err != nil {
		c.Error(err)
		return
	}

	// Begin the transaction
	tx := db.Begin()

	if tx.Error != nil {
		c.Error(tx.Error)
		return
	}

	// Insert the new item (BeforeCreate hook assigns its UUID)
	if err := tx.Create(&item).Error; err != nil {
		tx.Rollback()
		c.Error(err)
		return
	}

	// Attempt to apply the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.Error(err)
		return
	}

}

// Will update an existing inventory item with the fields provided in the request body.
func UpdateItem(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id := c.Param("id")

	var item models.Item
	// Ensure the item exists before attempting to update it
	if result := db.First(&item, "id = ?", id); result.Error != nil {
		c.Error(result.Error)
		return
	}

	var inputData models.Item
	// Bind and validate the incoming JSON payload with the updated fields
	if err := c.ShouldBindJSON(&inputData); err != nil {
		c.Error(err)
		return
	}

	// Begin the transaction
	tx := db.Begin()

	if tx.Error != nil {
		c.Error(tx.Error)
		return
	}

	// Apply the updates to the existing item and reload it into `item`.
	if err := tx.Model(&item).Updates(inputData).Scan(&item).Error; err != nil {
		tx.Rollback()
		c.Error(err)
		return
	}

	// Attempt to apply the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.Error(err)
		return
	}
}

// Will delete an existing inventory item by its id.
func DeleteItem(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id := c.Param("id")

	var item models.Item
	if result := db.First(&item, "id = ?", id); result.Error != nil {
		c.Error(result.Error)
		return
	}

	// Begin the transaction
	tx := db.Begin()

	if tx.Error != nil {
		c.Error(tx.Error)
		return
	}

	if err := tx.Delete(&item).Error; err != nil {
		tx.Rollback()
		c.Error(err)
		return
	}

	// Attempt to apply the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, utilities.MessageResponse{
		Message: "Item successfully deleted.",
	})
}
