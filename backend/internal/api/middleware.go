package api

import (
	"backend/internal/service"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Read JWT from cookie instead of Authorization header
		tokenStr, err := c.Cookie("auth_token")
		if err != nil || tokenStr == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing auth token"})
			c.Abort()
			return
		}

		// Validate token
		userID, err := service.ParseToken(tokenStr, []byte(os.Getenv("JWT_SECRET")))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		// Store userID in context so handlers can use it
		c.Set("userID", uint(userID))
		c.Next()
	}
}
