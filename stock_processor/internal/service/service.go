package service

import (
	"context"
	"encoding/json"
	"fmt"
	"stock-processor/internal/model"

	"sync"

	"github.com/segmentio/kafka-go"
)

type ProcessorService interface {
    StartWorkers(ctx context.Context, jobs <-chan kafka.Message, workerCount int) *sync.WaitGroup
    ProcessMessage(ctx context.Context, msg kafka.Message) error
}


type processorService struct {

}

// Constructor - create new service
func NewProcessorService() ProcessorService {
    return &processorService{}
}

func (s *processorService) StartWorkers(ctx context.Context, jobs <-chan kafka.Message, workerCount int) *sync.WaitGroup  {
	
	var wg sync.WaitGroup
  for w := 0; w < workerCount; w++ {
        wg.Add(1)
        go s.worker(ctx, w, jobs, &wg)
    }
	return &wg	
}


func (s *processorService) ProcessMessage(ctx context.Context, msg kafka.Message) error  {
	  var stock model.StockPrice
    
    if err := json.Unmarshal(msg.Value, &stock); err != nil {
        return fmt.Errorf("failed to parse message: %w", err)
    }
    
    fmt.Printf("Received: %s at $%.2f at %s\n", 
        stock.Symbol, 
        stock.Price, 
        stock.Timestamp.Format("15:04:05"))
    
    return nil

}

func (s *processorService) worker(ctx context.Context, id int, jobs <-chan kafka.Message, wg *sync.WaitGroup) {

    defer wg.Done()
		fmt.Printf("Worker %d started\n", id)
  for msg:= range jobs {
		 if err := s.ProcessMessage(ctx, msg); err != nil {
            fmt.Printf("Worker %d error: %v\n", id, err)
    }
	}
    fmt.Printf("Worker %d stopped\n", id)
}