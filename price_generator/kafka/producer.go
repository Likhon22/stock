package kafka

import (
	"context"
	"encoding/json"
	"log"
	"price_generator/config"
	"price_generator/generator"
	"time"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	Writer *kafka.Writer
}

func NewProducer() *Producer {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(config.KafkaBroker),
		Topic:    config.KafkaTopic,
		Balancer: &kafka.LeastBytes{},
	}
	return &Producer{Writer: writer}
}

func (p *Producer) Send(ctx context.Context, msg generator.StockMessage) {
	data, err := json.Marshal(msg)
	if err != nil {
		log.Println("[ERROR] Failed to marshal message:", err)
		return
	}

	var lastErr error
	for i := 0; i < config.MaxRetries; i++ {
		ctxWithTimeout, cancel := context.WithTimeout(ctx, config.KafkaTimeout)
		err := p.Writer.WriteMessages(ctxWithTimeout, kafka.Message{
			Key:   []byte(msg.Symbol),
			Value: data,
		})
		cancel()
		if err == nil {
			return
		}
		lastErr = err
		log.Printf("[WARN] Kafka send failed for %s (attempt %d): %v", msg.Symbol, i+1, err)
		time.Sleep(100 * time.Millisecond)
	}
	log.Printf("[ERROR] Failed to send message after %d attempts: %v", config.MaxRetries, lastErr)
}
