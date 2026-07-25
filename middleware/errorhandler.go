package middleware

import (
	"errors"
	"inv-backend/utilities"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// no errors. continue on.
		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err

		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			c.JSON(http.StatusNotFound, utilities.HTTPError{
				Code:    http.StatusNotFound,
				Message: "Item could not be found.",
			})
		default:
			c.JSON(http.StatusInternalServerError, utilities.HTTPError{
				Code:    http.StatusInternalServerError,
				Message: "An internal error occurred.",
			})
		}
	}
}
