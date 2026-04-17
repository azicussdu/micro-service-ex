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

	// 🔁 retry until RabbitMQ is ready
	for {
		conn, err = amqp.Dial(url)
		if err == nil {
			log.Println("✅ Connected to RabbitMQ")
			break
		}

		log.Println("⏳ Waiting for RabbitMQ...", err)
		time.Sleep(3 * time.Second)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, nil, err
	}

	err = ch.ExchangeDeclare(
		ExchangeName,
		ExchangeType,
		true,  // durable
		false, // auto-delete
		false, // internal
		false, // no-wait
		nil,
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, nil, err
	}

	return conn, ch, nil
}
