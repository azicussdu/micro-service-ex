package events

import (
	"encoding/json"
	"sync"

	"user-service/pkg/rabbitmq"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	mu   sync.Mutex
	conn *amqp.Connection
	ch   *amqp.Channel
}

type Message struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}

func NewPublisher() *Publisher {
	return &Publisher{}
}

func (p *Publisher) Publish(routingKey string, data interface{}) error {
	body, err := json.Marshal(Message{
		Event: routingKey,
		Data:  data,
	})
	if err != nil {
		return err
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	if err := p.ensureChannel(); err != nil {
		return err
	}

	if err := p.ch.Publish(
		rabbitmq.ExchangeName,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	); err != nil {
		p.close()
		return err
	}

	return nil
}

func (p *Publisher) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.close()
}

func (p *Publisher) ensureChannel() error {
	if p.conn != nil && !p.conn.IsClosed() && p.ch != nil && !p.ch.IsClosed() {
		return nil
	}

	conn, ch, err := rabbitmq.Connect()
	if err != nil {
		return err
	}

	p.conn = conn
	p.ch = ch
	return nil
}

func (p *Publisher) close() {
	if p.ch != nil && !p.ch.IsClosed() {
		_ = p.ch.Close()
	}
	if p.conn != nil && !p.conn.IsClosed() {
		_ = p.conn.Close()
	}
	p.ch = nil
	p.conn = nil
}
