package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"stock_websocket/internal/handler"
	"stock_websocket/internal/infra/redis"
	"stock_websocket/internal/websocket"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	setupGracefulShutdown(cancel)
	// Create hub
	hub := websocket.NewHub()
	rdb := redis.ConnectRedis("localhost:6379", 0, ctx)
	subscriber := redis.NewSubscriber(rdb, ctx)
	if err := subscriber.Subscribe("stock_updates"); err != nil {
		log.Fatalf("Failed to subscribe: %v", err)
	}
	go func() {
		log.Println("🌉 Starting bridge: Subscriber → Hub")
		for {
			select {
			case <-ctx.Done():
				// Main context cancelled (Ctrl+C or shutdown)
				log.Println("Bridge shutting down...")
				return

			case msg, ok := <-subscriber.Updates():
				if !ok {
					// Channel closed (subscriber stopped)
					log.Println("📭 Subscriber channel closed, bridge stopping...")
					return
				}
				// Forward message to all WebSocket clients
				hub.Broadcast <- msg
			}
		}
	}()

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

func setupGracefulShutdown(cancel context.CancelFunc) {
	sigChan := make(chan os.Signal, 1)

	// Notify on Interrupt (Ctrl+C) and SIGTERM (kill command)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		log.Printf("\n🛑 Received signal: %v", sig)
		log.Println("🔄 Initiating graceful shutdown...")
		cancel() // Cancel main context (stops everything)
	}()
}
