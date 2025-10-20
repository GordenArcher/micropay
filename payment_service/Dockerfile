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
