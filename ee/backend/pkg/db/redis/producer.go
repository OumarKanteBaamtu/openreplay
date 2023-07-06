package redis

import (
	"github.com/go-redis/redis"
	redis2 "openreplay/backend/pkg/db/redis"
	"openreplay/backend/pkg/queue/types"
)

type producerImpl struct {
	client *redis2.Client
}

func (c *producerImpl) Close(timeout int) {
	//TODO implement me
	panic("implement me")
}

func NewProducer(client *redis2.Client) types.Producer {
	return &producerImpl{
		client: client,
	}
}

func (c *producerImpl) Produce(topic string, key uint64, value []byte) error {
	args := &redis.XAddArgs{
		Stream: topic,
		Values: map[string]interface{}{
			"sessionID": key,
			"value":     value,
		},
		MaxLenApprox: c.client.Cfg.MaxLength,
	}
	_, err := c.client.Redis.XAdd(args).Result()
	return err
}

func (c *producerImpl) ProduceToPartition(topic string, partition, key uint64, value []byte) error {
	return c.Produce(topic, key, value)
}

func (c *producerImpl) Flush(timeout int) {}
