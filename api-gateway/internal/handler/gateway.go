package handler

import (
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"api-gateway/internal/config"
	"api-gateway/internal/middleware"
	"api-gateway/internal/types"

	"github.com/gin-gonic/gin"
)

// GatewayHandler - основной обработчик API Gateway.
type GatewayHandler struct {
	config     config.Config
	httpClient *http.Client
}

func NewGatewayHandler(cfg config.Config) *GatewayHandler {
	return &GatewayHandler{
		config: cfg,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func RegisterRoutes(router *gin.Engine, gatewayHandler *GatewayHandler, authMiddleware *middleware.AuthMiddleware) {
	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	api := router.Group("/api")
	{
		// gatewayHandler не обрабатывает запрос сам - а проксирует его в user-service (или в order-service)
		api.POST("/auth/register", gatewayHandler.ProxyAuth)
		api.POST("/auth/login", gatewayHandler.ProxyAuth)

		api.GET("/users/me", authMiddleware.RequireAuth(), func(c *gin.Context) {
			user := c.MustGet(middleware.ContextUserKey).(*types.User)
			c.JSON(http.StatusOK, user)
		})

		api.POST("/orders", authMiddleware.RequireAuth(), gatewayHandler.ProxyOrders)
		api.GET("/orders", authMiddleware.RequireAuth(), gatewayHandler.ProxyOrders)
		api.DELETE("/orders/:id", authMiddleware.RequireAuth(), gatewayHandler.ProxyOrders)
	}
}

func (h *GatewayHandler) ProxyAuth(c *gin.Context) {
	// ИТОГ: targetURL = http://user-service:8081/auth/login (api удалили)
	targetURL := h.config.UserServiceURL + strings.TrimPrefix(c.Request.URL.Path, "/api")
	h.forwardRequest(c, targetURL, func(request *http.Request) {})
}

func (h *GatewayHandler) ProxyOrders(c *gin.Context) {
	targetPath := "/orders"
	if c.Param("id") != "" {
		targetPath += "/" + c.Param("id")
	}

	targetURL := h.config.OrderServiceURL + targetPath
	userID := c.MustGet(middleware.ContextUserIDKey).(uint)

	h.forwardRequest(c, targetURL, func(request *http.Request) {
		// преобразует uint в строку в десятичной системе
		// Этот заголовок сообщает order-сервису, какой пользователь выполняет запрос
		request.Header.Set("X-User-ID", strconv.FormatUint(uint64(userID), 10))
		// Устанавливает внутренний токен - для аутентификации между сервисами (сервис-сервис).
		request.Header.Set("X-Internal-Token", h.config.InternalServiceToken)
		// Удаляет оригинальный заголовок Authorization - внутренние сервисы используют X-Internal-Token, не нужно передавать JWT дальше.
		request.Header.Del("Authorization")
	})
}

// Объявляет универсальный метод проксирования [mutate - функция для модификации запроса]
func (h *GatewayHandler) forwardRequest(c *gin.Context, targetURL string, mutate func(request *http.Request)) {
	// берёшь оригинальный запрос и создаёшь новый → в другой сервис
	// тот же method (POST), тот же body, тот же context (таймауты/отмена)
	request, err := http.NewRequestWithContext(c.Request.Context(), c.Request.Method, targetURL, c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "failed to build upstream request"})
		return
	}

	// переносишь ВСЕ headers (Authorization, Content-Type, и т.д.)
	request.Header = c.Request.Header.Clone()
	request.Header.Set("Content-Type", c.GetHeader("Content-Type"))
	// Вызывает функцию-мутатор - модифицирует запрос (добавляет/удаляет заголовки).
	mutate(request)

	// Выполняет HTTP-запрос к внутреннему сервису.
	response, err := h.httpClient.Do(request)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "upstream service unavailable"})
		return
	}
	defer response.Body.Close()

	// Итерируется по всем заголовкам ответа внутреннего сервиса.
	for key, values := range response.Header {
		for _, value := range values {
			// Добавляет заголовок к ответу клиенту - сохраняет все заголовки из ответа внутреннего сервиса.
			c.Writer.Header().Add(key, value)
		}
	}

	// Устанавливает HTTP-статус ответа - такой же, как вернул внутренний сервис.
	c.Status(response.StatusCode)
	// io.Copy - читает из response.Body и пишет в c.Writer (а когда ты пишешь в c.Writer то ответ уходит в браузер)
	if _, err := io.Copy(c.Writer, response.Body); err != nil {
		c.Abort()
	}
}
