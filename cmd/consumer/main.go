package main

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"

	"github.com/danni-popova/go-redis-streams/stream"
)

const (
	ConsumerName = "consumer-1"
	StreamName   = "demo"
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

	consumer, err := stream.NewConsumer(rdb, ConsumerName, StreamName).WithGroup(ctx, "group")
	if err != nil {
		log.Fatal("couldn't create consumer group")
	}

	for {
		msgId := ""
		message, msgId, err := consumer.Consume(ctx, msgId)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(message)
	}
}
