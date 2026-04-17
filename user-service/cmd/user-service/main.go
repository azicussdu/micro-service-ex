package main

import (
	"log"

	"user-service/internal/config"
	"user-service/internal/events"
	"user-service/internal/handler"
	"user-service/internal/model"
	"user-service/internal/repository"
	"user-service/internal/service"
	"user-service/pkg/database"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.MustLoad()

	db := database.MustConnect(cfg.DatabaseDSN)
	if err := db.AutoMigrate(&model.User{}); err != nil {
		log.Fatalf("failed to migrate user table: %v", err)
	}

	userRepository := repository.NewUserRepository(db)
	tokenService := service.NewTokenService(cfg.JWTSecret)
	publisher := events.NewPublisher()
	defer publisher.Close()
	authService := service.NewAuthService(userRepository, tokenService, publisher)

	router := gin.Default()
	handler.RegisterRoutes(router, cfg, authService)

	log.Printf("user service listening on :%s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("user service failed: %v", err)
	}
}
