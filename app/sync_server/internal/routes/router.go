package routes

import (
	"github.com/Mushroom-MSL1L/UltimateDesktopPet/app/sync_server/docs"
	. "github.com/Mushroom-MSL1L/UltimateDesktopPet/app/sync_server/internal/controllers"
	. "github.com/Mushroom-MSL1L/UltimateDesktopPet/app/sync_server/internal/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SwagSetUp(router *gin.Engine) {
	url := ginSwagger.URL("/swagger/doc.json")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
}

func SwagInfoSetup(host string) {
	docs.SwaggerInfo.Title = "Ultimate Desktop Pet Sync Server API"
	docs.SwaggerInfo.Description = "API documentation for Ultimate Desktop Pet Sync Server"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = host
}

func SetupRouter(host string) *gin.Engine {
	r := gin.Default()
	SwagSetUp(r)
	SwagInfoSetup(host)

	r.POST("/login", Login)
	r.GET("/profile", AuthMiddleware(), Profile)
	r.GET("/health", Health)
	return r
}
