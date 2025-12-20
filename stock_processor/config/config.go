package config



var (
    KafkaBroker   = "localhost:9092"
    KafkaTopic    = "stock_prices"
    ConsumerGroup = "stock-processor-group"  // NEW! Consumers need group ID
)