package handler

import (
	"net/http"
	"strconv"

	"order-service/internal/config"
	"order-service/internal/middleware"
	"order-service/internal/service"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderService *service.OrderService
}

type createOrderRequest struct {
	ProductName string `json:"product_name" binding:"required"`
}

func RegisterRoutes(router *gin.Engine, cfg config.Config, orderService *service.OrderService) {
	handler := &OrderHandler{orderService: orderService}

	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	orders := router.Group("/orders", middleware.RequireGateway(cfg.InternalServiceToken))
	{
		orders.POST("", handler.Create)
		orders.GET("", handler.List)
		orders.DELETE("/:id", handler.Delete)
	}
}

func (h *OrderHandler) Create(c *gin.Context) {
	var request createOrderRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := h.orderService.Create(userIDFromContext(c), request.ProductName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, order)
}

func (h *OrderHandler) List(c *gin.Context) {
	orders, err := h.orderService.List(userIDFromContext(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"orders": orders})
}

func (h *OrderHandler) Delete(c *gin.Context) {
	orderID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id"})
		return
	}

	if err := h.orderService.Delete(uint(orderID), userIDFromContext(c)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "order deleted"})
}

func userIDFromContext(c *gin.Context) uint {
	return c.MustGet(middleware.ContextUserIDKey).(uint)
}
