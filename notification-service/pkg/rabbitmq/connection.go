package rabbitmq

import (
	"log"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	ExchangeName = "events"
	ExchangeType = "topic"
	defaultURL   = "amqp://guest:guest@rabbitmq:5672/"
)

func Connect() (*amqp.Connection, *amqp.Channel, error) {
	url := os.Getenv("RABBITMQ_URL")
	if url == "" {
		url = defaultURL
	}

	var conn *amqp.Connection
	var err error

	// 🔁 retry logic
	for i := 0; i < 10; i++ {
		conn, err = amqp.Dial(url)
		if err == nil {
			log.Println("✅ Connected to RabbitMQ")
			break
		}

		log.Printf("❌ RabbitMQ not ready (%d/10): %v\n", i+1, err)
		time.Sleep(3 * time.Second)
	}

	if err != nil {
		return nil, nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, nil, err
	}

	err = ch.ExchangeDeclare(
		ExchangeName,
		ExchangeType,
		true, // durable
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, nil, err
	}

	return conn, ch, nil
}
