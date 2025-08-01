package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/shopspring/decimal"
)

type PostgresBot struct {
	db               *sql.DB
	mnemonic         string
	marketID         uint64
	matchingInterval time.Duration
}

type Order struct {
	OrderID         uint64
	MarketID        uint64
	Side            string
	Price           decimal.Decimal
	RemainingAmount decimal.Decimal
	Owner           string
	CreatedAt       time.Time
}

type Match struct {
	BuyOrderID  uint64
	SellOrderID uint64
	Price       decimal.Decimal
	Amount      decimal.Decimal
	ExecutedAt  time.Time
}

func NewPostgresBot(dbURL string, mnemonic string, marketID uint64) *PostgresBot {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Create tables if they don't exist
	if err := createTables(db); err != nil {
		log.Fatalf("Failed to create tables: %v", err)
	}

	return &PostgresBot{
		db:               db,
		mnemonic:         mnemonic,
		marketID:         marketID,
		matchingInterval: 500 * time.Millisecond,
	}
}

func createTables(db *sql.DB) error {
	// Create simple tables for testing
	queries := []string{
		`CREATE TABLE IF NOT EXISTS clob_orders (
			order_id BIGINT PRIMARY KEY,
			market_id BIGINT NOT NULL,
			side VARCHAR(10) NOT NULL,
			price DECIMAL(20,8) NOT NULL,
			remaining_amount DECIMAL(20,8) NOT NULL,
			owner VARCHAR(100) NOT NULL,
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NOW()
		)`,
		`CREATE TABLE IF NOT EXISTS clob_matches (
			match_id SERIAL PRIMARY KEY,
			buy_order_id BIGINT NOT NULL,
			sell_order_id BIGINT NOT NULL,
			market_id BIGINT NOT NULL,
			price DECIMAL(20,8) NOT NULL,
			amount DECIMAL(20,8) NOT NULL,
			executed_at TIMESTAMP DEFAULT NOW()
		)`,
		`CREATE INDEX IF NOT EXISTS idx_orders_market_side_price ON clob_orders(market_id, side, price)`,
		`CREATE INDEX IF NOT EXISTS idx_matches_market ON clob_matches(market_id)`,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return err
		}
	}

	return nil
}

func (b *PostgresBot) Start(ctx context.Context) {
	log.Printf("Starting PostgreSQL bot for market %d", b.marketID)
	log.Printf("Bot address derived from mnemonic: %s...", b.mnemonic[:20])

	// Seed test data
	b.seedTestData()

	ticker := time.NewTicker(b.matchingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("PostgreSQL bot stopped")
			return
		case <-ticker.C:
			if err := b.matchOrders(ctx); err != nil {
				log.Printf("Error matching orders: %v", err)
			}
		}
	}
}

func (b *PostgresBot) seedTestData() {
	// Check if we already have data
	var count int
	b.db.QueryRow("SELECT COUNT(*) FROM clob_orders WHERE market_id = $1", b.marketID).Scan(&count)
	if count > 0 {
		return // Already have data
	}

	// Insert test orders
	testOrders := []struct {
		orderID uint64
		side    string
		price   string
		amount  string
		owner   string
	}{
		{1, "buy", "100", "10", "test1"},
		{2, "buy", "99", "5", "test2"},
		{3, "sell", "101", "8", "test3"},
		{4, "sell", "102", "12", "test4"},
	}

	for _, order := range testOrders {
		_, err := b.db.Exec(`
			INSERT INTO clob_orders (order_id, market_id, side, price, remaining_amount, owner)
			VALUES ($1, $2, $3, $4, $5, $6)
			ON CONFLICT (order_id) DO NOTHING`,
			order.orderID, b.marketID, order.side, order.price, order.amount, order.owner)
		if err != nil {
			log.Printf("Failed to insert test order: %v", err)
		}
	}

	log.Println("Test data seeded")
}

func (b *PostgresBot) matchOrders(ctx context.Context) error {
	// Get best buy order (highest price)
	var buyOrder Order
	err := b.db.QueryRowContext(ctx, `
		SELECT order_id, market_id, side, price, remaining_amount, owner, created_at
		FROM clob_orders
		WHERE market_id = $1 AND side = 'buy' AND remaining_amount > 0
		ORDER BY price DESC, created_at ASC
		LIMIT 1`, b.marketID).Scan(
		&buyOrder.OrderID, &buyOrder.MarketID, &buyOrder.Side,
		&buyOrder.Price, &buyOrder.RemainingAmount, &buyOrder.Owner, &buyOrder.CreatedAt)

	if err == sql.ErrNoRows {
		return nil // No buy orders
	} else if err != nil {
		return fmt.Errorf("failed to get buy order: %w", err)
	}

	// Get best sell order (lowest price)
	var sellOrder Order
	err = b.db.QueryRowContext(ctx, `
		SELECT order_id, market_id, side, price, remaining_amount, owner, created_at
		FROM clob_orders
		WHERE market_id = $1 AND side = 'sell' AND remaining_amount > 0
		ORDER BY price ASC, created_at ASC
		LIMIT 1`, b.marketID).Scan(
		&sellOrder.OrderID, &sellOrder.MarketID, &sellOrder.Side,
		&sellOrder.Price, &sellOrder.RemainingAmount, &sellOrder.Owner, &sellOrder.CreatedAt)

	if err == sql.ErrNoRows {
		return nil // No sell orders
	} else if err != nil {
		return fmt.Errorf("failed to get sell order: %w", err)
	}

	// Check if prices cross
	if buyOrder.Price.GreaterThanOrEqual(sellOrder.Price) {
		match := b.createMatch(&buyOrder, &sellOrder, sellOrder.Price)
		if match != nil {
			if err := b.executeMatch(ctx, match); err != nil {
				return fmt.Errorf("failed to execute match: %w", err)
			}

			log.Printf("Matched: Buy #%d @ %s with Sell #%d @ %s, Amount: %s",
				buyOrder.OrderID, buyOrder.Price.String(),
				sellOrder.OrderID, sellOrder.Price.String(),
				match.Amount.String())
		}
	}

	return nil
}

func (b *PostgresBot) createMatch(buyOrder, sellOrder *Order, matchPrice decimal.Decimal) *Match {
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
		ExecutedAt:  time.Now(),
	}
}

func (b *PostgresBot) executeMatch(ctx context.Context, match *Match) error {
	// Use transaction for consistency
	tx, err := b.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert match record
	_, err = tx.ExecContext(ctx, `
		INSERT INTO clob_matches (buy_order_id, sell_order_id, market_id, price, amount, executed_at)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		match.BuyOrderID, match.SellOrderID, b.marketID,
		match.Price.String(), match.Amount.String(), match.ExecutedAt)
	if err != nil {
		return fmt.Errorf("failed to insert match: %w", err)
	}

	// Update buy order
	_, err = tx.ExecContext(ctx, `
		UPDATE clob_orders
		SET remaining_amount = remaining_amount - $1, updated_at = NOW()
		WHERE order_id = $2 AND remaining_amount >= $1`,
		match.Amount.String(), match.BuyOrderID)
	if err != nil {
		return fmt.Errorf("failed to update buy order: %w", err)
	}

	// Update sell order
	_, err = tx.ExecContext(ctx, `
		UPDATE clob_orders
		SET remaining_amount = remaining_amount - $1, updated_at = NOW()
		WHERE order_id = $2 AND remaining_amount >= $1`,
		match.Amount.String(), match.SellOrderID)
	if err != nil {
		return fmt.Errorf("failed to update sell order: %w", err)
	}

	return tx.Commit()
}

func (b *PostgresBot) Close() error {
	return b.db.Close()
}

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://indexer:123123123@localhost:5433/elys_indexer?sslmode=disable"
	}

	mnemonic := "choose isolate cruise nominee image peanut winter vacant enemy improve practice verb moon satisfy food fuel damage sugar load vendor mirror galaxy subject laptop"
	marketID := uint64(1)

	bot := NewPostgresBot(dbURL, mnemonic, marketID)
	defer bot.Close()

	ctx := context.Background()
	bot.Start(ctx)
}
