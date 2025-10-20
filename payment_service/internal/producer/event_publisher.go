package producer

import (
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
)

type EventPublisher struct {
	Conn     *amqp.Connection
	Exchange string
}

// NewPublisher creates a new publisher
func NewPublisher(conn *amqp.Connection, exchange string) *EventPublisher {
	return &EventPublisher{
		Conn:     conn,
		Exchange: exchange,
	}
}

// Publish sends an event to RabbitMQ
func (p *EventPublisher) Publish(routingKey string, payload interface{}) error {
	ch, err := p.Conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	err = ch.Publish(
		p.Exchange,
		routingKey,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
		},
	)
	if err != nil {
		log.Println("Failed to publish event:", err)
		return err
	}

	log.Printf("Published event %s: %s\n", routingKey, body)
	return nil
}
