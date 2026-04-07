package middleware

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

const ContextUserIDKey = "userID"

func RequireGateway(expectedToken string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("X-Internal-Token") != expectedToken {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid internal token"})
			return
		}

		userID, err := strconv.ParseUint(c.GetHeader("X-User-ID"), 10, 64)
		if err != nil || userID == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing user id"})
			return
		}

		c.Set(ContextUserIDKey, uint(userID))
		c.Next()
	}
}
