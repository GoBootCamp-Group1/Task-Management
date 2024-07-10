package ports

import (
	"context"
	"time"
)

type CacheRepository interface {
	// Set stores the value in the cache
	Set(ctx context.Context, key string, value string, ttl time.Duration) error
	// Get retrieves the value from the cache
	Get(ctx context.Context, key string) (string, error)
	// Delete removes the value from the cache
	Delete(ctx context.Context, key string) error
	// DeleteByPrefix removes the value from the cache with the given prefix
	DeleteByPrefix(ctx context.Context, prefix string) error
	// Close closes the connection to the cache server
	Close() error
}
