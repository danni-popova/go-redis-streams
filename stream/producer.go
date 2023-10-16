package stream

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type Producer struct {
	rdb    *redis.Client
	stream string
}

// NewProducer creates a producer that writes to the specified Redis stream
func NewProducer(rdb *redis.Client, stream string) *Producer {
	return &Producer{
		rdb:    rdb,
		stream: stream,
	}
}

func (p *Producer) Produce(ctx context.Context, message map[string]string) error {
	return p.rdb.XAdd(ctx, &redis.XAddArgs{
		Stream: p.stream,
		Values: message,
	}).Err()
}
