package main

import (
	"log"

	"io/ioutil"

	"gopkg.in/yaml.v3"

	"payment_service/internal/consumer"
	"payment_service/internal/rabbitmq"
)

type Config struct {
	RabbitMQ struct {
		URL          string `yaml:"url"`
		Exchange     string `yaml:"exchange"`
		ExchangeType string `yaml:"exchange_type"`
		Queue        string `yaml:"queue"`
		RoutingKey   string `yaml:"routing_key"`
	} `yaml:"rabbitmq"`
}

func main() {
	// Load config
	cfg := Config{}
	f, err := ioutil.ReadFile("configs/config.yaml")
	if err != nil {
		log.Fatal("Failed to read config:", err)
	}
	if err := yaml.Unmarshal(f, &cfg); err != nil {
		log.Fatal("Failed to parse config:", err)
	}

	// Connect RabbitMQ
	conn, err := rabbitmq.NewConnection(cfg.RabbitMQ.URL, cfg.RabbitMQ.Exchange, cfg.RabbitMQ.ExchangeType)
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ:", err)
	}
	defer conn.Close()

	// Start consumer
	if err := consumer.StartPaymentListener(conn); err != nil {
		log.Fatal("Failed to start listener:", err)
	}

	// Block forever
	select {}
}
