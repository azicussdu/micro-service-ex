package middleware

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"api-gateway/internal/config"
	"api-gateway/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	ContextUserIDKey = "userID"
	ContextUserKey   = "user"
)

type AuthMiddleware struct {
	config           config.Config
	userCacheService *service.UserCacheService
}

func NewAuthMiddleware(cfg config.Config, userCacheService *service.UserCacheService) *AuthMiddleware {
	return &AuthMiddleware{
		config:           cfg,
		userCacheService: userCacheService,
	}
}

func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := extractToken(c.GetHeader("Authorization"))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(m.config.JWTSecret), nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		userID, err := parseUserID(claims["user_id"])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid user claim"})
			return
		}

		user, err := m.userCacheService.GetOrFetchUser(c.Request.Context(), userID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
			return
		}

		c.Set(ContextUserIDKey, userID)
		c.Set(ContextUserKey, user)
		c.Next()
	}
}

func extractToken(header string) (string, error) {
	if header == "" {
		return "", errors.New("missing authorization header")
	}

	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", errors.New("invalid authorization header")
	}

	return parts[1], nil
}

func parseUserID(value interface{}) (uint, error) {
	switch v := value.(type) {
	case float64:
		return uint(v), nil
	case string:
		parsed, err := strconv.ParseUint(v, 10, 64)
		return uint(parsed), err
	default:
		return 0, errors.New("unsupported user_id claim")
	}
}
