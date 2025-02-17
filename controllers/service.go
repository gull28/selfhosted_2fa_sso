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
