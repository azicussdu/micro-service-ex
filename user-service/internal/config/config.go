package config

import (
	"log"
	"os"
)

type Config struct {
	Port                 string
	DatabaseDSN          string
	JWTSecret            string
	InternalServiceToken string
}

func MustLoad() Config {
	return Config{
		Port:                 getEnv("PORT", "8081"),
		DatabaseDSN:          mustGetEnv("DATABASE_DSN"),
		JWTSecret:            mustGetEnv("JWT_SECRET"),
		InternalServiceToken: mustGetEnv("INTERNAL_SERVICE_TOKEN"),
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
