package config

import (
	"log"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port                 string
	UserServiceURL       string
	OrderServiceURL      string
	RedisAddr            string
	RedisPassword        string
	RedisDB              int
	JWTSecret            string
	InternalServiceToken string
	UserCacheTTL         time.Duration
}

func MustLoad() Config {
	redisDB, err := strconv.Atoi(getEnv("REDIS_DB", "0"))
	if err != nil {
		log.Fatalf("invalid REDIS_DB: %v", err)
	}

	userCacheTTL, err := time.ParseDuration(getEnv("USER_CACHE_TTL", "10m"))
	if err != nil {
		log.Fatalf("invalid USER_CACHE_TTL: %v", err)
	}

	return Config{
		Port:                 getEnv("PORT", "8080"),
		UserServiceURL:       mustGetEnv("USER_SERVICE_URL"),
		OrderServiceURL:      mustGetEnv("ORDER_SERVICE_URL"),
		RedisAddr:            mustGetEnv("REDIS_ADDR"),
		RedisPassword:        getEnv("REDIS_PASSWORD", ""),
		RedisDB:              redisDB,
		JWTSecret:            mustGetEnv("JWT_SECRET"),
		InternalServiceToken: mustGetEnv("INTERNAL_SERVICE_TOKEN"),
		UserCacheTTL:         userCacheTTL,
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}

func mustGetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("missing required env var %s", key)
	}

	return value
}
