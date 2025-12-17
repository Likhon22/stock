package config

import "time"

var (
	KafkaBroker  = "localhost:9092"
	KafkaTopic   = "stock_prices"
	Symbols      = []string{"AAPL", "GOOGL", "AMZN", "TSLA", "MSFT"}
	Interval     = 1 * time.Second
	WorkerCount  = 5
	KafkaTimeout = 5 * time.Second
	MaxRetries   = 3
)
