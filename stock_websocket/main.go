package main

import (
	"log"
	"net/http"
	"stock_websocket/internal/handler"
	"stock_websocket/internal/websocket"
)

func main() {
	// Create hub
	hub := websocket.NewHub()

	// Start hub (runs forever)
	go hub.Run()

	// Create handler with hub
	h := handler.NewHandler(hub)

	// Setup route
	http.HandleFunc("/ws", h.HandleWebSocket)

	log.Println("WebSocket server starting on :8082")
	if err := http.ListenAndServe(":8082", nil); err != nil {
		log.Fatal(err)
	}
}
