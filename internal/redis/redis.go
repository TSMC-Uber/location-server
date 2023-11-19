package redis

import (
	"context"
	"location-server/internal/config"

	"github.com/redis/go-redis/v9"
)

//go:generate mockgen -destination="$(git rev-parse --show-toplevel)/internal/mocks/mock_rds_client.go" -package=mocks -mock_names Client=MockRedisClient git.synology.inc/synology/synotable/internal/rds Client

type Client interface {
	// redis client
	Close() error
	Subscribe(ctx context.Context, channels ...string) *redis.PubSub
	Publish(ctx context.Context, channel string, message interface{}) *redis.IntCmd
}

type ClientImpl struct {
	*redis.Client
}

func NewClient() Client {
	return &ClientImpl{
		Client: redis.NewClient(&redis.Options{
			Addr: config.MustGetEnv("REDIS_HOST"),
		}),
	}
}
