package queue

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

// MessageConsumer defines the structure to handle messages from a queue
type MessageConsumer struct {
	Handler func(string)
}

// NewMessageConsumer creates a new MessageConsumer instance
func NewMessageConsumer(handler func(string)) *MessageConsumer {
	return &MessageConsumer{
		Handler: handler,
	}
}

// Consume listens for messages from a queue and invokes the handler
func (mc *MessageConsumer) Consume(conn *amqp.Connection, queueName string) error {
	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open a channel: %w", err)
	}
	defer ch.Close()

	// Declare a queue to ensure it exists
	q, err := ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare a queue: %w", err)
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return fmt.Errorf("failed to register a consumer: %w", err)
	}

	// Run a goroutine to handle incoming messages
	go func() {
		for d := range msgs {
			mc.Handler(string(d.Body))
		}
	}()

	log.Printf("Waiting for messages on queue: %s", queueName)
	select {}
}
