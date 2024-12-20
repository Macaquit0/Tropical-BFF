package cache

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type RedisStore struct {
	client *redis.Client
}

type RedisConnectionConfig struct {
	Addr     string
	Password string
	DB       int
}

func NewRedisConnection(redisUrl string, redisPassword string) (CacheStore, error) {
	ctx := context.Background()
	options := redis.Options{
		Addr: redisUrl,
	}

	if redisPassword != "" {
		options.Password = redisPassword
	}

	client := redis.NewClient(&options)

	if _, err := client.Ping(ctx).Result(); err != nil {
		log.Error().Err(err).Msg("error while connecting to redis")
		return nil, err
	}

	return &RedisStore{
		client: client,
	}, nil
}

func (r *RedisStore) StoreKey(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	err := r.client.SetNX(ctx, key, value, expiration).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *RedisStore) StoreKeyStruct(ctx context.Context, key string, value any, expiration time.Duration) error {
	marshalled, err := json.Marshal(value)
	if err != nil {
		return err
	}

	if err := r.client.SetNX(ctx, key, marshalled, expiration).Err(); err != nil {
		return err
	}

	return nil
}

func (r *RedisStore) GetKey(ctx context.Context, key string) (*string, error) {
	command, err := r.client.Get(ctx, key).Result()

	if errors.Is(err, redis.Nil) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &command, nil
}

func (r *RedisStore) ReadCheck(ctx context.Context) error {
	_, err := r.client.Ping(ctx).Result()

	if errors.Is(err, redis.Nil) {
		return nil
	}

	if err != nil {
		return err
	}

	return nil
}

func (r *RedisStore) Subscribe(ctx context.Context, key string, channel *chan string) error {
	sub := r.client.Subscribe(ctx, key)

	go func() {
		defer func() {
			err := sub.Close()
			if err != nil {
				log.Error().Err(err).Msg("error while closing redis subscription")
			}
			close(*channel)
		}()

		for {
			msg, err := sub.ReceiveMessage(ctx)
			if err != nil {
				log.Error().Err(err).Msg("error while receiving message from redis subscription")

				time.Sleep(10 * time.Second)
				continue
			}

			log.Info().Msgf("Received message from %s channel.", msg.Channel)

			log.Info().Msgf("Message: %s", msg.Payload)
			*channel <- msg.Payload
			time.Sleep(1 * time.Second)
		}
	}()

	return nil
}

func (r *RedisStore) Publish(ctx context.Context, channel string, message string) error {
	err := r.client.Publish(ctx, channel, message).Err()
	if err != nil {
		return err
	}

	return nil
}
