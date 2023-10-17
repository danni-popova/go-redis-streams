package main

import (
	"context"

	"github.com/caarlos0/env/v9"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"github.com/danni-popova/go-redis-streams/stream"
)

type Config struct {
	ConsumerGroup string `env:"CONSUMER_GROUP"`
	ConsumerName  string `env:"CONSUMER_NAME"`
	StreamName    string `env:"STREAM_NAME"`
}

const (
	min = 500
	max = 5000
)

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
		log.Error(err)
	}

	consumer, err := stream.NewConsumer(rdb, cfg.ConsumerName, cfg.StreamName).
		WithLogger(log).
		WithGroup(ctx, cfg.ConsumerGroup)

	if err != nil {
		log.Warnw("couldn't create consumer group", "error", err)
	}

	for {
		msgId := ""
		message, msgId, err := consumer.Consume(ctx, msgId)
		if err != nil {
			log.Warnw("couldn't consume message", "error", err)
			continue
		}

		// This can simulate the consumer doing some kind of processing
		// using the message - analyzing, persisting to DB, actioning etc.
		//waitFor := rand.Intn(max-min) + min
		//log.Infow("waiting", "duration", waitFor)
		//time.Sleep(time.Millisecond * time.Duration(waitFor))

		err = consumer.Acknowledge(ctx, msgId)
		if err != nil {
			log.Warnw("couldn't ack message", "id", msgId, "error", err)
		}
		log.Info(message)
	}
}
