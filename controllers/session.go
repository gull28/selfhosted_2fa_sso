package controllers

import (
	"fmt"
	"net/http"
	"selfhosted_2fa_sso/config"
	"selfhosted_2fa_sso/models"
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
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
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

	token, err := generateJWT(user.ID, sc.cfg.JWT.Secret)

	if err != nil {
		fmt.Printf("%v", err)
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
