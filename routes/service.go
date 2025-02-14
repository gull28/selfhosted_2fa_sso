package routes

import (
	"selfhosted_2fa_sso/controllers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterServiceRoutes(serviceRoutes *gin.RouterGroup, db *gorm.DB) {
	serviceController := controllers.GetServiceController(db)

	serviceRoutes.GET(":id", serviceController.Fetch)
	serviceRoutes.GET("create", serviceController.Index)
	serviceRoutes.POST("create", serviceController.Create)
	serviceRoutes.POST("bind", serviceController.BindServiceTo2fa)
	serviceRoutes.DELETE(":id", serviceController.Delete)
	// userRoutes.POST("/unbind", serviceControllerController)
}
