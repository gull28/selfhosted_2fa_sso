package routes

import (
	"selfhosted_2fa_sso/controllers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterServiceRoutes(router *gin.Engine, db *gorm.DB) {
	serviceControllerController := controllers.GetServiceController(db)

	userRoutes := router.Group("/service")
	{
		userRoutes.POST("/bind", serviceControllerController.BindService)
		userRoutes.POST("/unbind", serviceControllerController.UnbindService)
	}
}
