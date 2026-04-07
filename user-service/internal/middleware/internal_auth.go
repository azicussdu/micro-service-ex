package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RequireInternalToken(expectedToken string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("X-Internal-Token") != expectedToken {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid internal token"})
			return
		}

		c.Next()
	}
}
