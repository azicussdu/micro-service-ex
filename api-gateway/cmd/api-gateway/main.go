package main

import (
	"log"

	"api-gateway/internal/config"
	"api-gateway/internal/handler"
	"api-gateway/internal/middleware"
	"api-gateway/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func main() {
	cfg := config.MustLoad()

	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	userCacheService := service.NewUserCacheService(redisClient, cfg)
	authMiddleware := middleware.NewAuthMiddleware(cfg, userCacheService)
	gatewayHandler := handler.NewGatewayHandler(cfg)

	router := gin.Default()
	handler.RegisterRoutes(router, gatewayHandler, authMiddleware)

	log.Printf("api gateway listening on :%s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("gateway failed: %v", err)
	}
}
