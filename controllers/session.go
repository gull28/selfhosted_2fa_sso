package controllers

import (
	"fmt"
	"net/http"
	"selfhosted_2fa_sso/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type SessionController struct {
	db *gorm.DB
}

type SessionRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func GetSessionController(db *gorm.DB) *SessionController {
	return &SessionController{db: db}
}

func (sc *SessionController) Index(c *gin.Context) {
	data := gin.H{
		"Message": "Welcome to the session creation page",
	}
	c.HTML(http.StatusOK, "session.html", data)
}

func (sc *SessionController) Create(c *gin.Context) {
	var req SessionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	var user models.SuperUser
	if err := sc.db.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	if !checkPasswordHash(req.Password, user.PasswordHash) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	fmt.Println("session created success !")
	c.JSON(http.StatusOK, gin.H{"message": "Session created successfully"})
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
