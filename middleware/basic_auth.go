package middleware

import (
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"strings"
)

func ConvertBasicAuth() string {
	username := "tuanuet1"
	password := "123456789"
	credentials := username + ":" + password
	token := base64.StdEncoding.EncodeToString([]byte(credentials))
	return token
}

func BasicAuth(expectedToken string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// Kiểm tra định dạng tiêu đề Authorization
		if !strings.HasPrefix(authHeader, "Basic ") {
			c.JSON(401, gin.H{"error": "Authorization header must start with Basic"})
			c.Abort()
			return
		}

		// Kiểm tra token với token mong đợi
		if authHeader != "Basic "+expectedToken {
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		c.Next()
	}
}
