package main

import (
	"context"
	"log"
	"time"

	"github.com/caarlos0/env/v9"
	"github.com/redis/go-redis/v9"

	"github.com/danni-popova/go-redis-streams/stream"
)

type Config struct {
	ConsumerGroups []string `env:"CONSUMER_GROUPS" envSeparator:","`
}

func main() {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err)
	}

	log.Println(cfg.ConsumerGroups)

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

	janitor := stream.NewJanitor(rdb, "demo")

	for {
		for _, group := range cfg.ConsumerGroups {
			janitor.GetPendingId(ctx, group)
		}

		time.Sleep(time.Second * 5)
	}
}
