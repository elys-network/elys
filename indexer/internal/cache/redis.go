package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/elys-network/elys/indexer/internal/config"
	"github.com/elys-network/elys/indexer/internal/models"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const (
	// Key prefixes
	prefixOrderBook    = "orderbook:"
	prefixSpotOrder    = "spot_order:"
	prefixPerpOrder    = "perp_order:"
	prefixPosition     = "position:"
	prefixMarketData   = "market:"
	prefixSubscription = "sub:"

	// TTL values
	ttlOrderBook    = 5 * time.Minute
	ttlOrder        = 1 * time.Hour
	ttlPosition     = 1 * time.Hour
	ttlMarketData   = 1 * time.Minute
	ttlSubscription = 24 * time.Hour
)

type Cache struct {
	client *redis.Client
	logger *zap.Logger
}

func New(cfg *config.RedisConfig, logger *zap.Logger) (*Cache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		MaxRetries:   cfg.MaxRetries,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
		PoolTimeout:  cfg.PoolTimeout,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &Cache{
		client: client,
		logger: logger,
	}, nil
}

func (c *Cache) Close() error {
	return c.client.Close()
}

// Order Book Cache
func (c *Cache) SetOrderBook(ctx context.Context, assetPair string, orderBook *models.OrderBookSnapshot) error {
	key := prefixOrderBook + assetPair

	data, err := json.Marshal(orderBook)
	if err != nil {
		return fmt.Errorf("failed to marshal order book: %w", err)
	}

	return c.client.Set(ctx, key, data, ttlOrderBook).Err()
}

func (c *Cache) GetOrderBook(ctx context.Context, assetPair string) (*models.OrderBookSnapshot, error) {
	key := prefixOrderBook + assetPair

	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var orderBook models.OrderBookSnapshot
	if err := json.Unmarshal(data, &orderBook); err != nil {
		return nil, fmt.Errorf("failed to unmarshal order book: %w", err)
	}

	return &orderBook, nil
}

// Spot Order Cache
func (c *Cache) SetSpotOrder(ctx context.Context, order *models.SpotOrder) error {
	key := fmt.Sprintf("%s%d", prefixSpotOrder, order.OrderID)

	data, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("failed to marshal spot order: %w", err)
	}

	return c.client.Set(ctx, key, data, ttlOrder).Err()
}

func (c *Cache) GetSpotOrder(ctx context.Context, orderID uint64) (*models.SpotOrder, error) {
	key := fmt.Sprintf("%s%d", prefixSpotOrder, orderID)

	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var order models.SpotOrder
	if err := json.Unmarshal(data, &order); err != nil {
		return nil, fmt.Errorf("failed to unmarshal spot order: %w", err)
	}

	return &order, nil
}

// Perpetual Order Cache
func (c *Cache) SetPerpetualOrder(ctx context.Context, order *models.PerpetualOrder) error {
	key := fmt.Sprintf("%s%d", prefixPerpOrder, order.OrderID)

	data, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("failed to marshal perpetual order: %w", err)
	}

	return c.client.Set(ctx, key, data, ttlOrder).Err()
}

func (c *Cache) GetPerpetualOrder(ctx context.Context, orderID uint64) (*models.PerpetualOrder, error) {
	key := fmt.Sprintf("%s%d", prefixPerpOrder, orderID)

	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var order models.PerpetualOrder
	if err := json.Unmarshal(data, &order); err != nil {
		return nil, fmt.Errorf("failed to unmarshal perpetual order: %w", err)
	}

	return &order, nil
}

// Position Cache
func (c *Cache) SetPosition(ctx context.Context, position *models.PerpetualPosition) error {
	key := fmt.Sprintf("%s%d", prefixPosition, position.MtpID)

	data, err := json.Marshal(position)
	if err != nil {
		return fmt.Errorf("failed to marshal position: %w", err)
	}

	return c.client.Set(ctx, key, data, ttlPosition).Err()
}

func (c *Cache) GetPosition(ctx context.Context, mtpID uint64) (*models.PerpetualPosition, error) {
	key := fmt.Sprintf("%s%d", prefixPosition, mtpID)

	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var position models.PerpetualPosition
	if err := json.Unmarshal(data, &position); err != nil {
		return nil, fmt.Errorf("failed to unmarshal position: %w", err)
	}

	return &position, nil
}

// Pub/Sub for WebSocket Updates
func (c *Cache) PublishOrderBookUpdate(ctx context.Context, assetPair string, update *models.WSOrderBookUpdate) error {
	channel := "orderbook:" + assetPair

	data, err := json.Marshal(update)
	if err != nil {
		return fmt.Errorf("failed to marshal order book update: %w", err)
	}

	return c.client.Publish(ctx, channel, data).Err()
}

func (c *Cache) PublishTradeUpdate(ctx context.Context, asset string, update *models.WSTradeUpdate) error {
	channel := "trades:" + asset

	data, err := json.Marshal(update)
	if err != nil {
		return fmt.Errorf("failed to marshal trade update: %w", err)
	}

	return c.client.Publish(ctx, channel, data).Err()
}

func (c *Cache) PublishOrderUpdate(ctx context.Context, ownerAddress string, update *models.WSOrderUpdate) error {
	channel := "orders:" + ownerAddress

	data, err := json.Marshal(update)
	if err != nil {
		return fmt.Errorf("failed to marshal order update: %w", err)
	}

	return c.client.Publish(ctx, channel, data).Err()
}

func (c *Cache) PublishPositionUpdate(ctx context.Context, ownerAddress string, update *models.WSPositionUpdate) error {
	channel := "positions:" + ownerAddress

	data, err := json.Marshal(update)
	if err != nil {
		return fmt.Errorf("failed to marshal position update: %w", err)
	}

	return c.client.Publish(ctx, channel, data).Err()
}

// Subscribe to channels
func (c *Cache) Subscribe(ctx context.Context, channels ...string) *redis.PubSub {
	return c.client.Subscribe(ctx, channels...)
}

// Market Data Cache
func (c *Cache) SetMarketData(ctx context.Context, asset string, data interface{}) error {
	key := prefixMarketData + asset

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal market data: %w", err)
	}

	return c.client.Set(ctx, key, jsonData, ttlMarketData).Err()
}

func (c *Cache) GetMarketData(ctx context.Context, asset string) ([]byte, error) {
	key := prefixMarketData + asset

	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	return data, nil
}

// Subscription Management
func (c *Cache) AddSubscription(ctx context.Context, clientID string, subscription *models.WebSocketSubscription) error {
	key := prefixSubscription + clientID

	data, err := json.Marshal(subscription)
	if err != nil {
		return fmt.Errorf("failed to marshal subscription: %w", err)
	}

	return c.client.Set(ctx, key, data, ttlSubscription).Err()
}

func (c *Cache) GetSubscriptions(ctx context.Context, clientID string) ([]*models.WebSocketSubscription, error) {
	pattern := prefixSubscription + clientID + ":*"

	keys, err := c.client.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, err
	}

	var subscriptions []*models.WebSocketSubscription
	for _, key := range keys {
		data, err := c.client.Get(ctx, key).Bytes()
		if err != nil {
			c.logger.Error("failed to get subscription", zap.String("key", key), zap.Error(err))
			continue
		}

		var sub models.WebSocketSubscription
		if err := json.Unmarshal(data, &sub); err != nil {
			c.logger.Error("failed to unmarshal subscription", zap.String("key", key), zap.Error(err))
			continue
		}

		subscriptions = append(subscriptions, &sub)
	}

	return subscriptions, nil
}

func (c *Cache) RemoveSubscription(ctx context.Context, clientID, subscriptionID string) error {
	key := fmt.Sprintf("%s%s:%s", prefixSubscription, clientID, subscriptionID)
	return c.client.Del(ctx, key).Err()
}

// Batch operations for performance
func (c *Cache) BatchSetOrders(ctx context.Context, orders interface{}) error {
	pipe := c.client.Pipeline()

	switch v := orders.(type) {
	case []*models.SpotOrder:
		for _, order := range v {
			key := fmt.Sprintf("%s%d", prefixSpotOrder, order.OrderID)
			data, err := json.Marshal(order)
			if err != nil {
				return err
			}
			pipe.Set(ctx, key, data, ttlOrder)
		}
	case []*models.PerpetualOrder:
		for _, order := range v {
			key := fmt.Sprintf("%s%d", prefixPerpOrder, order.OrderID)
			data, err := json.Marshal(order)
			if err != nil {
				return err
			}
			pipe.Set(ctx, key, data, ttlOrder)
		}
	default:
		return fmt.Errorf("unsupported order type for batch operation")
	}

	_, err := pipe.Exec(ctx)
	return err
}

// Basic cache operations to implement CacheInterface

// Get retrieves a value from the cache by key
func (c *Cache) Get(ctx context.Context, key string) (interface{}, error) {
	result, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	}
	return result, err
}

// Set stores a value in the cache with the specified expiration
func (c *Cache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}
	return c.client.Set(ctx, key, data, expiration).Err()
}

// Delete removes one or more keys from the cache
func (c *Cache) Delete(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}
	return c.client.Del(ctx, keys...).Err()
}
