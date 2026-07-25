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
		c.JSON(http.StatusNotFound, utilities.HTTPError{
			Code:    http.StatusNotFound,
			Message: "Item could not be found.",
		})
		return
	}

	c.JSON(http.StatusOK, item)
}

func CreateItem(c *gin.Context) {}

func UpdateItem(c *gin.Context) {}

func DeleteItem(c *gin.Context) {}
