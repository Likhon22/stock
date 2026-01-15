package config

import (
	"os"
	"time"
)

var (
	KafkaBroker  = getEnv("KAFKA_BROKER", "localhost:9092")
	KafkaTopic   = "stock_prices"
	Symbols      = []string{"AAPL", "GOOGL", "AMZN", "TSLA", "MSFT"}
	Interval     = 1 * time.Second
	WorkerCount  = 5
	KafkaTimeout = 5 * time.Second
	MaxRetries   = 3
)

// getEnv reads environment variable with fallback to default value
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
