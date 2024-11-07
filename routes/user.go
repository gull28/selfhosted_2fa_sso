package routes

import (
	"selfhosted_2fa_sso/controllers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterUserRoutes(router *gin.Engine, db *gorm.DB) {
	userController := controllers.GetUserController(db)

	userRoutes := router.Group("/users")
	{
		userRoutes.POST("/", userController.CreateUser)
		userRoutes.POST("verify", userController.VerifyUser)
	}
}
