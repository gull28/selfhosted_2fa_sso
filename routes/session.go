package routes

import (
	"selfhosted_2fa_sso/config"
	"selfhosted_2fa_sso/controllers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterSessionRoutes(sessionRoutes *gin.RouterGroup, db *gorm.DB, cfg *config.Config) {
	sessionController := controllers.GetSessionController(db, cfg)

	sessionRoutes.GET("create", sessionController.Index)
	sessionRoutes.POST("create", sessionController.Create)
	sessionRoutes.DELETE("logout", sessionController.Delete)
}
