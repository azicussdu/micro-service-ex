package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"api-gateway/internal/config"
	"api-gateway/internal/types"

	"github.com/redis/go-redis/v9"
)

type UserCacheService struct {
	redisClient *redis.Client
	httpClient  *http.Client
	config      config.Config
}

func NewUserCacheService(redisClient *redis.Client, cfg config.Config) *UserCacheService {
	return &UserCacheService{
		redisClient: redisClient,
		// httpClient - указатель на HTTP-клиент, используется для запросов к user-сервису
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
		config: cfg,
	}
}

// Метод который возвращает пользователя по ID (из кэша или из сервиса).
func (s *UserCacheService) GetOrFetchUser(ctx context.Context, userID uint) (*types.User, error) {
	cacheKey := fmt.Sprintf("user:%d", userID) // ключ для redis
	cachedUser, err := s.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var user types.User
		if json.Unmarshal([]byte(cachedUser), &user) == nil {
			return &user, nil // если есть в redis то возвращаем
		}
	}

	// готовим внутрений запрос к user-service
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/internal/users/%d", s.config.UserServiceURL, userID), nil)
	if err != nil {
		return nil, err
	}

	// добавляем Header X-Internal-Token: чтобы не все стучались к сервису
	request.Header.Set("X-Internal-Token", s.config.InternalServiceToken)

	// Выполняет HTTP-запрос - отправляет запрос, получает ответ.
	response, err := s.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close() // Откладывает закрытие тела ответа

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("user service returned status %d", response.StatusCode)
	}

	var user types.User
	// создаёт декодер из тела ответа - далее декодирует JSON в структуру user
	if err := json.NewDecoder(response.Body).Decode(&user); err != nil {
		return nil, err
	}

	// Сериализует пользователя в JSON - превращает структуру user в байты JSON. []byte
	payload, err := json.Marshal(user)
	if err == nil {
		_ = s.redisClient.Set(ctx, cacheKey, payload, s.config.UserCacheTTL).Err()
	}

	return &user, nil
}
