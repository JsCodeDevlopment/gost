package config

import (
	"fmt"
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

var RabbitMQConn *amqp.Connection

func ConnectRabbitMQ() {
	user := os.Getenv("RABBITMQ_USER")
	pass := os.Getenv("RABBITMQ_PASS")
	host := os.Getenv("RABBITMQ_HOST")
	port := os.Getenv("RABBITMQ_PORT")

	if user == "" {
		user = "guest"
	}
	if pass == "" {
		pass = "guest"
	}
	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "5672"
	}

	url := fmt.Sprintf("amqp://%s:%s@%s:%s/", user, pass, host, port)

	var err error
	RabbitMQConn, err = amqp.Dial(url)
	if err != nil {
		log.Printf("RabbitMQ connection failed: %v. Ensure RabbitMQ is running or ignore if not needed.\n", err)
		return
	}

	log.Println("RabbitMQ connected successfully!")
}

func CloseRabbitMQ() {
	if RabbitMQConn != nil && !RabbitMQConn.IsClosed() {
		RabbitMQConn.Close()
	}
}
