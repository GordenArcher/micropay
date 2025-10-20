package rabbitmq

import (
	"log"
	"time"

	"github.com/streadway/amqp"
)

// NewConnection connects to RabbitMQ and ensures the exchange exists
func NewConnection(url, exchange, exchangeType string) (*amqp.Connection, error) {
	var conn *amqp.Connection
	var err error

	// Retry with exponential backoff
	for attempt := 1; attempt <= 5; attempt++ {
		conn, err = amqp.Dial(url)
		if err == nil {
			break
		}
		sleep := time.Duration(2<<attempt) * time.Second
		log.Printf("RabbitMQ connection attempt %d failed: %v. Retrying in %s", attempt, err, sleep)
		time.Sleep(sleep)
	}
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	// Declare exchange
	if err := ch.ExchangeDeclare(
		exchange,
		exchangeType,
		true,  // durable
		false, // auto-delete
		false, // internal
		false, // no-wait
		nil,
	); err != nil {
		ch.Close()
		conn.Close()
		return nil, err
	}

	ch.Close() // temporary channel used only for setup
	log.Println("Connected to RabbitMQ and exchange declared:", exchange)
	return conn, nil
}

// GetChannel returns a channel for consuming or publishing
func GetChannel(conn *amqp.Connection) (*amqp.Channel, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	return ch, nil
}
