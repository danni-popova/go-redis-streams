package main

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/danni-popova/go-redis-streams/stream"
)

const (
	StreamName = "demo"
)

func main() {
	ctx := context.Background()

	// Connect to Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// Test we can read from Redis
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatal(err)
	}

	producer := stream.NewProducer(rdb, StreamName)

	for {
		err := producer.Produce(ctx, map[string]string{"messageId": "message"})
		if err != nil {
			log.Fatal("couldn't write to stream")
		}
		time.Sleep(time.Second)
	}
}
