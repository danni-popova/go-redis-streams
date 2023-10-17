package stream

import (
	"context"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type Consumer struct {
	rdb        *redis.Client
	group      string
	consumerId string
	stream     string
	log        *zap.SugaredLogger
}

// NewConsumer will return a consumer reading from the specified stream.
// After it's created WithGroup should be called to set what consumer group
// it should be a part of.
func NewConsumer(rdb *redis.Client, consumerId, stream string) *Consumer {
	return &Consumer{
		rdb:        rdb,
		stream:     stream,
		consumerId: consumerId,
	}
}

// WithLogger attaches the a zap Logger
func (c *Consumer) WithLogger(logger *zap.SugaredLogger) *Consumer {
	c.log = logger
	return c
}

// WithGroup will attempt to create the consumer group with the specified group name
// and will set it for the consumer, so messages can be read
func (c *Consumer) WithGroup(ctx context.Context, group string) (*Consumer, error) {
	// If we specify 0, the consumer group will consume all the messages
	// in the stream history to start with.
	result := c.rdb.XGroupCreateMkStream(ctx, c.stream, group, "0")

	// If the group was created successfully, set the group name
	c.group = group
	return c, result.Err()
}

// Consume will read a message from the set group and return its value and ID.
// If an empty ID is provided, the message consumed will be a message never delivered
// to any other consumer from the group so far.
func (c *Consumer) Consume(ctx context.Context, id string) (map[string]interface{}, string, error) {
	// Special ID > is only valid in the context of consumer groups,
	// and it means: messages never delivered to other consumers so far.
	if id == "" {
		id = ">"
	}

	args := &redis.XReadGroupArgs{
		Streams:  []string{c.stream, ">"},
		Group:    c.group,
		Consumer: c.consumerId,
		Count:    1,
	}

	result := c.rdb.XReadGroup(ctx, args)

	if result.Err() != nil {
		return nil, "", result.Err()
	}

	val := result.Val()
	if len(val) > 0 && len(val[0].Messages) > 0 {
		return val[0].Messages[0].Values, val[0].Messages[0].ID, nil
	}

	return nil, "", nil
}

// Acknowledge needs to be called once the message has been successfully processed.
// It will remove the message from the Pending Entries List for the set consumer group
// so that it won't be re-read again from the stream.
func (c *Consumer) Acknowledge(ctx context.Context, msgId string) error {
	return c.rdb.XAck(ctx, c.stream, c.group, msgId).Err()
}
