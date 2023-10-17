package main

import (
	"context"
	"log"

	"github.com/caarlos0/env/v9"
	"github.com/redis/go-redis/v9"

	"github.com/danni-popova/go-redis-streams/stream"
)

type Config struct {
	ConsumerGroup string `env:"CONSUMER_GROUPS"`
	ConsumerName  string `env:"CONSUMER_NAME"`
	StreamName    string `env:"STREAM_NAME"`
}

func main() {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err)
	}

	log.Println(cfg)

	ctx := context.Background()

	// Connect to Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis-stack:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// Test we can read from Redis
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatal(err)
	}

	consumer, err := stream.NewConsumer(rdb, cfg.ConsumerName, cfg.StreamName).
		WithGroup(ctx, cfg.ConsumerGroup)
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
