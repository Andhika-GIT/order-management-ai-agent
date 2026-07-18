package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type Publisher struct {
	client *redis.Client
}

func NewPublisher(client *redis.Client) *Publisher {
	return &Publisher{
		client: client,
	}
}

func (p *Publisher) PublishMessage(ctx context.Context, channel string, message string) error {
	return p.client.Publish(ctx, channel, message).Err()
}
