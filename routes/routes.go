package routes

import (
	"inv-backend/controllers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitRoutes(router *gin.Engine, db *gorm.DB) {
	router.Use(func(c *gin.Context) {
		c.Set("db", db)
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
