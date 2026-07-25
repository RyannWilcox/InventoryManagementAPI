package routes

import (
	"inv-backend/controllers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Create API routes
func InitRoutes(router *gin.Engine, db *gorm.DB) {
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
