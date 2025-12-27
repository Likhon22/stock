package websocket

import "time"

// Client → Server messages
type SubscribeMessage struct {
    Action  string   `json:"action"`  // "subscribe" or "unsubscribe"
    Symbols []string `json:"symbols"` // ["AAPL", "GOOGL"]
}

// Server → Client messages
type PriceUpdate struct {
    Type      string    `json:"type"`      // "price_update"
    Symbol    string    `json:"symbol"`
    Price     float64   `json:"price"`
    Timestamp time.Time `json:"timestamp"`
}

type ErrorMessage struct {
    Type  string `json:"type"`   // "error"
    Error string `json:"error"`
}