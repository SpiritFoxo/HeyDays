package config

import (
	"fmt"
	"heydays/handlers"
	"os"
)

func SetupRabbitMQManager() *handlers.RabbitMQManager {
	amqpURI := fmt.Sprintf(
		"amqp://%s:%s@%s:%s/",
		os.Getenv("RABBITMQ_USER"),
		os.Getenv("RABBITMQ_PASS"),
		os.Getenv("RABBITMQ_HOST"),
		os.Getenv("RABBITMQ_PORT"),
	)

	return handlers.NewRabbitMQManager(amqpURI)
}
