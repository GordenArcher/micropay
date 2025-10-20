
# User Service (Django)

The **User Service** handles **user registration** and publishes `user.created` events to **RabbitMQ**. This is part of an **event-driven microservices architecture** with Go-based services consuming events.

---

## Project Structure

```
user_service/
├── users/
│   ├── apps.py         # Loads signals
│   ├── models.py       # User model (Django's default)
│   ├── signals.py      # Publishes user.created events
│   └── views.py        # Functional view for registration
├── manage.py
├── requirements.txt
└── ...
```

---

## ⚙️ Prerequisites

* Python 3.13+
* Postgres database
* RabbitMQ broker

---

## Running Locally

1. Clone the repository:

```bash
git clone <repo-url>
cd user_service
```

2. Set environment variables:

```bash
export DATABASE_URL=postgres://postgres:postgres@localhost:5432/user_service
export RABBITMQ_URL=amqp://guest:guest@localhost:5672/
```

3. Install dependencies:

```bash
pip install -r requirements.txt
```

4. Run database migrations:

```bash
python manage.py migrate
```

5. Start the Django server:

```bash
python manage.py runserver 0.0.0.0:8000
```

---

## 🔗 Environment Variables

| Variable       | Description                |
| -------------- | -------------------------- |
| `DATABASE_URL` | Postgres connection string |
| `RABBITMQ_URL` | RabbitMQ connection URL    |

---

## 📝 Functionality

* **User Registration:** `POST /api/users/register/`
  Example request:

```json
{
  "username": "gorden",
  "email": "gorden@example.com",
  "password": "securepassword"
}
```

* **Event Publishing:**
  After a successful registration, a `user.created` event is published to RabbitMQ.
  Example payload:

```json
{
  "user_id": 1,
  "username": "gorden",
  "email": "gorden@example.com"
}
```

* **Signals & Transactions:**
  Django signals ensure events are published **only after database commit**.

---

## 🛠️ Development Notes

* Uses **Django signals** in `users/signals.py` to publish events.
* Uses **functional views** for simplicity; no serializers are used to filter fields—all user fields are included in events.
* `EventPublisher` class handles **resilient RabbitMQ publishing**, with retries and exponential backoff.
* All RabbitMQ messages are marked **persistent** to prevent loss.

---

## Docker (Optional)

```dockerfile
FROM python:3.13-slim

WORKDIR /app
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY . .
EXPOSE 8000

CMD ["python", "manage.py", "runserver", "0.0.0.0:8000"]
```

---

## Production Considerations

* Secure RabbitMQ with credentials and TLS.
* Use **persistent queues** and consider **dead-letter queues** for failed events.
* Add monitoring and structured logging for RabbitMQ events.
* Use Django’s **custom User model** if you plan to extend fields.
* Implement rate-limiting and input validation for registration endpoints.

---

## References

* [Django Signals Documentation](https://docs.djangoproject.com/en/stable/topics/signals/)
* [RabbitMQ Tutorials](https://www.rabbitmq.com/getstarted.html)
* [Event-driven Microservices Pattern](https://microservices.io/patterns/data/event-driven.html)
