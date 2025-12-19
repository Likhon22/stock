package bootstrap

import (
	"context"
	"os"
	"os/signal"
	"price_generator/config"
	"price_generator/internal/domain"
	"price_generator/internal/kafka"
	"price_generator/internal/repository"
	"price_generator/internal/service"
	"price_generator/logger"
	"syscall"
	"time"
)

type Application struct {
	service     service.PriceService
	kafkaClient *kafka.KafkaProducer
	ctx         context.Context
	cancel      context.CancelFunc
}

func NewApp() *Application {
	ctx, cancel := context.WithCancel(context.Background())
	kafkaProducer := kafka.NewKafkaProducer(config.KafkaBroker, config.KafkaTopic)
	priceGen := &domain.RandomPriceGenerator{}
	priceRepo := repository.NewPriceRepository(config.Symbols)
	priceService := service.NewPriceService(priceGen, priceRepo, kafkaProducer)
	return &Application{
		service:     priceService,
		kafkaClient: kafkaProducer,
		ctx:         ctx,
		cancel:      cancel,
	}
}

func (a *Application) Run() error {
	logger.Info("Price Generator Service started")

	defer a.kafkaClient.Close()
	defer a.cancel()

	// Start worker pool
	jobs := make(chan string, len(config.Symbols)*2)
	wg := a.service.StartWorkers(a.ctx, jobs, config.WorkerCount)

	// Setup signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Ticker loop
	ticker := time.NewTicker(config.Interval)
	defer ticker.Stop()

	logger.Info("Workers started, beginning price generation")

	// Main loop
	running := true
	for running {
		select {
		case <-ticker.C:
			for _, symbol := range config.Symbols {
				jobs <- symbol
			}
		case sig := <-sigChan:
			logger.Info("Received signal: %v, shutting down...", sig)
			running = false
		}
	}

	// Graceful shutdown
	a.cancel()
	close(jobs)
	logger.Info("Waiting for workers to finish...")
	wg.Wait()
	logger.Info("Shutdown complete!")

	return nil
}
