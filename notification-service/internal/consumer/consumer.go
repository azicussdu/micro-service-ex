package consumer

import (
	"log"

	"notification-service/pkg/rabbitmq"
)

func Run() error {
	conn, ch, err := rabbitmq.Connect()
	if err != nil {
		return err
	}
	defer conn.Close()
	defer ch.Close()

	queue, err := ch.QueueDeclare("notification-queue", true, false, false, false, nil)
	if err != nil {
		return err
	}

	for _, routingKey := range []string{"order.*", "user.*"} {
		if err := ch.QueueBind(queue.Name, routingKey, rabbitmq.ExchangeName, false, nil); err != nil {
			return err
		}
	}

	messages, err := ch.Consume(queue.Name, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	log.Println("Notification consumer started, waiting for messages...")
	for msg := range messages {
		log.Printf("Notification received event: %s", string(msg.Body))

		msg.Ack(false)
	}
	log.Println("Subscribed to queue:", queue.Name)

	return nil
}
