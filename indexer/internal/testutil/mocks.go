package testutil

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"sync"
	"testing"
	"time"

	"cosmossdk.io/math"
	"github.com/elys-network/elys/indexer/internal/models"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

// MockCache provides a mock implementation of the cache.Cache interface
type MockCache struct {
	mu            sync.RWMutex
	data          map[string]interface{}
	subscriptions map[string][]*models.WebSocketSubscription
}

func NewMockCache() *MockCache {
	return &MockCache{
		data:          make(map[string]interface{}),
		subscriptions: make(map[string][]*models.WebSocketSubscription),
	}
}

func (m *MockCache) Get(ctx context.Context, key string) (interface{}, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	val, ok := m.data[key]
	if !ok {
		return nil, redis.Nil
	}
	return val, nil
}

func (m *MockCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.data[key] = value
	return nil
}

func (m *MockCache) Delete(ctx context.Context, keys ...string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, key := range keys {
		delete(m.data, key)
	}
	return nil
}

func (m *MockCache) AddSubscription(ctx context.Context, clientID string, sub *models.WebSocketSubscription) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.subscriptions[clientID] = append(m.subscriptions[clientID], sub)
	return nil
}

func (m *MockCache) GetSubscriptions(ctx context.Context, clientID string) ([]*models.WebSocketSubscription, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.subscriptions[clientID], nil
}

func (m *MockCache) Subscribe(ctx context.Context, channels ...string) *redis.PubSub {
	// Return a mock PubSub that doesn't do anything
	return nil
}

// MockDB provides a mock database connection for testing
type MockDB struct {
	*sql.DB
	execFunc  func(query string, args ...interface{}) (sql.Result, error)
	queryFunc func(query string, args ...interface{}) (*sql.Rows, error)
}

func (m *MockDB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	if m.execFunc != nil {
		return m.execFunc(query, args...)
	}
	return nil, nil
}

func (m *MockDB) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	if m.queryFunc != nil {
		return m.queryFunc(query, args...)
	}
	return nil, nil
}

// TestHelpers provides utility functions for tests
type TestHelpers struct{}

func NewTestHelpers() *TestHelpers {
	return &TestHelpers{}
}

// CreateTestSpotOrder creates a test spot order with default values
func (h *TestHelpers) CreateTestSpotOrder(overrides ...func(*models.SpotOrder)) *models.SpotOrder {
	order := &models.SpotOrder{
		OrderID:          12345,
		OrderType:        models.OrderTypeLimitBuy,
		OwnerAddress:     "elys1test123",
		OrderTargetDenom: "USDC",
		OrderPrice:       models.JSONB{"price": "10.50"},
		Status:           models.OrderStatusPending,
		CreatedAt:        time.Now(),
		BlockHeight:      100,
		TxHash:           "0xtest123",
	}

	for _, override := range overrides {
		override(order)
	}

	return order
}

// CreateTestPerpetualPosition creates a test perpetual position with default values
func (h *TestHelpers) CreateTestPerpetualPosition(overrides ...func(*models.PerpetualPosition)) *models.PerpetualPosition {
	pos := &models.PerpetualPosition{
		MtpID:           98765,
		OwnerAddress:    "elys1test456",
		AmmPoolID:       1,
		Position:        models.PositionTypeLong,
		CollateralAsset: "USDC",
		OpenedAt:        time.Now(),
		BlockHeight:     200,
		TxHash:          "0xpos123",
	}

	for _, override := range overrides {
		override(pos)
	}

	return pos
}

// CreateTestTrade creates a test trade with default values
func (h *TestHelpers) CreateTestTrade(overrides ...func(*models.Trade)) *models.Trade {
	trade := &models.Trade{
		TradeType:    "SPOT",
		ReferenceID:  123,
		OwnerAddress: "elys1test789",
		Asset:        "ATOM",
		Fees:         models.JSONB{"amount": "0.25"},
		ExecutedAt:   time.Now(),
		BlockHeight:  300,
		TxHash:       "0xtrade123",
		EventType:    "OrderExecuted",
	}

	for _, override := range overrides {
		override(trade)
	}

	return trade
}

// CreateTestOrderBookSnapshot creates a test order book snapshot
func (h *TestHelpers) CreateTestOrderBookSnapshot(assetPair string) *models.OrderBookSnapshot {
	return &models.OrderBookSnapshot{
		AssetPair: assetPair,
		Bids: []models.OrderBookEntry{
			{Price: math.LegacyNewDec(995).Quo(math.LegacyNewDec(100)), Amount: math.NewInt(100)},
			{Price: math.LegacyNewDec(990).Quo(math.LegacyNewDec(100)), Amount: math.NewInt(200)},
			{Price: math.LegacyNewDec(985).Quo(math.LegacyNewDec(100)), Amount: math.NewInt(300)},
		},
		Asks: []models.OrderBookEntry{
			{Price: math.LegacyNewDec(1005).Quo(math.LegacyNewDec(100)), Amount: math.NewInt(150)},
			{Price: math.LegacyNewDec(1010).Quo(math.LegacyNewDec(100)), Amount: math.NewInt(250)},
			{Price: math.LegacyNewDec(1015).Quo(math.LegacyNewDec(100)), Amount: math.NewInt(350)},
		},
		SnapshotTime: time.Now(),
		BlockHeight:  400,
	}
}

// CreateTestWebSocketMessage creates a test WebSocket message
func (h *TestHelpers) CreateTestWebSocketMessage(msgType models.WSMessageType, data interface{}) *models.WSMessage {
	jsonData, _ := json.Marshal(data)
	return &models.WSMessage{
		Type: msgType,
		Data: json.RawMessage(jsonData),
	}
}

// AssertEventuallyTrue waits for a condition to become true
func (h *TestHelpers) AssertEventuallyTrue(t *testing.T, condition func() bool, timeout time.Duration, msg string) {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if condition() {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("Condition never became true: %s", msg)
}

// WaitForWebSocketConnection waits for a WebSocket connection to be established
func (h *TestHelpers) WaitForWebSocketConnection(ws *websocket.Conn, timeout time.Duration) error {
	done := make(chan struct{})
	go func() {
		// Send a ping to verify connection
		err := ws.WriteMessage(websocket.PingMessage, nil)
		if err == nil {
			close(done)
		}
	}()

	select {
	case <-done:
		return nil
	case <-time.After(timeout):
		return fmt.Errorf("websocket connection timeout")
	}
}
