package middleware

import (
	"fmt"
	"net/http"
	"selfhosted_2fa_sso/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

func AuthMiddleware(db *gorm.DB, jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		isRequestingLoginScreen := c.Request.URL.Path == "/session/create"

		tokenString, err := c.Cookie("auth_token")

		if err != nil {
			handleUnauthorized(c, isRequestingLoginScreen)
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			handleUnauthorized(c, isRequestingLoginScreen)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			handleUnauthorized(c, isRequestingLoginScreen)
			return
		}

		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			handleUnauthorized(c, isRequestingLoginScreen)
			return
		}

		userID := uint(userIDFloat)
		user, err := models.FindSuperUserById(db, userID)
		if err != nil {
			handleUnauthorized(c, isRequestingLoginScreen)
			return
		}

		if isRequestingLoginScreen {
			c.Redirect(http.StatusFound, "/services")
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Next()
	}
}

func handleUnauthorized(c *gin.Context, isRequestingLoginScreen bool) {
	if !isRequestingLoginScreen {
		c.Redirect(http.StatusFound, "/session/create")
		c.Abort()
	} else {
		c.Next()
	}
}
