package bootstrap

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"stock-processor/config"
	infra "stock-processor/internal/kafka"
	"stock-processor/internal/service"

	"github.com/segmentio/kafka-go"
)

type Application struct {
    service  service.ProcessorService
    consumer *infra.KafkaConsumer
    ctx      context.Context
    cancel   context.CancelFunc
}

func NewApp() *Application {
    ctx, cancel := context.WithCancel(context.Background())
    consumer := infra.NewKafkaConsumer(
        config.KafkaBroker,
        config.KafkaTopic,
        config.ConsumerGroup,
    )
    processorService := service.NewProcessorService();
    
    return &Application{
        service:  processorService,
        consumer: consumer,
        ctx:      ctx,
        cancel:   cancel,
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