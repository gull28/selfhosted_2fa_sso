package controllers

import (
	"fmt"
	"net/http"
	"selfhosted_2fa_sso/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ServiceController struct {
	db *gorm.DB
}

type BindRequest struct {
	ServiceID  string `json:"serviceId" binding:"required"`
	UserID     string `json:"userId" binding:"required"`     // userId created in the service server
	AuthUserID uint   `json:"authUserId" binding:"required"` // userId created in the 2fa server
}

type CreateServiceRequest struct {
	ID          uint   `json:"serviceId" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
}

func GetServiceController(db *gorm.DB) *ServiceController {
	return &ServiceController{db: db}
}

func (sc *ServiceController) Create(c *gin.Context) {
	var service models.Service2fa
	if err := c.ShouldBindJSON(&service); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	if err := service.Create(sc.db); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create service"})
		return
	}

	data, err := models.GetAllServices(sc.db)

	if err != nil {
		c.JSON(http.StatusCreated, []string{})
		return
	}
	c.JSON(http.StatusCreated, data)
}

func (sc *ServiceController) Delete(c *gin.Context) {
	idParam := c.Param("id")

	if idParam == "" {

		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	idUint64, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		fmt.Printf("%v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	id := uint(idUint64)

	if err := models.DeleteService(sc.db, uint(id)); err != nil {
		fmt.Printf("%v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
	}

	data, err := models.GetAllServices(sc.db)

	if err != nil {
		fmt.Printf("%v", err)
		c.JSON(http.StatusCreated, []string{})
		return
	}
	c.JSON(http.StatusOK, data)
}

func (sc *ServiceController) Index(c *gin.Context) {
	data, err := models.GetAllServices(sc.db)
	// items := []string{"hello", "world", "123312"}
	if err != nil {
		// render error page
		return
	}

	c.HTML(http.StatusOK, "service.html", gin.H{
		"items": data,
	})
}

func (sc *ServiceController) BindServiceTo2fa(c *gin.Context) {
	var bindRequest BindRequest

	if err := c.ShouldBindJSON(&bindRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	service, err := models.GetServiceByID(sc.db, bindRequest.ServiceID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Service not found"})
		return
	}

	user, err := models.GetUserByID(sc.db, bindRequest.AuthUserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}

	var userServiceLink = models.UserServiceLink{
		UserID:       bindRequest.UserID,
		Service2faID: service.ID,
		User2faID:    user.ID,
	}

	if err := sc.db.Create(&userServiceLink).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create link"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Service linked to 2FA successfully", "link": userServiceLink})
}
