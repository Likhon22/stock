package redis

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

var redisInstance *redis.Client

func ConnectRedis(addr string, db int, ctx context.Context) *redis.Client {
	if redisInstance != nil {
		return redisInstance

	}
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   db,
	})
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	return rdb
}
