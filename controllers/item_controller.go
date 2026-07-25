package controllers

import (
	"inv-backend/models"
	"inv-backend/utilities"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetItems(c *gin.Context) {}

func GetItem(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id := c.Param("id")

	var item models.Item
	if result := db.First(&item, "id = ?", id); result.Error != nil {
		c.Error(result.Error)
		return
	}

	c.JSON(http.StatusOK, item)
}

func CreateItem(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var item models.Item

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

func UpdateItem(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id := c.Param("id")

	var item models.Item
	if result := db.First(&item, "id = ?", id); result.Error != nil {
		c.Error(result.Error)
		return
	}

	var inputData models.Item
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
