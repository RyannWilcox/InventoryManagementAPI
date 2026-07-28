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

	validSortOptions := map[string]bool{
		"name":  true,
		"stock": true,
		"price": true,
	}
	sort := strings.Split(c.DefaultQuery("sort", "name"), ",")
	order := strings.Split(c.DefaultQuery("order", "asc"), ",")

	var orderParams []string
	for i, field := range sort {
		if !validSortOptions[field] {
			continue // skip over invalid sort params
		}
		direction := "asc"
		if i < len(order) && (order[i] == "asc" || order[i] == "desc") {
			direction = order[i]
		}
		orderParams = append(orderParams, fmt.Sprintf("%s %s", field, direction))

	}

	// Didn't find any valid params
	// append a default value
	if len(orderParams) == 0 {
		orderParams = append(orderParams, "name asc")
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
	minStock := c.DefaultQuery("min_stock", "0")

	if name != "" {
		query = query.Where("name = ?", name)
	}
	if stock, err := strconv.Atoi(minStock); err == nil {
		query = query.Where("stock >= ?", stock)
	}

	return query
}

// Will add pagination to the query
func applyPagination(c *gin.Context, query *gorm.DB) (*gorm.DB, utilities.PaginationMetaData, error) {

	limitParam := c.DefaultQuery("limit", "20")
	offsetParam := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitParam)
	if err != nil {
		return nil, utilities.PaginationMetaData{}, err
	}
	if limit <= 0 {
		return nil, utilities.PaginationMetaData{}, fmt.Errorf("%w", utilities.LimitError)
	}

	offset, err := strconv.Atoi(offsetParam)
	if err != nil {
		return nil, utilities.PaginationMetaData{}, err
	}
	if offset < 0 {
		return nil, utilities.PaginationMetaData{}, fmt.Errorf("%w", utilities.OffsetError)
	}

	// Add pagination to the query
	query = query.Limit(limit).Offset(offset)

	return query,
		utilities.PaginationMetaData{
			Limit:  limit,
			Offset: offset,
		}, nil
}

// GetItems godoc
// @Summary      List inventory items
// @Description  Retrieves a paginated, filterable, sortable list of inventory items
// @Tags         inventory
// @Accept       json
// @Produce      json
// @Param        sort      query     string  false  "Comma-separated fields to sort by (name, stock, price, created_at)"
// @Param        order     query     string  false  "Comma-separated sort directions matching `sort`: asc or desc"
// @Param        name      query     string  false  "Filter items by exact name match"
// @Param        min_stock query     int     false  "Filter items with stock greater than or equal to this value"
// @Param        limit     query     int     false  "Max number of items to return (default 20)"
// @Param        offset    query     int     false  "Number of items to skip (default 0)"
// @Success      200  {object}  utilities.PaginatedItemsResponse
// @Failure      400  {object}  utilities.HTTPError
// @Failure      500  {object}  utilities.HTTPError
// @Router       /inventory [get]
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

	c.JSON(http.StatusOK,
		utilities.PaginatedItemsResponse{
			Meta:  metaData,
			Items: items,
		})
}

// GetItem godoc
// @Summary      Get an inventory item
// @Description  Retrieves the details of a single inventory item by its UUID
// @Tags         inventory
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Item ID (UUID)"
// @Success      200  {object}  models.Item
// @Failure      404  {object}  utilities.HTTPError
// @Failure      500  {object}  utilities.HTTPError
// @Router       /inventory/{id} [get]
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

// CreateItem godoc
// @Summary      Create an inventory item
// @Description  Creates a new inventory item from the request body
// @Tags         inventory
// @Accept       json
// @Produce      json
// @Param        item  body      models.Item  true  "Item to create"
// @Success      201  {object}  models.Item
// @Failure      400  {object}  utilities.HTTPError
// @Failure      500  {object}  utilities.HTTPError
// @Router       /inventory [post]
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

	c.JSON(http.StatusCreated, item)
}

// UpdateItem godoc
// @Summary      Update an inventory item
// @Description  Updates an existing inventory item with the fields provided in the request body
// @Tags         inventory
// @Accept       json
// @Produce      json
// @Param        id    path      string       true  "Item ID (UUID)"
// @Param        item  body      models.Item  true  "Fields to update"
// @Success      200  {object}  models.Item
// @Failure      400  {object}  utilities.HTTPError
// @Failure      404  {object}  utilities.HTTPError
// @Failure      500  {object}  utilities.HTTPError
// @Router       /inventory/{id} [put]
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

	c.JSON(http.StatusOK, item)
}

// DeleteItem godoc
// @Summary      Delete an inventory item
// @Description  Deletes an existing inventory item by its UUID
// @Tags         inventory
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Item ID (UUID)"
// @Success      200  {object}  utilities.MessageResponse
// @Failure      404  {object}  utilities.HTTPError
// @Failure      500  {object}  utilities.HTTPError
// @Router       /inventory/{id} [delete]
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

	c.JSON(http.StatusNoContent, utilities.MessageResponse{
		Message: "Item successfully deleted.",
	})
}
