package repository

import (
	"math/rand/v2"
	"sync"
)

type PriceRepository interface {
	Get(symbol string) (float64, bool)
	Set(symbol string, price float64)
}
type priceRepository struct {
	mu     sync.RWMutex
	prices map[string]float64
}

func NewPriceRepository(symbols []string) PriceRepository {

	repo := &priceRepository{
		prices: make(map[string]float64),
	}
	for _, symbol := range symbols {
		repo.prices[symbol] = 100 + 50*rand.Float64()
	}
	return repo
}

func (r *priceRepository) Get(symbol string) (float64, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	price, ok := r.prices[symbol]
	return price, ok

}
func (r *priceRepository) Set(symbol string, price float64) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	r.prices[symbol] = price
}
