package service

import (
	"context"
	"encoding/json"
	"fmt"
	"stock-processor/internal/model"
	"stock-processor/internal/repository"

	"sync"

	"github.com/segmentio/kafka-go"
)

type ProcessorService interface {
	StartWorkers(ctx context.Context, jobs <-chan kafka.Message, workerCount int) *sync.WaitGroup
	ProcessMessage(ctx context.Context, msg kafka.Message) error
	GetHistory(ctx context.Context, symbol string, limit int) ([]float64, error)
	GetCache(ctx context.Context, symbol string) (float64, error)
	GetAll(ctx context.Context) (map[string]float64, error)
}

type processorService struct {
	repo repository.Repository
}

// Constructor - create new service
func NewProcessorService(repo repository.Repository) ProcessorService {
	return &processorService{
		repo: repo,
	}
}

func (s *processorService) StartWorkers(ctx context.Context, jobs <-chan kafka.Message, workerCount int) *sync.WaitGroup {

	var wg sync.WaitGroup
	for w := 0; w < workerCount; w++ {
		wg.Add(1)
		go s.worker(ctx, w, jobs, &wg)
	}
	return &wg
}

func (s *processorService) ProcessMessage(ctx context.Context, msg kafka.Message) error {
	var stock model.StockPrice

	if err := json.Unmarshal(msg.Value, &stock); err != nil {
		return fmt.Errorf("failed to parse message: %w", err)
	}

	if err := s.repo.SetCache(ctx, stock); err != nil {
		fmt.Printf("Cache error for %s: %v\n", stock.Symbol, err)
	}
	if err := s.repo.AddToHistory(ctx, stock); err != nil {
		fmt.Printf("history error for %s: %v\n", stock.Symbol, err)
	}
	if err := s.repo.PublishUpdate(ctx, stock); err != nil {
		fmt.Printf("📡 Pub/Sub error for %s: %v\n", stock.Symbol, err)
	}
	fmt.Println("published to redis")

	fmt.Printf("Received: %s at $%.2f at %s\n",
		stock.Symbol,
		stock.Price,
		stock.Timestamp.Format("15:04:05"))

	return nil

}

func (s *processorService) worker(ctx context.Context, id int, jobs <-chan kafka.Message, wg *sync.WaitGroup) {

	defer wg.Done()
	fmt.Printf("Worker %d started\n", id)
	for msg := range jobs {
		if err := s.ProcessMessage(ctx, msg); err != nil {
			fmt.Printf("Worker %d error: %v\n", id, err)
		}
	}
	fmt.Printf("Worker %d stopped\n", id)
}

func (s *processorService) GetCache(ctx context.Context, symbol string) (float64, error) {
	return s.repo.GetCache(ctx, symbol)
}
func (s *processorService) GetAll(ctx context.Context) (map[string]float64, error) {
	return s.repo.GetAll(ctx)
}

func (s *processorService) GetHistory(ctx context.Context, symbol string, limit int) ([]float64, error) {
	return s.repo.GetHistory(ctx, symbol, limit)
}
