package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"gost/src/config"

	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishMessage(exchange, routingKey string, payload interface{}) error {
	if config.RabbitMQConn == nil || config.RabbitMQConn.IsClosed() {
		return fmt.Errorf("rabbitmq connection is not available")
	}

	ch, err := config.RabbitMQConn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open a channel: %v", err)
	}
	defer ch.Close()

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = ch.PublishWithContext(ctx,
		exchange,   // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent, // Persistent delivery
		})

	if err != nil {
		return fmt.Errorf("failed to publish a message: %v", err)
	}

	log.Printf(" [x] Sent to %s: %s", routingKey, body)
	return nil
}
