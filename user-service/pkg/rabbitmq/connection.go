package rabbitmq

import (
	"os"

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

	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, nil, err
	}

	if err := ch.ExchangeDeclare(ExchangeName, ExchangeType, true, false, false, false, nil); err != nil {
		ch.Close()
		conn.Close()
		return nil, nil, err
	}

	return conn, ch, nil
}
