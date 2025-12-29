package redis

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

type Subscriber struct {
	client  *redis.Client
	pubsub  *redis.PubSub
	updates chan []byte
	ctx     context.Context
	cancel  context.CancelFunc
}

func NewSubscriber(client *redis.Client, parentCtx context.Context) *Subscriber {
	ctx, cancel := context.WithCancel(parentCtx)
	updates := make(chan []byte, 100)

	return &Subscriber{
		client:  client,
		updates: updates,
		ctx:     ctx,
		cancel:  cancel,
	}

}

func (s *Subscriber) Subscribe(channel string) error {
	s.pubsub = s.client.Subscribe(s.ctx, channel)

	_, err := s.pubsub.Receive(s.ctx)
	if err != nil {
		return fmt.Errorf("failed to subscribe: %w", err)
	}
	go s.listen()
	log.Printf("✅ Subscribed to Redis channel: %s", channel)
	return nil
}

func (s *Subscriber) listen() {
	ch := s.pubsub.Channel()
	for {

		select {
		case <-s.ctx.Done():
			log.Println(" Stopping Redis subscriber...")
			close(s.updates)
			return
		case msg := <-ch:
			select {
			case s.updates <- []byte(msg.Payload):
				log.Printf("📨 Forwarded message: %s", msg.Payload[:50])

			default:
				log.Println("⚠️ Updates channel full, dropping message")
			}
		}

	}

}

func (s *Subscriber) Updates() <-chan []byte {
	return s.updates

}

func (s *Subscriber) Close() error {
	log.Println("closing the subscriber")
	s.cancel()
	if s.pubsub != nil {
		if err := s.pubsub.Close(); err != nil {
			return fmt.Errorf("failed to close pubsub: %w", err)
		}
	}
	log.Println("✅ Subscriber closed")
	return nil

}
