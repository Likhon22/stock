package main

import (
	"context"
	"log"
	"math/rand/v2"
	"price_generator/config"
	"price_generator/generator"
	"price_generator/kafka"
	"price_generator/logger"
	"sync"
	"time"
)

var priceLock sync.RWMutex

func main() {
	logger.Info("Price Generator Service started")

	producer := kafka.NewProducer()
	lastPrices := make(map[string]float64)

	for _, symbol := range config.Symbols {

		lastPrices[symbol] = 100 + 50*rand.Float64()

		log.Printf("Initialized %s with starting price: %.2f", symbol, lastPrices[symbol])
	}
	ctx := context.Background()
	ticker := time.NewTicker(config.Interval)
	defer ticker.Stop()
	jobs := make(chan string, len(config.Symbols))
	for w := 0; w < 5; w++ {
		go worker(w, jobs, lastPrices, producer, ctx)
	}
	for range ticker.C {
		for _, symbol := range config.Symbols {

			jobs <- symbol
		}
	}

	logger.Info("Price Generator Service stopped")
}

func worker(id int, jobs <-chan string, lastPrices map[string]float64, producer *kafka.Producer, ctx context.Context) {
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
