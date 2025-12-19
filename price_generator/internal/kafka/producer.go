// internal/infrastructure/kafka/producer.go
package kafka

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

// Generic - doesn't know about domain!
type MessageProducer interface {
	Send(ctx context.Context, key string, value []byte) error
	Close() error
}

type KafkaProducer struct {
	writer *kafka.Writer
}

func NewKafkaProducer(broker, topic string) *KafkaProducer {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(broker),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
	return &KafkaProducer{writer: writer}
}

func (p *KafkaProducer) Send(ctx context.Context, key string, value []byte) error {
	var lastErr error

	for i := 0; i < 3; i++ {
		ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
		err := p.writer.WriteMessages(ctxWithTimeout, kafka.Message{
			Key:   []byte(key),
			Value: value,
		})
		cancel()

		if err == nil {
			return nil
		}

		lastErr = err
		log.Printf("[WARN] Kafka send failed (attempt %d): %v", i+1, err)
		time.Sleep(100 * time.Millisecond)
	}

	return fmt.Errorf("failed after 3 attempts: %w", lastErr)
}
func (p *KafkaProducer) Close() error {
	return p.writer.Close()
}
