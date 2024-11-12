package middleware

import (
	"net/http"
	"selfhosted_2fa_sso/internal/ratelimit"

	"github.com/gin-gonic/gin"
)

func RateLimiterMiddleware(rl *ratelimit.RateLimiter) gin.HandlerFunc {

	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := rl.GetLimiter(ip)

		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
			c.Abort()
			return
		}

		c.Next()
	}
}
