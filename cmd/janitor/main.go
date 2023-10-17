package main

import (
	"context"
	"time"

	"github.com/caarlos0/env/v9"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"github.com/danni-popova/go-redis-streams/stream"
)

type Config struct {
	ConsumerGroups []string `env:"CONSUMER_GROUPS" envSeparator:","`
	StreamName     string   `env:"STREAM_NAME"`
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
		log.Fatalw("couldn't ping redis", "error", err)
	}

	janitor := stream.NewJanitor(rdb, cfg.StreamName)

	for {
		for _, group := range cfg.ConsumerGroups {
			result, err := janitor.GetPendingId(ctx, group)
			if err != nil {
				log.Errorw("XPENDING returned an error",
					"group", group,
					"error", err)
				continue
			}
			log.Infow("XPENDING returned info",
				"XPENDING", result,
				"group", group)
		}

		time.Sleep(time.Second)
	}
}
