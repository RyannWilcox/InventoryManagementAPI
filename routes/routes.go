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
		inventory := v1.Group("/inventory")
		{
			inventory.GET("", controllers.GetItems)
			inventory.GET(":id", controllers.GetItem)
			inventory.POST("", controllers.CreateItem)
			inventory.PUT(":id", controllers.UpdateItem)
			inventory.DELETE(":id", controllers.DeleteItem)
		}
	}
}
