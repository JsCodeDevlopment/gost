package messaging

import (
	"fmt"
	"log"

	"gost/src/config"

	amqp "github.com/rabbitmq/amqp091-go"
)

type ConsumerHandler func(d amqp.Delivery)

func RegisterConsumer(queueName string, handler ConsumerHandler) error {
	if config.RabbitMQConn == nil || config.RabbitMQConn.IsClosed() {
		return fmt.Errorf("rabbitmq connection is not available")
	}

	ch, err := config.RabbitMQConn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open a channel: %v", err)
	}

	q, err := ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare a queue: %v", err)
	}

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		return fmt.Errorf("failed to set QoS: %v", err)
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack (we manual ack in handler or here)
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return fmt.Errorf("failed to register a consumer: %v", err)
	}

	go func() {
		for d := range msgs {
			handler(d)
		}
		log.Printf("Consumer for %s has stopped", queueName)
		ch.Close()
	}()

	log.Printf(" [*] Waiting for messages on %s. To exit press CTRL+C", queueName)
	return nil
}

func SimpleAckHandler(handler func([]byte) error) ConsumerHandler {
	return func(d amqp.Delivery) {
		err := handler(d.Body)
		if err != nil {
			log.Printf("Error processing message: %v", err)
			d.Nack(false, false)
			return
		}
		d.Ack(false)
	}
}
