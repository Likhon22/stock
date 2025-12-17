package generator

import (
	"math/rand"
	"time"
)

type StockMessage struct {
	Symbol    string  `json:"symbol"`
	Price     float64 `json:"price"`
	Timestamp int64   `json:"timestamp"`
	Version   int     `json:"version"`
}

func GeneratePrice(symbol string, lastPrice float64) StockMessage {

	change := lastPrice * (rand.Float64()*0.1 - 0.05)
	newPrice := lastPrice + change
	if newPrice < 1 {
		newPrice = 1

	}
	return StockMessage{
		Symbol:    symbol,
		Price:     newPrice,
		Timestamp: time.Now().Unix(),
		Version:   1,
	}
}
