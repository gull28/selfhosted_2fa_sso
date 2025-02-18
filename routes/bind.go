package routes

import (
	"selfhosted_2fa_sso/config"
	"selfhosted_2fa_sso/controllers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterBindRoutes(bindRoutes *gin.RouterGroup, db *gorm.DB, cfg *config.Config) {
	bindRequestController := controllers.GetBindController(db, cfg)

	bindRoutes.GET("", bindRequestController.Fetch)
	bindRoutes.POST("", bindRequestController.Create)
	bindRoutes.POST("/accept", bindRequestController.Accept)
	bindRoutes.DELETE("/decline", bindRequestController.Decline)
}
