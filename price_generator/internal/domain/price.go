package domain

import "time"

type Price struct {
	Symbol    string    `json:"symbol"`
	Value     float64   `json:"price"`
	Timestamp time.Time `json:"timestamp"`
}
