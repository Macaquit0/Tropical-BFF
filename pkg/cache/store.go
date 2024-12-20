package cache

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
)

type CacheStore interface {
	StoreKey(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	StoreKeyStruct(ctx context.Context, key string, value any, expiration time.Duration) error
	GetKey(ctx context.Context, key string) (*string, error)
	ReadCheck(ctx context.Context) error
	Subscribe(ctx context.Context, key string, channel *chan string) error
	Publish(ctx context.Context, channel string, message string) error
}

func NewCacheStore(redisUrl string, redisPassword string) (CacheStore, error) {
	cache, err := NewRedisConnection(redisUrl, redisPassword)
	if err != nil {
		log.Error().Err(err).Msg("failed to establish redis connection")
		return nil, err
	}

	return cache, nil
}
