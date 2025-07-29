package database

import (
	"context"
	"testing"

	"cosmossdk.io/math"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/elys-network/elys/indexer/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// TestSimpleSpotOrderCreate tests basic spot order creation
func TestSimpleSpotOrderCreate(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	logger := zap.NewNop()
	repo := NewRepository(&DB{DB: db, logger: logger}, logger)

	ctx := context.Background()
	order := &models.SpotOrder{
		OrderID:          12345,
		OrderType:        models.OrderTypeLimitBuy,
		OwnerAddress:     "elys1test",
		OrderTargetDenom: "USDC",
		OrderPrice:       models.JSONB{"price": "1.25"},
		OrderAmount:      math.NewInt(1000),
		Status:           models.OrderStatusPending,
	}

	// Expect the INSERT query
	mock.ExpectExec("INSERT INTO spot_orders").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.CreateSpotOrder(ctx, order)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestSimplePerpetualOrderCreate tests basic perpetual order creation
func TestSimplePerpetualOrderCreate(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	logger := zap.NewNop()
	repo := NewRepository(&DB{DB: db, logger: logger}, logger)

	ctx := context.Background()
	order := &models.PerpetualOrder{
		OrderID:            54321,
		OwnerAddress:       "elys1test",
		PerpetualOrderType: models.PerpetualOrderTypeLimitOpen,
		Position:           models.PositionTypeLong,
		TriggerPrice:       models.JSONB{"price": "1500.00"},
		Collateral:         math.NewInt(10000),
		Leverage:           math.LegacyNewDec(10),
		TakeProfitPrice:    math.LegacyNewDec(1600),
		StopLossPrice:      math.LegacyNewDec(1400),
		Status:             models.OrderStatusPending,
	}

	// Expect the INSERT query
	mock.ExpectExec("INSERT INTO perpetual_orders").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.CreatePerpetualOrder(ctx, order)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestSimpleTradeCreate tests basic trade creation
func TestSimpleTradeCreate(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	logger := zap.NewNop()
	repo := NewRepository(&DB{DB: db, logger: logger}, logger)

	ctx := context.Background()
	trade := &models.Trade{
		TradeType:    "SPOT",
		ReferenceID:  123,
		OwnerAddress: "elys1test",
		Asset:        "ATOM",
		Amount:       math.NewInt(100),
		Price:        math.LegacyNewDec(10),
		Fees:         models.JSONB{"amount": "0.25"},
		EventType:    "OrderExecuted",
	}

	// Expect the INSERT query
	mock.ExpectExec("INSERT INTO trades").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.CreateTrade(ctx, trade)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestSimpleWebSocketSubscription tests websocket subscription operations
func TestSimpleWebSocketSubscription(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	logger := zap.NewNop()
	repo := NewRepository(&DB{DB: db, logger: logger}, logger)

	ctx := context.Background()
	sub := &models.WebSocketSubscription{
		ID:               "sub123",
		ClientID:         "client123",
		SubscriptionType: "orders",
		Filters:          models.JSONB{"owner": "elys1test"},
	}

	// Test Create
	mock.ExpectExec("INSERT INTO websocket_subscriptions").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.CreateSubscription(ctx, sub)
	assert.NoError(t, err)

	// Test Update Ping
	mock.ExpectExec("UPDATE websocket_subscriptions SET last_ping").
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.UpdateSubscriptionPing(ctx, sub.ClientID)
	assert.NoError(t, err)

	// Test Delete
	mock.ExpectExec("DELETE FROM websocket_subscriptions WHERE id").
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.DeleteSubscription(ctx, sub.ID)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}