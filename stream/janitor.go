package stream

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

type Janitor struct {
	rdb    *redis.Client
	stream string
}

func NewJanitor(rdb *redis.Client, stream string) *Janitor {
	return &Janitor{
		rdb:    rdb,
		stream: stream,
	}
}

func (j *Janitor) GetPendingId(ctx context.Context, group string) (*redis.XPending, error) {
	result := j.rdb.XPending(ctx, j.stream, group)
	return result.Val(), result.Err()
}

func (j *Janitor) Cleanup(ctx context.Context, minId string) {
	result := j.rdb.XTrimMinIDApprox(ctx, j.stream, minId, 1)

	if result.Err() != nil {
		return
	}
	log.Println(result.Val())
}
