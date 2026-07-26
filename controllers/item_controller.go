package controllers

import (
	"fmt"
	"inv-backend/models"
	"inv-backend/utilities"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Will add sorting options to the query
func applySorting(c *gin.Context, query *gorm.DB) *gorm.DB {
	/*
	 Sorting
	 Your API supports sorting by fields such as
	 name, stock, and price (in either ascending *or *descending order).
	*/
	sort := strings.Split(c.DefaultQuery("sort", "created_at"), ",")
	order := strings.Split(c.DefaultQuery("order", "asc"), ",")

	var orderParams []string
	for i, field := range sort {
		direction := "asc"
		if i < len(order) && (order[i] == "asc" || order[i] == "desc") {
			direction = order[i]
		}
		orderParams = append(orderParams, fmt.Sprintf("%s %s", field, direction))

	}

	return query.Order(strings.Join(orderParams, ", "))
}

// Will add filter optins to the query
func applyFiltering(c *gin.Context, query *gorm.DB) *gorm.DB {
	/*
	 filtering:
	 Your users can also filter results by
	 the item's name or minimum stock levels.
	*/
	name := c.Query("name")
	minStock := c.DefaultQuery("minStock", "0")

	if name != "" {
		query = query.Where("name = ?", name)
	}
	if stock, err := strconv.Atoi(minStock); err == nil {
		query = query.Where("stock >= ?", stock)
	}

	return query
}

type PaginationMetaData struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

// Will add pagination to the query
func applyPagination(c *gin.Context, query *gorm.DB) (*gorm.DB, PaginationMetaData, error) {

	limitParam := c.DefaultQuery("limit", "20")
	offsetParam := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitParam)
	if err != nil || limit <= 0 {
		return nil, PaginationMetaData{}, err
	}

	offset, err := strconv.Atoi(offsetParam)
	if err != nil || offset < 0 {
		return nil, PaginationMetaData{}, err
	}

	// Add pagination to the query
	query = query.Limit(limit).Offset(offset)

	return query, PaginationMetaData{
		Limit:  limit,
		Offset: offset,
	}, nil
}

// Retieve all inventory items
func GetItems(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var items []models.Item

	query := db.Model(&models.Item{})

	query = applyFiltering(c, query)
	query = applySorting(c, query)
	query, metaData, paginationError := applyPagination(c, query)

	if paginationError != nil {
		c.Error(paginationError)
		return
	}

	// Query the database after filterting, sorting and pagination
	if err := query.Find(&items).Error; err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"meta":  metaData,
		"items": items,
	})
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
