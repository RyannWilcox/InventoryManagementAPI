package routes

import (
	"inv-backend/controllers"
	"inv-backend/utilities"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InjectDB(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	}
}

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

func InitRoutes(router *gin.Engine, db *gorm.DB, rateLimiter *utilities.TokenBucket) {
	router.Use(InjectDB(db))
	router.Use(RateLimit(rateLimiter))

	router.Use(func(c *gin.Context) {
		c.Set("db", db)

		if !rateLimiter.Allow() {
			c.JSON(http.StatusTooManyRequests, utilities.HTTPError{
				Code:    http.StatusTooManyRequests,
				Message: "Too many requests",
			})
			c.Abort()
			return
		}
		c.Next()
	})

	v1 := router.Group("api/v1")
	{
		items := v1.Group("/items")
		{
			items.GET("", controllers.GetItems)
			items.GET(":id", controllers.GetItem)
			items.POST("", controllers.CreateItem)
			items.PUT(":id", controllers.UpdateItem)
			items.DELETE(":id", controllers.DeleteItem)
		}
	}
}
