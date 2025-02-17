package controllers

import (
	"net/http"
	"selfhosted_2fa_sso/config"
	"selfhosted_2fa_sso/models"
	"selfhosted_2fa_sso/requests"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BindController struct {
	db  *gorm.DB
	cfg *config.Config
}

func GetBindController(db *gorm.DB) *BindController {
	return &BindController{db: db}
}

func (bc *BindController) Create(c *gin.Context) {
	var req requests.CreateBindRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	user, err := models.GetUserByUsername(bc.db, req.Username)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}

	service, err := models.GetServiceByID(bc.db, req.ServiceID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Service not found"})
		return
	}

	// delete all bind requests

	if err := models.DeleteBindRequestsForService(bc.db, user.ID, service.ID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Internal server error"})
		return
	}

	// create user service link

	userServiceLink := models.UserServiceLink{
		ServiceUserID: req.ServiceUserID,
		ValidUntil:    time.Now().Add(time.Duration(bc.cfg.Auth.ValidFor) * time.Minute),
	}

	if err := userServiceLink.CreateUserServiceLink(bc.db); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create User service link"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully created user service link"})

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
