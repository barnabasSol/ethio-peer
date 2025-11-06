package cache

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	client *redis.Client
	ttl    time.Duration
}

func New(addr string, ttl time.Duration) (*Redis, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   0,
	})

	ctx := context.Background()
	var err error
	maxRetries := 5
	backoff := time.Second
	for i := range maxRetries {
		err = rdb.Ping(ctx).Err()
		if err == nil {
			log.Println("bridge connected to redis")
			return &Redis{
				client: rdb,
				ttl:    ttl,
			}, nil
		}
		log.Printf(
			"failed to ping redis (attempt %d/%d): %v",
			i+1,
			maxRetries,
			err,
		)
		time.Sleep(backoff)
		backoff *= 2
	}
	return nil, err
}

var ErrCacheMiss = errors.New("cache miss")

func (c *Redis) Get(ctx context.Context, key string, dest any) error {
	val, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return ErrCacheMiss
	}
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(val), dest)
}

func (c *Redis) Set(
	ctx context.Context,
	key string,
	obj any,
) error {
	data, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, key, data, c.ttl).Err()
}

func (c *Redis) SetWithTTL(
	ctx context.Context,
	key string,
	obj any,
	ttl time.Duration,
) error {
	data, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, key, data, ttl).Err()
}

func (c *Redis) Delete(
	ctx context.Context,
	key string,
) (bool, error) {
	n, err := c.client.Del(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return n > 0, nil
}
