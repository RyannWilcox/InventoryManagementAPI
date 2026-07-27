package utilities

import "inv-backend/models"

// MessageResponse defines the structure for success messages
type MessageResponse struct {
	Message string `json:"message" example:"Item deleted"`
}

// HTTPError defines the structure for error responses
type HTTPError struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"Invalid input"`
}

// Meta data provided by paginated items
type PaginationMetaData struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

// Represents the response body for a paginated list of items.
type PaginatedItemsResponse struct {
	Meta  PaginationMetaData `json:"meta"`
	Items []models.Item      `json:"items"`
}
