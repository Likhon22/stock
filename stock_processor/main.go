package main

import (
	"context"
	"encoding/json"
	"fmt"
	"stock-processor/config"
	infra "stock-processor/internal/kafka"
	"stock-processor/internal/model"
	"sync"

	"github.com/segmentio/kafka-go"
)

func main() {
    fmt.Println("Stock Processor started")
    var wg sync.WaitGroup
    // Create consumer
    consumer := infra.NewKafkaConsumer(
        config.KafkaBroker,
        config.KafkaTopic,
        config.ConsumerGroup,
    )
    defer consumer.Close()  // Don't forget to close!
    
ctx, cancel := context.WithCancel(context.Background())
defer cancel()
    jobs:=make(chan kafka.Message, len(config.Symbols)*2)
  
		 for i := 0; i < config.WorkerCount; i++ {
			wg.Add(1)
			go worker(jobs,&wg)
		 }

		 for{

			msg,err:=consumer.ReadMessage(ctx)
			  if err != nil {
        break
    }
	   	jobs<-msg
		 }
   
  close(jobs)		 
	wg.Wait()	
}


func worker(jobs <-chan kafka.Message,wg *sync.WaitGroup)  {
	defer wg.Done()
	for msg := range jobs {
		var stock model.StockPrice
            if err := json.Unmarshal(msg.Value, &stock); err != nil {
                fmt.Printf("Error parsing: %v\n", err)
               continue 
            }
            
						   fmt.Printf("Received: %s at $%.2f\n", stock.Symbol, stock.Price)
	}
	
}