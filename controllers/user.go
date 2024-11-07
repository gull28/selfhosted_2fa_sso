package controllers

import (
	"net/http"
	"selfhosted_2fa_sso/models"

	"github.com/gin-gonic/gin"
	"github.com/pquerna/otp/totp"
	"gorm.io/gorm"
)

type UserController struct {
	db *gorm.DB
}

type VerifyRequest struct {
	Username string `json:"username" binding:"required"`
	Code     string `json:"code" binding:"required"`
}

type CreateRequest struct {
	Username string `json:"username" binding:"required"`
}

func GetUserController(db *gorm.DB) *UserController {
	return &UserController{db: db}
}

func (uc *UserController) CreateUser(c *gin.Context) {
	var user models.User2fa
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	var existingUser models.User2fa
	if err := uc.db.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Username or email already exists"})
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

func (uc *UserController) VerifyUser(c *gin.Context) {
	var req VerifyRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	var user models.User2fa
	if err := uc.db.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if totp.Validate(req.Code, user.TOTPSecret) {
		c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Code is valid"})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "failure", "message": "Invalid code"})
	}
}
