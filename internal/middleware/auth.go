package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware — JWT stub (Day 2: always passes, sets demo-user)
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		
		// TEMP: accept any Bearer token, set demo-user
		if strings.HasPrefix(auth, "Bearer ") {
			c.Set("user_id", "demo-user")
			c.Next()
			return
		}

		// No token = 401
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
		c.Abort()
	}
}

// TODO Day 3: Replace with real JWT validation:
// - Parse JWT from Authorization header
// - Verify signature (RS256/HS256)
// - Extract user_id from claims
// - Set c.Set("user_id", claims.UserID)