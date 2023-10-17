package main

import (
	"context"
	"fmt"
	"time"

	"github.com/caarlos0/env/v9"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"github.com/danni-popova/go-redis-streams/stream"
)

type Config struct {
	StreamName string `env:"STREAM_NAME"`
}

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	log := logger.Sugar()

	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err)
	}

	log.Infow("config parsed", "config", cfg)

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

	producer := stream.NewProducer(rdb, cfg.StreamName).WithLogger(log)
	counter := 0

	for {
		err := producer.Produce(ctx, map[string]string{
			"message": fmt.Sprintf("message-%d", counter),
		})
		if err != nil {
			log.Errorw("couldn't write to stream", "error", err)
		}
		counter++
		time.Sleep(time.Second / 8)
	}
}
