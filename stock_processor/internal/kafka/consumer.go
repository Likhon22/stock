package infra

import (
	"context"

	"github.com/segmentio/kafka-go"
)


type KafkaConsumer struct {
	reader *kafka.Reader
}

func NewKafkaConsumer(broker, topic, groupID string) *KafkaConsumer {
    r := kafka.NewReader(kafka.ReaderConfig{
        Brokers: []string{broker},
        Topic:   topic,
        GroupID: groupID,   
        MaxBytes: 10e6,
    })
    return &KafkaConsumer{reader: r}  
}

func (c *KafkaConsumer) ReadMessage(ctx context.Context) (kafka.Message, error) {
    return c.reader.ReadMessage(ctx)
}

func (c *KafkaConsumer) Close() error {
    return c.reader.Close()
}