package cache

import (
	"context"
	"time"

	"github.com/elys-network/elys/indexer/internal/models"
	"github.com/go-redis/redis/v8"
)

// CacheInterface defines the cache operations
type CacheInterface interface {
	// Basic operations
	Get(ctx context.Context, key string) (interface{}, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Delete(ctx context.Context, keys ...string) error

	// Subscription operations
	AddSubscription(ctx context.Context, clientID string, sub *models.WebSocketSubscription) error
	GetSubscriptions(ctx context.Context, clientID string) ([]*models.WebSocketSubscription, error)
	Subscribe(ctx context.Context, channels ...string) *redis.PubSub
}
