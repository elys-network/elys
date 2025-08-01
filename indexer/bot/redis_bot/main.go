package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/shopspring/decimal"
)

type RedisBot struct {
	client           *redis.Client
	mnemonic         string
	marketID         uint64
	matchingInterval time.Duration
}

type Order struct {
	OrderID         uint64          `json:"order_id"`
	MarketID        uint64          `json:"market_id"`
	Side            string          `json:"side"` // "buy" or "sell"
	Price           decimal.Decimal `json:"price"`
	RemainingAmount decimal.Decimal `json:"remaining_amount"`
	Owner           string          `json:"owner"`
	Timestamp       time.Time       `json:"timestamp"`
}

type Match struct {
	BuyOrderID  uint64          `json:"buy_order_id"`
	SellOrderID uint64          `json:"sell_order_id"`
	Price       decimal.Decimal `json:"price"`
	Amount      decimal.Decimal `json:"amount"`
	Timestamp   time.Time       `json:"timestamp"`
}

func NewRedisBot(redisURL string, mnemonic string, marketID uint64) *RedisBot {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Fatalf("Failed to parse Redis URL: %v", err)
	}

	client := redis.NewClient(opt)

	// Test connection
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	return &RedisBot{
		client:           client,
		mnemonic:         mnemonic,
		marketID:         marketID,
		matchingInterval: 500 * time.Millisecond,
	}
}

func (b *RedisBot) Start(ctx context.Context) {
	log.Printf("Starting Redis bot for market %d", b.marketID)
	log.Printf("Bot address derived from mnemonic: %s...", b.mnemonic[:20])

	// Seed some test data
	b.seedTestData(ctx)

	ticker := time.NewTicker(b.matchingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Redis bot stopped")
			return
		case <-ticker.C:
			if err := b.matchOrders(ctx); err != nil {
				log.Printf("Error matching orders: %v", err)
			}
		}
	}
}

func (b *RedisBot) seedTestData(ctx context.Context) {
	// Create some test orders
	testOrders := []Order{
		{OrderID: 1, MarketID: b.marketID, Side: "buy", Price: decimal.NewFromInt(100), RemainingAmount: decimal.NewFromInt(10), Owner: "test1", Timestamp: time.Now()},
		{OrderID: 2, MarketID: b.marketID, Side: "buy", Price: decimal.NewFromInt(99), RemainingAmount: decimal.NewFromInt(5), Owner: "test2", Timestamp: time.Now()},
		{OrderID: 3, MarketID: b.marketID, Side: "sell", Price: decimal.NewFromInt(101), RemainingAmount: decimal.NewFromInt(8), Owner: "test3", Timestamp: time.Now()},
		{OrderID: 4, MarketID: b.marketID, Side: "sell", Price: decimal.NewFromInt(102), RemainingAmount: decimal.NewFromInt(12), Owner: "test4", Timestamp: time.Now()},
	}

	for _, order := range testOrders {
		// Store order details
		orderKey := fmt.Sprintf("order:%d", order.OrderID)
		orderData, _ := json.Marshal(order)
		b.client.Set(ctx, orderKey, orderData, 0)

		// Add to sorted set based on side
		if order.Side == "buy" {
			buyBookKey := fmt.Sprintf("market:%d:buy", b.marketID)
			// For buy orders, use negative price for descending sort
			b.client.ZAdd(ctx, buyBookKey, redis.Z{
				Score:  -order.Price.InexactFloat64(),
				Member: order.OrderID,
			})
		} else {
			sellBookKey := fmt.Sprintf("market:%d:sell", b.marketID)
			// For sell orders, use positive price for ascending sort
			b.client.ZAdd(ctx, sellBookKey, redis.Z{
				Score:  order.Price.InexactFloat64(),
				Member: order.OrderID,
			})
		}
	}

	log.Println("Test data seeded")
}

func (b *RedisBot) matchOrders(ctx context.Context) error {
	// Get best buy order (highest price)
	buyBookKey := fmt.Sprintf("market:%d:buy", b.marketID)
	buyOrders, err := b.client.ZRangeWithScores(ctx, buyBookKey, 0, 0).Result()
	if err != nil || len(buyOrders) == 0 {
		return nil // No buy orders
	}

	// Get best sell order (lowest price)
	sellBookKey := fmt.Sprintf("market:%d:sell", b.marketID)
	sellOrders, err := b.client.ZRangeWithScores(ctx, sellBookKey, 0, 0).Result()
	if err != nil || len(sellOrders) == 0 {
		return nil // No sell orders
	}

	// Get order details
	buyOrderID := uint64(buyOrders[0].Member.(int64))
	sellOrderID := uint64(sellOrders[0].Member.(int64))

	buyOrder, err := b.getOrder(ctx, buyOrderID)
	if err != nil {
		return err
	}

	sellOrder, err := b.getOrder(ctx, sellOrderID)
	if err != nil {
		return err
	}

	// Check if prices cross
	if buyOrder.Price.GreaterThanOrEqual(sellOrder.Price) {
		// Create and execute match
		match := b.createMatch(buyOrder, sellOrder, sellOrder.Price)
		if match != nil {
			if err := b.executeMatch(ctx, match); err != nil {
				return err
			}

			log.Printf("Matched: Buy #%d @ %s with Sell #%d @ %s, Amount: %s",
				buyOrder.OrderID, buyOrder.Price.String(),
				sellOrder.OrderID, sellOrder.Price.String(),
				match.Amount.String())

			// Update or remove orders based on remaining amount
			buyOrder.RemainingAmount = buyOrder.RemainingAmount.Sub(match.Amount)
			sellOrder.RemainingAmount = sellOrder.RemainingAmount.Sub(match.Amount)

			if buyOrder.RemainingAmount.IsZero() {
				b.removeOrder(ctx, buyOrder)
			} else {
				b.updateOrder(ctx, buyOrder)
			}

			if sellOrder.RemainingAmount.IsZero() {
				b.removeOrder(ctx, sellOrder)
			} else {
				b.updateOrder(ctx, sellOrder)
			}
		}
	}

	return nil
}

func (b *RedisBot) getOrder(ctx context.Context, orderID uint64) (*Order, error) {
	key := fmt.Sprintf("order:%d", orderID)
	data, err := b.client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var order Order
	if err := json.Unmarshal([]byte(data), &order); err != nil {
		return nil, err
	}

	return &order, nil
}

func (b *RedisBot) createMatch(buyOrder, sellOrder *Order, matchPrice decimal.Decimal) *Match {
	// Calculate match amount
	matchAmount := buyOrder.RemainingAmount
	if sellOrder.RemainingAmount.LessThan(buyOrder.RemainingAmount) {
		matchAmount = sellOrder.RemainingAmount
	}

	if matchAmount.IsZero() {
		return nil
	}

	return &Match{
		BuyOrderID:  buyOrder.OrderID,
		SellOrderID: sellOrder.OrderID,
		Price:       matchPrice,
		Amount:      matchAmount,
		Timestamp:   time.Now(),
	}
}

func (b *RedisBot) executeMatch(ctx context.Context, match *Match) error {
	// Store match record
	matchKey := fmt.Sprintf("match:%d:%d:%d", b.marketID, match.BuyOrderID, match.SellOrderID)
	matchData, err := json.Marshal(match)
	if err != nil {
		return err
	}

	if err := b.client.Set(ctx, matchKey, matchData, 24*time.Hour).Err(); err != nil {
		return err
	}

	// Add to matches list
	matchesKey := fmt.Sprintf("market:%d:matches", b.marketID)
	if err := b.client.LPush(ctx, matchesKey, matchKey).Err(); err != nil {
		return err
	}

	return nil
}

func (b *RedisBot) updateOrder(ctx context.Context, order *Order) error {
	orderKey := fmt.Sprintf("order:%d", order.OrderID)
	orderData, err := json.Marshal(order)
	if err != nil {
		return err
	}
	return b.client.Set(ctx, orderKey, orderData, 0).Err()
}

func (b *RedisBot) removeOrder(ctx context.Context, order *Order) error {
	// Remove from order hash
	orderKey := fmt.Sprintf("order:%d", order.OrderID)
	b.client.Del(ctx, orderKey)

	// Remove from sorted set
	if order.Side == "buy" {
		buyBookKey := fmt.Sprintf("market:%d:buy", b.marketID)
		b.client.ZRem(ctx, buyBookKey, order.OrderID)
	} else {
		sellBookKey := fmt.Sprintf("market:%d:sell", b.marketID)
		b.client.ZRem(ctx, sellBookKey, order.OrderID)
	}

	return nil
}

func (b *RedisBot) Close() error {
	return b.client.Close()
}

func main() {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "redis://localhost:6380"
	}

	mnemonic := "choose isolate cruise nominee image peanut winter vacant enemy improve practice verb moon satisfy food fuel damage sugar load vendor mirror galaxy subject laptop"
	marketID := uint64(1)

	bot := NewRedisBot(redisURL, mnemonic, marketID)
	defer bot.Close()

	ctx := context.Background()
	bot.Start(ctx)
}
