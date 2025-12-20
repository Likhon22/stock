package main

import (
	"context"
	"encoding/json"
	"fmt"
	"stock-processor/config"
	infra "stock-processor/internal/kafka"
	"stock-processor/internal/model"

	"github.com/segmentio/kafka-go"
)

func main() {
    fmt.Println("Stock Processor started")
    
    // Create consumer
    consumer := infra.NewKafkaConsumer(
        config.KafkaBroker,
        config.KafkaTopic,
        config.ConsumerGroup,
    )
    defer consumer.Close()  // Don't forget to close!
    
    ctx := context.Background()
    
    // Sequential loop - read one by one
    for {
        msg, err := consumer.ReadMessage(ctx)
        if err != nil {
            fmt.Printf("Error reading: %v\n", err)
            break
        }
        
        // Parse JSON
   go func(message kafka.Message) {
            var stock model.StockPrice
            if err := json.Unmarshal(message.Value, &stock); err != nil {
                fmt.Printf("Error parsing: %v\n", err)
                return
            }
            
            // Process concurrently (MULTIPLE at same time)
            fmt.Printf("Received: %s at $%.2f\n", stock.Symbol, stock.Price)
        }(msg) 
    }
}