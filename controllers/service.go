package controllers

import (
	"fmt"
	"net/http"
	"selfhosted_2fa_sso/models"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ServiceController struct {
	db *gorm.DB
}

type BindRequest struct {
	ServiceID string `json:"serviceId" binding:"required"`
	// used for hooks in service so that user doesnt have to enter an id or username on each totp request
	UserID   string `json:"userId" binding:"required"`   // userId created in the service server
	Username string `json:"username" binding:"required"` // userId created in the 2fa server
}

type ServiceItem struct {
	ServiceID   string    `json:"serviceId" binding:"required"`
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"username" binding:"required"`
	Enabled     bool      `json:"enabled"`
	ValidUntil  time.Time `json:"validUntil"`
}

// type CreateServiceRequest struct {
// 	ID          string `json:"serviceId" binding:"required"`
// 	Name        string `json:"name" binding:"required"`
// 	Description string `json:"description" binding:"required"`
// }

func GetServiceController(db *gorm.DB) *ServiceController {
	return &ServiceController{db: db}
}

func (sc *ServiceController) Create(c *gin.Context) {
	var service models.Service2fa
	if err := c.ShouldBindJSON(&service); err != nil {
		fmt.Printf("%v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data "})
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

func (sc *ServiceController) Fetch(c *gin.Context) {

	userID := c.Param("id")

	userServiceLinks := models.FetchUserServiceLinks(sc.db, userID)

	nonBindedServices, err := models.GetAllServices(sc.db)

	if err != nil {
		fmt.Printf("error 1 %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "User with given userId not found"})
	}

	serviceItems := []ServiceItem{}

	for _, v := range userServiceLinks {
		serviceItems = append(serviceItems, ServiceItem{
			ServiceID:   v.Service2faID,
			Name:        v.Service2fa.Name,
			Description: v.Service2fa.Description,
			ValidUntil:  v.ValidUntil,
			Enabled:     v.Enabled,
		})
	}

	nonBindedServiceItems := []ServiceItem{}

	for _, nonBindedService := range nonBindedServices {
		nonBindedServiceItems = append(nonBindedServiceItems, ServiceItem{
			ServiceID:   nonBindedService.ID,
			Name:        nonBindedService.Name,
			Description: nonBindedService.Description,
		})
	}

	fmt.Printf("%v serviceItmes", serviceItems)
	c.JSON(http.StatusOK, gin.H{"services": serviceItems, "nonBoundServices": nonBindedServiceItems})
}

func (sc *ServiceController) Delete(c *gin.Context) {
	idParam := c.Param("id")

	if idParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := models.DeleteService(sc.db, idParam); err != nil {
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

	if err != nil {
		return
	}

	c.HTML(http.StatusOK, "service.html", gin.H{
		"items": data,
	})
}

func (sc *ServiceController) BindServiceTo2fa(c *gin.Context) {
	var bindRequest BindRequest

	if err := c.ShouldBindJSON(&bindRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request data"})
		return
	}

	service, err := models.GetServiceByID(sc.db, bindRequest.ServiceID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Service not found"})
		return
	}

	user, err := models.GetUserByUsername(sc.db, bindRequest.Username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User not found"})
		return
	}

	var userServiceLink = models.UserServiceLink{
		ValidUntil:    time.Now(),
		ServiceUserID: bindRequest.UserID,
		Service2faID:  service.ID,
		User2faID:     user.ID,
	}

	if userServiceLink.IsUserAlreadyBound(sc.db) {
		c.JSON(http.StatusOK, gin.H{"message": "User already bound!"})
		return
	}

	if err := sc.db.Create(&userServiceLink).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create link"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Service linked to 2FA successfully", "username": user.Username})
}

func (sc *ServiceController) CreateBindRequest(c *gin.Context) {
	// username into 2fa user id
	// serviceID
	// user id from service
}

func (sc *ServiceController) VerifyBindRequest(c *gin.Context) {
	//
}
