// +build integration

package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"cosmossdk.io/math"
	"github.com/elys-network/elys/indexer/internal/config"
	"github.com/elys-network/elys/indexer/internal/models"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// TestDatabase provides integration tests against a real PostgreSQL database
// Run with: go test -tags=integration ./internal/database/...
type TestDatabase struct {
	DB   *DB
	Repo *Repository
	ctx  context.Context
}

func setupIntegrationTest(t *testing.T) *TestDatabase {
	// Get database connection from environment
	dbHost := os.Getenv("TEST_DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}
	
	dbPort := os.Getenv("TEST_DB_PORT")
	if dbPort == "" {
		dbPort = "5432"
	}
	
	dbUser := os.Getenv("TEST_DB_USER")
	if dbUser == "" {
		dbUser = "postgres"
	}
	
	dbPass := os.Getenv("TEST_DB_PASS")
	if dbPass == "" {
		dbPass = "postgres"
	}
	
	dbName := os.Getenv("TEST_DB_NAME")
	if dbName == "" {
		dbName = "elys_indexer_test"
	}

	// Create test database if it doesn't exist
	masterDSN := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=postgres sslmode=disable",
		dbHost, dbPort, dbUser, dbPass)
	
	masterDB, err := sql.Open("postgres", masterDSN)
	require.NoError(t, err)
	defer masterDB.Close()
	
	// Create test database
	_, err = masterDB.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))
	if err != nil && !isAlreadyExistsError(err) {
		t.Fatalf("Failed to create test database: %v", err)
	}

	// Connect to test database
	port, _ := strconv.Atoi(dbPort)
	cfg := &config.DatabaseConfig{
		Host:            dbHost,
		Port:            port,
		User:            dbUser,
		Password:        dbPass,
		Database:        dbName,
		SSLMode:         "disable",
		MaxOpenConns:    10,
		MaxIdleConns:    5,
		ConnMaxLifetime: time.Hour,
	}
	
	logger := zap.NewNop()
	db, err := New(cfg, logger)
	require.NoError(t, err)
	
	// Run migrations
	err = runMigrations(db.DB, "./../../sql/schema.sql")
	require.NoError(t, err)
	
	// Clear all tables before tests
	clearAllTables(t, db.DB)
	
	repo := NewRepository(db, logger)
	
	return &TestDatabase{
		DB:   db,
		Repo: repo,
		ctx:  context.Background(),
	}
}

func (td *TestDatabase) cleanup() {
	td.DB.Close()
}

func isAlreadyExistsError(err error) bool {
	return err != nil && err.Error() == "pq: database \"elys_indexer_test\" already exists"
}

func runMigrations(db *sql.DB, schemaPath string) error {
	schema, err := os.ReadFile(schemaPath)
	if err != nil {
		return fmt.Errorf("failed to read schema file: %w", err)
	}
	
	_, err = db.Exec(string(schema))
	return err
}

func clearAllTables(t *testing.T, db *sql.DB) {
	tables := []string{
		"trades",
		"spot_orders",
		"perpetual_orders",
		"perpetual_positions",
		"order_book_snapshots",
		"websocket_subscriptions",
		"indexer_state",
	}
	
	for _, table := range tables {
		_, err := db.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table))
		require.NoError(t, err)
	}
}

func TestSpotOrderLifecycle(t *testing.T) {
	td := setupIntegrationTest(t)
	defer td.cleanup()
	
	// Create a spot order
	order := &models.SpotOrder{
		OrderID:          uint64(time.Now().UnixNano()),
		OrderType:        models.OrderTypeLimitBuy,
		OwnerAddress:     "elys1test123",
		OrderTargetDenom: "USDC",
		OrderPrice:       models.JSONB{"price": "10.50"},
		OrderAmount:      math.NewInt(1000),
		Status:           models.OrderStatusPending,
		CreatedAt:        time.Now(),
		BlockHeight:      100,
		TxHash:          "0xabc123",
	}
	
	// Test Create
	err := td.Repo.CreateSpotOrder(td.ctx, order)
	assert.NoError(t, err)
	
	// Test Get by Owner
	orders, err := td.Repo.GetSpotOrdersByOwner(td.ctx, order.OwnerAddress, 10)
	assert.NoError(t, err)
	assert.Len(t, orders, 1)
	assert.Equal(t, order.OrderID, orders[0].OrderID)
	
	// Test Update Status
	err = td.Repo.UpdateSpotOrderStatus(td.ctx, order.OrderID, models.OrderStatusExecuted)
	assert.NoError(t, err)
	
	// Verify update
	orders, err = td.Repo.GetSpotOrdersByOwner(td.ctx, order.OwnerAddress, 10)
	assert.NoError(t, err)
	assert.Equal(t, models.OrderStatusExecuted, orders[0].Status)
	assert.NotNil(t, orders[0].ExecutedAt)
}

func TestPerpetualPositionLifecycle(t *testing.T) {
	td := setupIntegrationTest(t)
	defer td.cleanup()
	
	// Create a perpetual position
	pos := &models.PerpetualPosition{
		MtpID:           uint64(time.Now().UnixNano()),
		OwnerAddress:    "elys1test456",
		AmmPoolID:       1,
		Position:        models.PositionTypeLong,
		CollateralAsset: "USDC",
		Collateral:      math.NewInt(5000),
		Liabilities:     math.NewInt(45000),
		Custody:         math.NewInt(50000),
		MtpHealth:       math.LegacyNewDec(100),
		OpenPrice:       math.LegacyNewDec(1450),
		StopLossPrice:   math.LegacyNewDec(1400),
		TakeProfitPrice: math.LegacyNewDec(1600),
		OpenedAt:        time.Now(),
		BlockHeight:     200,
		TxHash:          "0xpos123",
	}
	
	// Test Create
	err := td.Repo.CreatePerpetualPosition(td.ctx, pos)
	assert.NoError(t, err)
	
	// Test Get Open Positions
	positions, err := td.Repo.GetOpenPositions(td.ctx, pos.OwnerAddress)
	assert.NoError(t, err)
	assert.Len(t, positions, 1)
	assert.Equal(t, pos.MtpID, positions[0].MtpID)
	
	// Test Close Position
	err = td.Repo.ClosePerpetualPosition(td.ctx, pos.MtpID, "1500.00", "50.00", "USER", "MANUAL")
	assert.NoError(t, err)
	
	// Verify position is closed
	positions, err = td.Repo.GetOpenPositions(td.ctx, pos.OwnerAddress)
	assert.NoError(t, err)
	assert.Len(t, positions, 0)
}

func TestTradeHistory(t *testing.T) {
	td := setupIntegrationTest(t)
	defer td.cleanup()
	
	owner := "elys1trader"
	
	// Create multiple trades
	for i := 0; i < 5; i++ {
		trade := &models.Trade{
			TradeType:    "SPOT",
			ReferenceID:  uint64(i),
			OwnerAddress: owner,
			Asset:        "ATOM",
			Amount:       math.NewInt(int64(100 * (i + 1))),
			Price:        math.LegacyNewDec(int64(10 + i)),
			Fees:         models.JSONB{"amount": fmt.Sprintf("0.%d", i)},
			ExecutedAt:   time.Now().Add(time.Duration(i) * time.Minute),
			BlockHeight:  int64(300 + i),
			TxHash:       fmt.Sprintf("0xtrade%d", i),
			EventType:    "OrderExecuted",
		}
		
		err := td.Repo.CreateTrade(td.ctx, trade)
		assert.NoError(t, err)
	}
	
	// Test Get Trade History
	trades, err := td.Repo.GetTradeHistory(td.ctx, owner, 3)
	assert.NoError(t, err)
	assert.Len(t, trades, 3)
	
	// Verify order (should be descending by executed_at)
	for i := 0; i < len(trades)-1; i++ {
		assert.True(t, trades[i].ExecutedAt.After(trades[i+1].ExecutedAt))
	}
}

func TestOrderBookSnapshots(t *testing.T) {
	td := setupIntegrationTest(t)
	defer td.cleanup()
	
	assetPair := "ATOM/USDC"
	
	// Create snapshots at different times
	for i := 0; i < 3; i++ {
		snapshot := &models.OrderBookSnapshot{
			AssetPair: assetPair,
			Bids: []models.OrderBookEntry{
				{Price: math.LegacyMustNewDecFromStr(fmt.Sprintf("9.%d5", 9-i)), Amount: math.NewInt(100)},
				{Price: math.LegacyMustNewDecFromStr(fmt.Sprintf("9.%d0", 9-i)), Amount: math.NewInt(200)},
			},
			Asks: []models.OrderBookEntry{
				{Price: math.LegacyMustNewDecFromStr(fmt.Sprintf("10.%d5", i)), Amount: math.NewInt(150)},
				{Price: math.LegacyMustNewDecFromStr(fmt.Sprintf("10.%d0", 1+i)), Amount: math.NewInt(250)},
			},
			SnapshotTime: time.Now().Add(time.Duration(i) * time.Second),
			BlockHeight:  int64(400 + i),
		}
		
		err := td.Repo.SaveOrderBookSnapshot(td.ctx, snapshot)
		assert.NoError(t, err)
		
		time.Sleep(10 * time.Millisecond) // Ensure different timestamps
	}
	
	// Test Get Latest Snapshot
	latest, err := td.Repo.GetLatestOrderBookSnapshot(td.ctx, assetPair)
	assert.NoError(t, err)
	assert.NotNil(t, latest)
	assert.Equal(t, int64(402), latest.BlockHeight)
	assert.Len(t, latest.Bids, 2)
	assert.Len(t, latest.Asks, 2)
}

func TestIndexerState(t *testing.T) {
	td := setupIntegrationTest(t)
	defer td.cleanup()
	
	// Test initial state (no rows)
	state, err := td.Repo.GetIndexerState(td.ctx)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), state.LastProcessedHeight)
	
	// Initialize indexer state
	_, err = td.DB.ExecContext(td.ctx, 
		"INSERT INTO indexer_state (last_processed_height, last_processed_time) VALUES ($1, $2)",
		int64(0), time.Now())
	assert.NoError(t, err)
	
	// Test Update State
	newHeight := int64(1000)
	err = td.Repo.UpdateIndexerState(td.ctx, newHeight)
	assert.NoError(t, err)
	
	// Verify update
	state, err = td.Repo.GetIndexerState(td.ctx)
	assert.NoError(t, err)
	assert.Equal(t, newHeight, state.LastProcessedHeight)
	assert.WithinDuration(t, time.Now(), state.LastProcessedTime, 5*time.Second)
}

func TestWebSocketSubscriptionManagement(t *testing.T) {
	td := setupIntegrationTest(t)
	defer td.cleanup()
	
	clientID := "test-client-123"
	
	// Create subscriptions
	sub1 := &models.WebSocketSubscription{
		ID:               "sub1",
		ClientID:         clientID,
		SubscriptionType: "orders",
		Filters:          map[string]interface{}{"owner": "elys1test"},
		CreatedAt:        time.Now(),
		LastPing:         time.Now(),
	}
	
	sub2 := &models.WebSocketSubscription{
		ID:               "sub2",
		ClientID:         clientID,
		SubscriptionType: "orderbook",
		Filters:          map[string]interface{}{"pair": "ATOM/USDC"},
		CreatedAt:        time.Now(),
		LastPing:         time.Now(),
	}
	
	// Test Create
	err := td.Repo.CreateSubscription(td.ctx, sub1)
	assert.NoError(t, err)
	
	err = td.Repo.CreateSubscription(td.ctx, sub2)
	assert.NoError(t, err)
	
	// Test Update Ping
	time.Sleep(100 * time.Millisecond)
	err = td.Repo.UpdateSubscriptionPing(td.ctx, clientID)
	assert.NoError(t, err)
	
	// Test Delete
	err = td.Repo.DeleteSubscription(td.ctx, sub1.ID)
	assert.NoError(t, err)
	
	// Test Cleanup Stale
	// Create an old subscription
	oldSub := &models.WebSocketSubscription{
		ID:               "old-sub",
		ClientID:         "old-client",
		SubscriptionType: "trades",
		Filters:          map[string]interface{}{},
		CreatedAt:        time.Now().Add(-time.Hour),
		LastPing:         time.Now().Add(-time.Hour),
	}
	
	err = td.Repo.CreateSubscription(td.ctx, oldSub)
	assert.NoError(t, err)
	
	// Cleanup subscriptions older than 5 minutes
	err = td.Repo.CleanupStaleSubscriptions(td.ctx, 5*time.Minute)
	assert.NoError(t, err)
	
	// Verify old subscription was deleted
	var count int
	err = td.DB.QueryRowContext(td.ctx, 
		"SELECT COUNT(*) FROM websocket_subscriptions WHERE id = $1", oldSub.ID).Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 0, count)
}

func TestConcurrentOperations(t *testing.T) {
	td := setupIntegrationTest(t)
	defer td.cleanup()
	
	// Test concurrent writes
	done := make(chan bool, 3)
	
	// Concurrent spot order creation
	go func() {
		for i := 0; i < 10; i++ {
			order := &models.SpotOrder{
				OrderID:          uint64(time.Now().UnixNano() + int64(i)),
				OrderType:        models.OrderTypeLimitBuy,
				OwnerAddress:     fmt.Sprintf("elys1concurrent%d", i),
				OrderTargetDenom: "USDC",
				OrderPrice:       models.JSONB{"price": "10.50"},
				OrderAmount:      math.NewInt(1000),
				Status:           models.OrderStatusPending,
				CreatedAt:        time.Now(),
				BlockHeight:      int64(100 + i),
				TxHash:          fmt.Sprintf("0xconcurrent%d", i),
			}
			
			err := td.Repo.CreateSpotOrder(td.ctx, order)
			assert.NoError(t, err)
		}
		done <- true
	}()
	
	// Concurrent trade creation
	go func() {
		for i := 0; i < 10; i++ {
			trade := &models.Trade{
				TradeType:    "SPOT",
				ReferenceID:  uint64(1000 + i),
				OwnerAddress: "elys1concurrent",
				Asset:        "ATOM",
				Amount:       math.NewInt(100),
				Price:        math.LegacyNewDec(10),
				Fees:         models.JSONB{"amount": "0.25"},
				ExecutedAt:   time.Now(),
				BlockHeight:  int64(200 + i),
				TxHash:       fmt.Sprintf("0xtradeconcurrent%d", i),
				EventType:    "OrderExecuted",
			}
			
			err := td.Repo.CreateTrade(td.ctx, trade)
			assert.NoError(t, err)
		}
		done <- true
	}()
	
	// Concurrent indexer state updates
	go func() {
		for i := 0; i < 10; i++ {
			err := td.Repo.UpdateIndexerState(td.ctx, int64(1000+i))
			assert.NoError(t, err)
			time.Sleep(10 * time.Millisecond)
		}
		done <- true
	}()
	
	// Wait for all goroutines to complete
	for i := 0; i < 3; i++ {
		<-done
	}
	
	// Verify final state
	state, err := td.Repo.GetIndexerState(td.ctx)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, state.LastProcessedHeight, int64(1009))
}