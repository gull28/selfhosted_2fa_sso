package routes

import (
	"selfhosted_2fa_sso/controllers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterSessionRoutes(sessionRoutes *gin.RouterGroup, db *gorm.DB) {
	sessionController := controllers.GetSessionController(db)

	sessionRoutes.GET("create", sessionController.Index)
	sessionRoutes.POST("create", sessionController.Create)
}
