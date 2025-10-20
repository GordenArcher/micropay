

# Event-Driven Microservices: User Service & Payment Service

This project demonstrates a **minimal event-driven microservices architecture** using **Django** and **Go**, integrated via **RabbitMQ**.

* **User Service** (Django): Handles user registration and publishes `user.created` events.
* **Payment Service** (Go): Listens to `user.created` events and creates payment accounts.
* **Postgres**: Stores payment data.
* **RabbitMQ**: Event broker for asynchronous communication.

---

## Project Structure

```
event-driven-payments/
│
├── user_service/                # Django user service
│   ├── users/
│   │   ├── models.py
│   │   ├── views.py
│   │   ├── signals.py          # publishes events to RabbitMQ
│   │   └── apps.py
│   ├── manage.py
│   ├── requirements.txt
│   └── ...
│
├── payment_service/             # Go payment service
│   ├── main.go
│   ├── go.mod
│   ├── internal/
│   │   ├── consumer/           # listens to RabbitMQ events
│   │   ├── producer/           # optional publisher
│   │   ├── service/            # payment logic
│   │   └── rabbitmq/           # RabbitMQ connection helper
    ├── configs/
    ├── README.md
│   └── Dockerfile
│
├── docker-compose.yml           # orchestrates Postgres, RabbitMQ, and services
└── README.md
```

---

## Prerequisites

* Docker & Docker Compose
* Go 1.21+ (for payment service, if building locally)
* Python 3.13+ (for user service, if running locally)

---

## Running Locally with Docker Compose

1. Clone the repo:

```bash
git clone <repo-url>
cd event-driven-payments
```

2. Build and start all services:

```bash
docker-compose up --build
```

3. Access services:

* **User Service (Django)**: `http://localhost:8000`
* **RabbitMQ Management UI**: `http://localhost:15672` (guest/guest)
* **Postgres**: `localhost:5432` (user: `postgres`, password: `postgres`, db: `payments`)

---

## User Registration Flow

1. User registers via **User Service API**:

```
POST /api/users/register/
{
  "username": "gorden",
  "email": "gorden@example.com",
  "password": "password123"
}
```

2. **User Service** triggers `post_save` signal for the user.
3. Signal publishes `user.created` event to **RabbitMQ**.
4. **Payment Service** consumes the event and creates a payment account in Postgres.

---

## Environment Variables

### User Service

| Variable       | Description                |
| -------------- | -------------------------- |
| `RABBITMQ_URL` | RabbitMQ connection URL    |
| `DATABASE_URL` | Postgres connection string |

### Payment Service

| Variable       | Description                |
| -------------- | -------------------------- |
| `RABBITMQ_URL` | RabbitMQ connection URL    |
| `DATABASE_URL` | Postgres connection string |

---

## Dockerfiles

### User Service

```dockerfile
FROM python:3.13-slim

WORKDIR /app
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt
COPY . .
CMD ["python", "manage.py", "runserver", "0.0.0.0:8000"]
```

### Payment Service (multi-stage, production-ready)

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

## Development Notes

* **RabbitMQ**: Default exchange `events` is topic-based.
* **Postgres**: Database `payments` is required for Payment Service.
* **EventPublisher** (Django) ensures messages are sent **after transaction commit**.
* **Payment Service** listens continuously for `user.created` events.

---

## Production Considerations

* Use **healthchecks** for Postgres and RabbitMQ in Docker Compose.
* Add **retry logic** and **dead-letter queues** for RabbitMQ events.
* Secure services with **TLS** and **authentication**.
* Use **multi-stage Dockerfiles** to minimize image size.
* Implement **migrations** for database schemas.

---

## References

* [Django Signals](https://docs.djangoproject.com/en/stable/topics/signals/)
* [RabbitMQ Tutorials](https://www.rabbitmq.com/getstarted.html)
* [Go `streadway/amqp` library](https://github.com/streadway/amqp)


