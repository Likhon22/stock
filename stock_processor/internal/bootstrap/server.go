package bootstrap

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"stock-processor/config"
	"stock-processor/db"

	"stock-processor/internal/handler"
	infra "stock-processor/internal/kafka"
	"stock-processor/internal/repository"
	"stock-processor/internal/routes"

	"stock-processor/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/segmentio/kafka-go"
)

type Application struct {
    service  service.ProcessorService
    consumer *infra.KafkaConsumer
    ctx      context.Context
    cancel   context.CancelFunc
    r        *chi.Mux
}

// getEnv reads environment variable with fallback to default value
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func NewApp() *Application {
    ctx, cancel := context.WithCancel(context.Background())
    consumer := infra.NewKafkaConsumer(
        config.KafkaBroker,
        config.KafkaTopic,
        config.ConsumerGroup,
    )
 
    redisAddr := getEnv("REDIS_ADDR", "localhost:6379")
    rdb :=db.ConnectRedis(redisAddr, 0)
 
    
    processorRepo:=repository.NewPriceRepository(rdb)
    processorService := service.NewProcessorService(processorRepo);
    processorHandler:=handler.NewHandler(processorService)
    processorRoutes:=routes.SetUpRoutes(processorHandler)

    return &Application{
        service:  processorService,
        consumer: consumer,
        ctx:      ctx,
        cancel:   cancel,
        r:processorRoutes,
    }
}

func (a *Application) Run() error {
    fmt.Println("Stock Processor started")
    
    defer a.consumer.Close()
    defer a.cancel()
    
    // Create jobs channel
    jobs := make(chan kafka.Message, config.WorkerCount*2)
    
    // Start worker pool
    wg := a.service.StartWorkers(a.ctx, jobs, config.WorkerCount)
    
    // Setup signal handling
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
    
    // Goroutine to handle signals
    go func() {
        sig := <-sigChan
        fmt.Printf("\nReceived signal %v, shutting down...\n", sig)
        a.cancel()
    }()
    
    fmt.Println("Workers started, consuming messages...")
 // Start HTTP server in background (non-blocking)
  go func() {
    fmt.Println("HTTP server starting on :3000")
    if err := http.ListenAndServe(":3000", a.r); err != nil {
        fmt.Printf("HTTP server error: %v\n", err)
    }
  }()
    // Main loop - BOOTSTRAP reads from Kafka!
    for {
        msg, err := a.consumer.ReadMessage(a.ctx) 
        if err != nil {
            if err == context.Canceled {
                fmt.Println("Consumer stopped by signal")
            } else {
                fmt.Printf("Error reading: %v\n", err)
            }
            break
        }
        
        jobs <- msg  
    }
    
    // Graceful shutdown
    close(jobs)
    fmt.Println("Waiting for workers to finish...")
    wg.Wait()
    fmt.Println("Shutdown complete!")
    
    return nil
}