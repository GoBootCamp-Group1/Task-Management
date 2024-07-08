package cache

import (
	"context"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/port"
	"github.com/redis/go-redis/v9"
	"sync"
	"time"
)

var (
	instance *CacheRepository
	once     sync.Once
)

type CacheRepository struct {
	client *redis.Client
}

func NewCacheRepository(client *redis.Client) port.CacheRepository {
	once.Do(func() {
		instance = &CacheRepository{
			client: client,
		}
	})
	return instance
}

// Set stores the value in the redis database
func (c *CacheRepository) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	return c.client.Set(ctx, key, value, ttl).Err()
}

// Get retrieves the value from the redis database
func (c *CacheRepository) Get(ctx context.Context, key string) (string, error) {
	return c.client.Get(ctx, key).Result()
}

// Delete removes the value from the redis database
func (c *CacheRepository) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

// DeleteByPrefix removes the value from the redis database with the given prefix
func (c *CacheRepository) DeleteByPrefix(ctx context.Context, prefix string) error {
	var cursor uint64
	var keys []string

	for {
		var err error
		keys, cursor, err = c.client.Scan(ctx, cursor, prefix, 100).Result()
		if err != nil {
			return err
		}

		for _, key := range keys {
			err := c.client.Del(ctx, key).Err()
			if err != nil {
				return err
			}
		}

		if cursor == 0 {
			break
		}
	}

	return nil
}

// Close closes the connection to the redis database
func (c *CacheRepository) Close() error {
	return c.client.Close()
}
