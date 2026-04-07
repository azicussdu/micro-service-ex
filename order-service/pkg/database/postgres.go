package database

import (
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func MustConnect(dsn string) *gorm.DB {
	var db *gorm.DB
	var err error

	for attempt := 1; attempt <= 10; attempt++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			return db
		}

		log.Printf("database connection attempt %d failed: %v", attempt, err)
		time.Sleep(3 * time.Second)
	}

	log.Fatalf("database connection failed after retries: %v", err)
	return nil
}
