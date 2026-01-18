
# Payment Service (Go)

The **Payment Service** listens for `user.created` events from **RabbitMQ** and creates payment accounts in **Postgres**. It demonstrates a **minimal, event-driven microservice** using Go.

---

## Project Structure

```
payment_service/
├── main.go
├── go.mod
├── configs/
│   ├── config.yaml/  
├── internal/
│   ├── consumer/      # Listens to RabbitMQ events
│   │   └── payment_listener.go
│   ├── producer/      # Optional event publisher
│   │   └── event_publisher.go
│   ├── service/       # Payment logic
│   │   └── payment.go
│   └── rabbitmq/      # RabbitMQ connection helper
│       └── connect.go
└── Dockerfile
```

---

## Prerequisites

* Go 1.21+
* Postgres database (with database `payments`)
* RabbitMQ broker

---

## Running Locally

1. Clone the repository:

```bash
git clone https://github.com/GordenArcher/micropay/tree/main/payment_service
cd payment_service
```

2. Set environment variables:

```bash
export DATABASE_URL=postgres://postgres:postgres@localhost:5432/payments
export RABBITMQ_URL=amqp://guest:guest@localhost:5672/
```

3. Run the service:

```bash
go run main.go
```

The service will connect to RabbitMQ, listen for `user.created` events, and create payment accounts.

---

## Environment Variables

| Variable       | Description                |
| -------------- | -------------------------- |
| `DATABASE_URL` | Postgres connection string |
| `RABBITMQ_URL` | RabbitMQ connection URL    |

---

## Functionality

* **RabbitMQ Connection:** Connects to a durable exchange (`events`) and subscribes to a queue.
* **Event Consumption:** Listens for `user.created` events.
* **Payment Creation:** Logs or saves a payment account for each user received.
* **Logging:** Prints successful account creation to stdout.

Example log:

```
2025/10/20 13:23:43 Created payment account for user 1 (gorden, gorden@example.com)
```

---

## Development Notes

* Uses `streadway/amqp` for RabbitMQ integration.
* Payment logic is in `internal/service/payment.go`.
* RabbitMQ connection is in `internal/rabbitmq/connect.go`.
* Idempotency and database integration should be implemented for production.
* Dockerfile is multi-stage for smaller, production-ready images.

---

## Docker

```dockerfile
# Build stage
FROM golang:1.21-bullseye AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o payment_service main.go

# Final stage
FROM debian:bullseye-slim
WORKDIR /app
COPY --from=builder /app/payment_service ./
EXPOSE 8080
CMD ["./payment_service"]
```

---

## Production Considerations

* Ensure Postgres database `payments` exists before running.
* Use **persistent queues** and **dead-letter queues** in RabbitMQ.
* Add proper error handling and retries for RabbitMQ messages.
* Secure RabbitMQ connections with TLS and credentials.
* Use structured logging for monitoring and debugging.

---

## References

* [Go `streadway/amqp` library](https://github.com/streadway/amqp)
* [RabbitMQ Tutorials](https://www.rabbitmq.com/getstarted.html)
* [Event-driven Microservices Pattern](https://microservices.io/patterns/data/event-driven.html)

