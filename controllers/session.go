package controllers

import (
	"fmt"
	"net/http"
	"selfhosted_2fa_sso/config"
	"selfhosted_2fa_sso/models"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type SessionController struct {
	db  *gorm.DB
	cfg *config.Config
}

type SessionRequest struct {
	Username string `json:"username" binding:"required,min=3,max=20"`
	Password string `json:"password" binding:"required,min=8,max=100"`
}

func GetSessionController(db *gorm.DB, cfg *config.Config) *SessionController {
	return &SessionController{db: db, cfg: cfg}
}

func (sc *SessionController) Index(c *gin.Context) {
	data := gin.H{
		"Message": "Welcome to the session creation page",
	}
	c.HTML(http.StatusOK, "session.html", data)
}

func (sc *SessionController) Delete(c *gin.Context) {
	c.SetCookie(
		"auth_token",
		"",
		-1,
		"/",
		"",
		true,
		true,
	)

	c.Set("user", nil)

	c.JSON(http.StatusOK, gin.H{
		"message": "Logged out successfully",
	})
}

func (sc *SessionController) Create(c *gin.Context) {
	var req SessionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	var user models.SuperUser
	if err := sc.db.Where("username = ?", req.Username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	fmt.Printf("Password from request: %s\n", req.Password)
	fmt.Printf("Hashed password from DB: %s\n", user.PasswordHash)

	password := strings.TrimSpace(req.Password)

	if !checkPasswordHash(password, user.PasswordHash) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	token, err := generateJWT(user.ID, sc.cfg.JWT.Secret)
	if err != nil {
		fmt.Printf("JWT generation error: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
		return
	}

	c.SetCookie(
		"auth_token",
		token,
		3600,
		"/",
		"",
		true,
		true,
	)

	c.JSON(http.StatusOK, gin.H{"message": "Session created successfully"})
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func generateJWT(userID uint, secretKey string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 1).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}
