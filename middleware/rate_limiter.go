package middleware

import (
	"inv-backend/utilities"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RateLimit(rateLimiter *utilities.TokenBucket) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !rateLimiter.Allow() {
			c.JSON(http.StatusTooManyRequests, utilities.HTTPError{
				Code:    http.StatusTooManyRequests,
				Message: "Too many requests",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
