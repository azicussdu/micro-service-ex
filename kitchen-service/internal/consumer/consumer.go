package consumer

import (
	"log"

	"kitchen-service/pkg/rabbitmq"
)

func Run() error {
	conn, ch, err := rabbitmq.Connect()
	if err != nil {
		return err
	}
	defer conn.Close()
	defer ch.Close()

	queue, err := ch.QueueDeclare("kitchen-queue", true, false, false, false, nil)
	if err != nil {
		return err
	}

	if err := ch.QueueBind(queue.Name, "order.*", rabbitmq.ExchangeName, false, nil); err != nil {
		return err
	}

	messages, err := ch.Consume(queue.Name, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	log.Println("Kitchen consumer started, waiting for messages...")
	for msg := range messages {
		log.Printf("Kitchen received order event: %s\n", string(msg.Body))

		msg.Ack(false)
	}

	return nil
}
