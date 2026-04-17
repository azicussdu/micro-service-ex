package main

import (
	"log"

	"kitchen-service/internal/consumer"
)

func main() {
	if err := consumer.Run(); err != nil {
		log.Fatalf("kitchen service failed: %v", err)
	}
}
