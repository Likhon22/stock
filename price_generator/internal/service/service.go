package service

import (
	"context"
	"encoding/json"
	"fmt"
	"price_generator/internal/domain"
	"price_generator/internal/kafka"
	"price_generator/internal/repository"
	"price_generator/logger"
	"sync"
)

type PriceService interface {
	StartWorkers(ctx context.Context, jobs <-chan string, workerCount int) *sync.WaitGroup
	ProcessSymbol(ctx context.Context, symbol string) error
}

type priceService struct {
	generator  domain.PriceGenerator
	repository repository.PriceRepository
	publisher  *kafka.KafkaProducer
}

func NewPriceService(
	gen domain.PriceGenerator,
	repo repository.PriceRepository,
	publisher *kafka.KafkaProducer,

) PriceService {
	return &priceService{
		generator:  gen,
		repository: repo,
		publisher:  publisher,
	}
}

func (s *priceService) ProcessSymbol(ctx context.Context, symbol string) error {

	lastPrice, ok := s.repository.Get(symbol)
	if !ok {
		return fmt.Errorf("symbol not found: %s", symbol)
	}
	newPrice := s.generator.Generate(symbol, lastPrice)
	s.repository.Set(symbol, newPrice.Value)
	data, err := json.Marshal(newPrice)
	if err != nil {
		return err

	}
	return s.publisher.Send(ctx, symbol, data)
}

func (s *priceService) StartWorkers(ctx context.Context, jobs <-chan string, workerCount int) *sync.WaitGroup {
	var wg sync.WaitGroup
	for w := 0; w < workerCount; w++ {
		wg.Add(1)
		go s.worker(ctx, w, jobs, &wg)
	}
	return &wg
}

func (s *priceService) worker(ctx context.Context, id int, jobs <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	logger.Info("Worker %d started", id)

	for symbol := range jobs {
		if err := s.ProcessSymbol(ctx, symbol); err != nil {
			logger.Error("Worker %d failed to process %s: %v", id, symbol, err)
		} else {
			logger.Info("Worker %d processed %s", id, symbol)
		}
	}

	logger.Info("Worker %d stopped", id)
}
