package bootstrap

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

type Application struct {
	hub        *websocket.Hub
	subscriber *redis.Subscriber
	handler    *handler.Handler
	ctx        context.Context
	cancel     context.CancelFunc
}

func NewApp() *Application {
	ctx, cancel := context.WithCancel(context.Background())

	// Create Redis connection
	rdb := redis.ConnectRedis("localhost:6379", 0, ctx)

	// Create subscriber
	subscriber := redis.NewSubscriber(rdb, ctx)

	// Create WebSocket hub
	hub := websocket.NewHub()

	// Create HTTP handler
	h := handler.NewHandler(hub)

	return &Application{
		hub:        hub,
		subscriber: subscriber,
		handler:    h,
		ctx:        ctx,
		cancel:     cancel,
	}
}

func (a *Application) Run() error {
	log.Println("🚀 WebSocket Server starting...")

	defer a.subscriber.Close()
	defer a.cancel()

	// Subscribe to Redis channel
	if err := a.subscriber.Subscribe("stock_updates"); err != nil {
		return err
	}

	// Setup signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Goroutine to handle signals
	go func() {
		sig := <-sigChan
		log.Printf("\n🛑 Received signal: %v", sig)
		log.Println("🔄 Initiating graceful shutdown...")
		a.cancel()
	}()

	// Start Hub
	go a.hub.Run()

	// Start Bridge: Subscriber → Hub
	go a.startBridge()

	// Setup HTTP routes
	http.HandleFunc("/ws", a.handler.HandleWebSocket)

	// Start HTTP server
	log.Println("📡 WebSocket server listening on :8082")
	log.Println("💡 Connect via: ws://localhost:8082/ws")

	if err := http.ListenAndServe(":8082", nil); err != nil {
		return err
	}

	return nil
}

// startBridge connects Redis Pub/Sub to WebSocket Hub
func (a *Application) startBridge() {
	log.Println("🌉 Starting bridge: Subscriber → Hub")

	for {
		select {
		case <-a.ctx.Done():
			log.Println("🛑 Bridge shutting down...")
			return

		case msg, ok := <-a.subscriber.Updates():
			if !ok {
				log.Println("📭 Subscriber channel closed, bridge stopping...")
				return
			}

			a.hub.Broadcast <- msg
		}
	}
}
