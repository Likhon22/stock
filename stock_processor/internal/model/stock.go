package model

import "time"

type StockPrice struct {
	Symbol    string    `json:"symbol"`
	Price     float64   `json:"price"`
	Timestamp time.Time `json:"timestamp"`
}
type StockHistoryItem struct {
	Timestamp string
	Price     float64
}
