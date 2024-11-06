package controllers

import (
	"gorm.io/gorm"
)

type ServiceController struct {
	db *gorm.DB
}

func GetServiceController(db *gorm.DB) *ServiceController {
	return &ServiceController{DB: db}
}

func (sc *ServiceController) BindServiceTo2fa(c *gin.Config) {
	var userServiceLink models.UserServiceLink
}
