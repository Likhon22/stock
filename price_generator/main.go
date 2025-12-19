package main

import (
	"context"
	"os"
	"os/signal"
	"price_generator/config"
	"price_generator/generator"
	"price_generator/kafka"
	"price_generator/logger"
	"price_generator/utils"
	"sync"
	"syscall"
	"time"
)

var priceLock sync.RWMutex

func main() {
	logger.Info("Price Generator Service started")

	producer := kafka.NewProducer()
	lastPrices := utils.InitializePrices(config.Symbols)
	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()
	ticker := time.NewTicker(config.Interval)
	defer ticker.Stop()

	jobs := make(chan string, len(config.Symbols)*2)
	var wg sync.WaitGroup
	for w := 0; w < 5; w++ {
		wg.Add(1)
		go worker(w, jobs, lastPrices, producer, ctx, &wg)
	}
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	running := true
	for running {
		select {
		case <-ticker.C:
			for _, symbol := range config.Symbols {

				jobs <- symbol
			}
		case sig := <-sigChan:
			logger.Info("Received signal: %v, shutting down gracefully...", sig)
			running = false
		}
	}

	ticker.Stop()
	close(jobs)

	logger.Info("Waiting for workers to finish...")
	wg.Wait()

	// Clean up Kafka connection
	logger.Info("Closing Kafka producer...")
	if err := producer.Close(); err != nil {
		logger.Info("Error closing producer: %v", err)
	}

	logger.Info("Shutdown complete!")
}

func worker(id int, jobs <-chan string, lastPrices map[string]float64, producer *kafka.Producer, ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	logger.Info("Worker %d started", id)
	for symbol := range jobs {
		priceLock.RLock()
		last := lastPrices[symbol]
		priceLock.RUnlock()

		msg := generator.GeneratePrice(symbol, last)

		priceLock.Lock()
		lastPrices[symbol] = msg.Price
		priceLock.Unlock()

		producer.Send(ctx, msg)
		logger.Info("Worker %d processed %s", id, symbol)
	}
}
