package main

import (
	"log"

	"order-service/internal/config"
	"order-service/internal/events"
	"order-service/internal/handler"
	"order-service/internal/model"
	"order-service/internal/repository"
	"order-service/internal/service"
	"order-service/pkg/database"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.MustLoad()

	db := database.MustConnect(cfg.DatabaseDSN)
	if err := db.AutoMigrate(&model.Order{}); err != nil {
		log.Fatalf("failed to migrate order table: %v", err)
	}

	orderRepository := repository.NewOrderRepository(db)
	publisher := events.NewPublisher()
	defer publisher.Close()

	orderService := service.NewOrderService(orderRepository, publisher)

	router := gin.Default()
	handler.RegisterRoutes(router, cfg, orderService)

	log.Printf("order service listening on :%s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("order service failed: %v", err)
	}
}
