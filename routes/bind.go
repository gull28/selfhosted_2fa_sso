package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func BindRoutes(bindRoutes *gin.RouterGroup, db *gorm.DB) {
	bindRoutes.GET("")
}
