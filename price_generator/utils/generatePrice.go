package utils

import (
	"math/rand/v2"
	"price_generator/logger"
)

func InitializePrices(symbols []string) map[string]float64 {
	prices := make(map[string]float64)
	for _, symbol := range symbols {
		prices[symbol] = 100 + 50*rand.Float64()
		logger.Info("Initialized %s with starting price: %.2f", symbol, prices[symbol])
	}
	return prices
}
