package routes

import (
	"selfhosted_2fa_sso/config"
	"selfhosted_2fa_sso/controllers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterUserRoutes(router *gin.Engine, db *gorm.DB, cfg *config.Config) {
	userController := controllers.GetUserController(db, cfg)

	userRoutes := router.Group("/users")
	{
		userRoutes.POST("", userController.Create)
		userRoutes.POST("verify", userController.Verify)
		userRoutes.POST("session/check", userController.CheckSession)
	}
}
