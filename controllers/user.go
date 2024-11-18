package controllers

import (
	"net/http"
	"selfhosted_2fa_sso/config"
	"selfhosted_2fa_sso/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pquerna/otp/totp"
	"gorm.io/gorm"
)

type UserController struct {
	db  *gorm.DB
	cfg *config.Config
}

type VerifyRequest struct {
	ServiceUserID string `json:"serviceUserId" binding:"required"`
	Code          string `json:"code" binding:"required"`
}

type CreateRequest struct {
	Username string `json:"username" binding:"required"`
}

type UpdateRequest struct {
	Username    string `json:"username" binding:"required"`
	OldUsername string `json:"oldUsername" binding:"required"`
	Code        string `json:"code" binding:"required"`
}

type CheckValidRequest struct {
	ServiceID uint   `json:"serviceId" binding:"required"`
	UserID    string `json:"userId" binding:"required"`
}

func GetUserController(db *gorm.DB, cfg *config.Config) *UserController {
	return &UserController{db: db, cfg: cfg}
}

func (uc *UserController) Create(c *gin.Context) {
	var user models.User2fa
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	var existingUser models.User2fa
	if err := uc.db.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
		return
	}

	totpSecret, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "selfhosted_2fa_sso",
		AccountName: req.Username,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate TOTP secret"})
		return
	}
	user.TOTPSecret = totpSecret.Secret()
	user.Username = req.Username

	if err := user.CreateUser(uc.db); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (uc *UserController) Update(c *gin.Context) {
	var req UpdateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid data"})
		return
	}

	var existingUser models.User2fa
	if err := uc.db.Where("username = ?", req.Username).First(&existingUser); err == nil {
		c.JSON(http.StatusConflict, gin.H{"message": "Username already exists"})
		return
	}

	var user models.User2fa
	if err := uc.db.Where("username = ?", req.OldUsername).First(&user); err != nil {
		c.JSON(http.StatusConflict, gin.H{"message": "Couldn't find current user"})
		return
	}

	totpSecret, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "selfhosted_2fa_sso",
		AccountName: req.Username,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Couldn't generate TOTP"})
		return
	}

	user.Username = req.Username
	user.TOTPSecret = totpSecret.Secret()

	if err := user.UpdateUser(uc.db); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (uc *UserController) Verify(c *gin.Context) {
	var req VerifyRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	var userService models.UserServiceLink

	if err := uc.db.Where("service_user_id = ?", req.ServiceUserID).First(&userService).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	userService.ValidUntil = time.Now().Add(time.Duration(uc.cfg.Auth.ValidFor) * time.Minute)

	if err := uc.db.Save(&userService).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user service link"})
		return
	}

	if totp.Validate(req.Code, userService.User2fa.TOTPSecret) {

		c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Code is valid"})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "failure", "message": "Invalid code"})
	}
}

func (uc *UserController) CheckSession(c *gin.Context) {
	var request CheckValidRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	isValid, err := models.IsAuthValid(uc.db, request.UserID, request.ServiceID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": "Failed to fetch user-service link"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"valid": isValid})
}
