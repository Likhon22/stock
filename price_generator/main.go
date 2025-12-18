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

var wg sync.WaitGroup
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
	for range ticker.C {
		log.Println("--New tick--")
		for _, symbol := range config.Symbols {
			wg.Add(1)
			go func(sym string) {
				defer wg.Done()
				priceLock.Lock()
				last := lastPrices[sym]
				priceLock.Unlock()
				msg := generator.GeneratePrice(sym, last)
				priceLock.Lock()
				lastPrices[sym] = msg.Price
				priceLock.Unlock()
				producer.Send(ctx, msg)
				logger.Info("Sent %s price %.2f", msg.Symbol, msg.Price)
			}(symbol)
		}
		wg.Wait()
		log.Println("Tick complete")
	}

	logger.Info("Price Generator Service stopped")
}
