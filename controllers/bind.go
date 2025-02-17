package controllers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BindController struct {
	db *gorm.DB
}

func GetBindController(db *gorm.DB) *BindController {
	return &BindController{db: db}
}

func (bc *BindController) Create(c *gin.Context) {
	//

}

func (bc *BindController) Fetch(c *gin.Context) {
	//

}

func (bc *BindController) Accept(c *gin.Context) {
	//
}

func (bc *BindController) Decline(c *gin.Context) {
	//
}
