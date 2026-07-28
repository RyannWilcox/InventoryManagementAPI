package routes

import (
	"inv-backend/controllers"

	"inv-backend/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

// Create API routes
func InitRoutes(router *gin.Engine, db *gorm.DB) {
	v1BasePath := "api/v1"

	// Setup swagger documentation endpoint
	docs.SwaggerInfo.BasePath = "/" + v1BasePath
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := router.Group(v1BasePath)
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
