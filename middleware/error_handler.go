package middleware

import (
	"encoding/json"
	"errors"
	"inv-backend/utilities"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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
		case errors.As(err, new(*json.SyntaxError)), errors.As(err, new(*json.UnmarshalTypeError)):
			c.JSON(http.StatusBadRequest, utilities.HTTPError{
				Code:    http.StatusBadRequest,
				Message: "Invalid JSON provided in the request body.",
			})
		case errors.Is(err, utilities.LimitError), errors.Is(err, utilities.OffsetError):
			c.JSON(http.StatusBadRequest, utilities.HTTPError{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			})
		case errors.As(err, new(*strconv.NumError)):
			c.JSON(http.StatusBadRequest, utilities.HTTPError{
				Code:    http.StatusBadRequest,
				Message: "An invalid numeric value was provided.",
			})
		case errors.As(err, new(validator.ValidationErrors)):
			c.JSON(http.StatusBadRequest, utilities.HTTPError{
				Code:    http.StatusBadRequest,
				Message: "Required fields are either missing or invalid",
			})
		default:
			c.JSON(http.StatusInternalServerError, utilities.HTTPError{
				Code:    http.StatusInternalServerError,
				Message: "An internal error occurred.",
			})
		}
	}
}
