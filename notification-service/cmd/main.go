package main

import (
	"log"

	"notification-service/internal/consumer"
)

func main() {
	if err := consumer.Run(); err != nil {
		log.Fatalf("notification service failed: %v", err)
	}
}
