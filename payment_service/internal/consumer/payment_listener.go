package consumer

import (
	"encoding/json"
	"log"

	"payment_service/internal/service"

	"github.com/streadway/amqp"
)

type UserCreatedEvent struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

func StartPaymentListener(conn *amqp.Connection) error {
	ch, err := conn.Channel()
	if err != nil {
		return err
	}

	if err := ch.ExchangeDeclare("events", "topic", true, false, false, false, nil); err != nil {
		return err
	}

	q, err := ch.QueueDeclare(
		"user-created-queue",
		true,  // durable
		false, // auto delete
		false, // exclusive
		false,
		nil,
	)
	if err != nil {
		return err
	}

	if err := ch.QueueBind(q.Name, "user.created", "events", false, nil); err != nil {
		return err
	}

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	// Init DB connection
	db, err := service.InitDB("postgres://postgres:postgres@localhost:5432/payments?sslmode=disable")
	if err != nil {
		return err
	}

	go func() {
		for d := range msgs {
			var event UserCreatedEvent
			if err := json.Unmarshal(d.Body, &event); err != nil {
				log.Println("Failed to parse event:", err)
				continue
			}
			log.Printf("Received user.created event: %+v", event)
			if err := service.CreatePaymentAccount(event.UserID, event.Username, event.Email, db); err != nil {
				log.Println("Failed to create payment account:", err)
			}
		}
	}()

	return nil
}
