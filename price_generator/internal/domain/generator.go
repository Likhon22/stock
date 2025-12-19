package domain

import (
	"math/rand/v2"
	"time"
)

type PriceGenerator interface {
	Generate(symbol string, lastPrice float64) Price
}

type RandomPriceGenerator struct{}

func (g *RandomPriceGenerator) Generate(symbol string, last float64) Price {
	change := last * (rand.Float64()*0.1 - 0.05)
	newPrice := last + change
	if newPrice < 1 {
		newPrice = 1
	}
	return Price{
		Symbol:    symbol,
		Value:     newPrice,
		Timestamp: time.Now(),
	}
}
