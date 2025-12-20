package config

import "time"



var (
    KafkaBroker   = "localhost:9092"
    KafkaTopic    = "stock_prices"
    ConsumerGroup = "stock-processor-group"  
		Symbols      = []string{"AAPL", "GOOGL", "AMZN", "TSLA", "MSFT"}
	 Interval     = 1 * time.Second
	 WorkerCount  = 5
)