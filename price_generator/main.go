package main

import (
	"context"
	"log"
	"math/rand"
	"price_generator/config"
	"price_generator/generator"
	"price_generator/kafka"
	"price_generator/logger"
	"time"
)

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
			last := lastPrices[symbol]
			msg := generator.GeneratePrice(symbol, last)
			lastPrices[symbol] = msg.Price
			producer.Send(ctx, msg)
			logger.Info("Sent %s price %.2f", msg.Symbol, msg.Price)
		}
	}
	logger.Info("Price Generator Service stopped")
}
