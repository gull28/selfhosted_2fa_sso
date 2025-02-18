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

func GetBindController(db *gorm.DB, cfg *config.Config) *BindController {
	return &BindController{db: db, cfg: cfg}
}

func (bc *BindController) Create(c *gin.Context) {
	var req requests.CreateBindRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// check user
	user, err := models.GetUserByUsername(bc.db, req.Username)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}

	// check service
	service, err := models.GetServiceByID(bc.db, req.ServiceID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Service not found"})
		return
	}

	// check if there already isnt an active bind request

	bindRequest := models.BindRequest{
		ValidUntil:    time.Now().Add(time.Duration(bc.cfg.Auth.ValidFor) * time.Minute),
		ServiceUserID: req.ServiceUserID,
		User2faID:     user.ID,
		Service2faID:  service.ID,
	}

	if err := bindRequest.Create(bc.db); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error creating bindRequest"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully created user service link"})
}

func (bc *BindController) Accept(c *gin.Context) {
	var req requests.ActionBindRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	bindRequest, err := models.GetBindRequestByID(bc.db, req.BindRequestID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}

	user, err := models.GetUserByID(bc.db, bindRequest.User2faID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}

	service, err := models.GetServiceByID(bc.db, bindRequest.Service2faID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Service not found"})
		return
	}

	if err := models.DeleteBindRequestsForService(bc.db, user.ID, service.ID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Internal server error"})
		return
	}

	// create user service link

	userServiceLink := models.UserServiceLink{
		ServiceUserID: bindRequest.Service2faID,
		ValidUntil:    time.Now().Add(time.Duration(bc.cfg.Auth.ValidFor) * time.Minute),
	}

	if err := userServiceLink.CreateUserServiceLink(bc.db); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create User service link"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully accepted bind request"})

}

func (bc *BindController) Fetch(c *gin.Context) {
	var req requests.FetchActiveBindRequests

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	bindRequests := models.GetBindRequestsByUserID(bc.db, req.UserID)

	c.JSON(http.StatusOK, gin.H{"bindRequests": bindRequests})

}

func (bc *BindController) Decline(c *gin.Context) {
	var req requests.ActionBindRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	bindRequest, err := models.GetBindRequestByID(bc.db, req.BindRequestID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}

	user, err := models.GetUserByID(bc.db, bindRequest.User2faID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}

	service, err := models.GetServiceByID(bc.db, bindRequest.Service2faID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Service not found"})
		return
	}

	// delete all bind requests

	if err := models.DeleteBindRequestsForService(bc.db, user.ID, service.ID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Bind request declined successfully"})
}
