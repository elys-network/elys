package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authclient "github.com/cosmos/cosmos-sdk/x/auth/client"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/go-bip39"
	clobtypes "github.com/elys-network/elys/v7/x/clob/types"
	_ "github.com/lib/pq"
	"github.com/shopspring/decimal"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type PostgresBot struct {
	db               *sql.DB
	mnemonic         string
	marketID         uint64
	matchingInterval time.Duration
	// Blockchain interaction
	grpcConn   *grpc.ClientConn
	txClient   tx.ServiceClient
	authClient authtypes.QueryClient
	privKey    cryptotypes.PrivKey
	address    sdk.AccAddress
	clientCtx  client.Context
}

type Order struct {
	OrderID         uint64
	MarketID        uint64
	Side            string
	Price           decimal.Decimal
	RemainingAmount decimal.Decimal
	Owner           string
	CreatedAt       time.Time
	Counter         uint64 // Order counter for OrderKey
}

type Match struct {
	BuyOrderID  uint64
	SellOrderID uint64
	Price       decimal.Decimal
	Amount      decimal.Decimal
	ExecutedAt  time.Time
}

func NewPostgresBot(dbURL string, mnemonic string, marketID uint64, grpcEndpoint string, chainID string) (*PostgresBot, error) {
	// Database setup
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Create tables if they don't exist
	if err := createTables(db); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	// Derive private key from mnemonic
	if !bip39.IsMnemonic(mnemonic) {
		return nil, fmt.Errorf("invalid mnemonic")
	}

	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, "")
	if err != nil {
		return nil, fmt.Errorf("failed to create seed: %w", err)
	}

	master, ch := hd.ComputeMastersFromSeed(seed)
	path := "m/44'/118'/0'/0/0" // Standard Cosmos derivation path
	privKey, err := hd.DerivePrivateKeyForPath(master, ch, path)
	if err != nil {
		return nil, fmt.Errorf("failed to derive private key: %w", err)
	}

	// Get address
	pubKey := privKey.PubKey()
	address := sdk.AccAddress(pubKey.Address())

	// Create gRPC connection
	grpcConn, err := grpc.Dial(
		grpcEndpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to gRPC: %w", err)
	}

	// Create client context
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	authtypes.RegisterInterfaces(interfaceRegistry)
	clobtypes.RegisterInterfaces(interfaceRegistry)

	cdc := codec.NewProtoCodec(interfaceRegistry)
	txConfig := tx.NewTxConfig(cdc, tx.DefaultSignModes)

	clientCtx := client.Context{}.
		WithChainID(chainID).
		WithCodec(cdc).
		WithTxConfig(txConfig).
		WithInterfaceRegistry(interfaceRegistry).
		WithAccountRetriever(authtypes.AccountRetriever{})

	return &PostgresBot{
		db:               db,
		mnemonic:         mnemonic,
		marketID:         marketID,
		matchingInterval: 500 * time.Millisecond,
		grpcConn:         grpcConn,
		txClient:         tx.NewServiceClient(grpcConn),
		authClient:       authtypes.NewQueryClient(grpcConn),
		privKey:          privKey,
		address:          address,
		clientCtx:        clientCtx,
	}, nil
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
			updated_at TIMESTAMP DEFAULT NOW(),
			counter BIGINT NOT NULL DEFAULT 0
		)`,
		`CREATE TABLE IF NOT EXISTS clob_matches (
			match_id SERIAL PRIMARY KEY,
			buy_order_id BIGINT NOT NULL,
			sell_order_id BIGINT NOT NULL,
			market_id BIGINT NOT NULL,
			price DECIMAL(20,8) NOT NULL,
			amount DECIMAL(20,8) NOT NULL,
			executed_at TIMESTAMP DEFAULT NOW(),
			tx_hash VARCHAR(100)
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
	log.Printf("Bot address: %s", b.address.String())

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

	// Insert test orders with counters
	testOrders := []struct {
		orderID uint64
		side    string
		price   string
		amount  string
		owner   string
		counter uint64
	}{
		{1, "buy", "100", "10", "test1", 1},
		{2, "buy", "99", "5", "test2", 2},
		{3, "sell", "101", "8", "test3", 3},
		{4, "sell", "102", "12", "test4", 4},
	}

	for _, order := range testOrders {
		_, err := b.db.Exec(`
			INSERT INTO clob_orders (order_id, market_id, side, price, remaining_amount, owner, counter)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			ON CONFLICT (order_id) DO NOTHING`,
			order.orderID, b.marketID, order.side, order.price, order.amount, order.owner, order.counter)
		if err != nil {
			log.Printf("Failed to insert test order: %v", err)
		}
	}

	log.Println("Test data seeded")
}

func (b *PostgresBot) matchOrders(ctx context.Context) error {
	// Get all crossing orders
	crossingBuyOrders, crossingSellOrders, err := b.findCrossingOrders(ctx)
	if err != nil {
		return fmt.Errorf("failed to find crossing orders: %w", err)
	}

	if len(crossingBuyOrders) == 0 || len(crossingSellOrders) == 0 {
		return nil // No orders to match
	}

	// Execute match on blockchain by calling ExecuteMarket for this market
	txHash, err := b.executeMarketOnChain(ctx)
	if err != nil {
		log.Printf("Failed to execute market on chain: %v", err)
		// Continue with local recording even if blockchain tx fails
		txHash = ""
	}

	// Record matches locally
	for _, buyOrder := range crossingBuyOrders {
		for _, sellOrder := range crossingSellOrders {
			if buyOrder.Price.GreaterThanOrEqual(sellOrder.Price) {
				match := b.createMatch(&buyOrder, &sellOrder, sellOrder.Price)
				if match != nil {
					if err := b.executeMatch(ctx, match, txHash); err != nil {
						log.Printf("Failed to execute match: %v", err)
					}

					log.Printf("Matched: Buy #%d @ %s with Sell #%d @ %s, Amount: %s",
						buyOrder.OrderID, buyOrder.Price.String(),
						sellOrder.OrderID, sellOrder.Price.String(),
						match.Amount.String())

					// Update remaining amounts for next iteration
					buyOrder.RemainingAmount = buyOrder.RemainingAmount.Sub(match.Amount)
					sellOrder.RemainingAmount = sellOrder.RemainingAmount.Sub(match.Amount)

					if buyOrder.RemainingAmount.IsZero() {
						break // Move to next buy order
					}
				}
			}
		}
	}

	return nil
}

func (b *PostgresBot) findCrossingOrders(ctx context.Context) ([]Order, []Order, error) {
	// Get all buy orders sorted by price descending
	buyQuery := `
		SELECT order_id, market_id, side, price, remaining_amount, owner, created_at, counter
		FROM clob_orders
		WHERE market_id = $1 AND side = 'buy' AND remaining_amount > 0
		ORDER BY price DESC, created_at ASC`

	buyRows, err := b.db.QueryContext(ctx, buyQuery, b.marketID)
	if err != nil {
		return nil, nil, err
	}
	defer buyRows.Close()

	var buyOrders []Order
	for buyRows.Next() {
		var order Order
		err := buyRows.Scan(&order.OrderID, &order.MarketID, &order.Side,
			&order.Price, &order.RemainingAmount, &order.Owner, &order.CreatedAt, &order.Counter)
		if err != nil {
			return nil, nil, err
		}
		buyOrders = append(buyOrders, order)
	}

	// Get all sell orders sorted by price ascending
	sellQuery := `
		SELECT order_id, market_id, side, price, remaining_amount, owner, created_at, counter
		FROM clob_orders
		WHERE market_id = $1 AND side = 'sell' AND remaining_amount > 0
		ORDER BY price ASC, created_at ASC`

	sellRows, err := b.db.QueryContext(ctx, sellQuery, b.marketID)
	if err != nil {
		return nil, nil, err
	}
	defer sellRows.Close()

	var sellOrders []Order
	for sellRows.Next() {
		var order Order
		err := sellRows.Scan(&order.OrderID, &order.MarketID, &order.Side,
			&order.Price, &order.RemainingAmount, &order.Owner, &order.CreatedAt, &order.Counter)
		if err != nil {
			return nil, nil, err
		}
		sellOrders = append(sellOrders, order)
	}

	if len(buyOrders) == 0 || len(sellOrders) == 0 {
		return nil, nil, nil
	}

	// Find crossing orders
	var crossingBuyOrders []Order
	var crossingSellOrders []Order

	lowestSellPrice := sellOrders[0].Price
	highestBuyPrice := buyOrders[0].Price

	// Only proceed if prices cross
	if highestBuyPrice.GreaterThanOrEqual(lowestSellPrice) {
		// Get all buy orders that cross
		for _, buyOrder := range buyOrders {
			if buyOrder.Price.GreaterThanOrEqual(lowestSellPrice) {
				crossingBuyOrders = append(crossingBuyOrders, buyOrder)
			} else {
				break
			}
		}

		// Get all sell orders that cross
		for _, sellOrder := range sellOrders {
			if sellOrder.Price.LessThanOrEqual(highestBuyPrice) {
				crossingSellOrders = append(crossingSellOrders, sellOrder)
			} else {
				break
			}
		}
	}

	return crossingBuyOrders, crossingSellOrders, nil
}

func (b *PostgresBot) executeMarketOnChain(ctx context.Context) (string, error) {
	// Get account info
	accountRes, err := b.authClient.Account(ctx, &authtypes.QueryAccountRequest{
		Address: b.address.String(),
	})
	if err != nil {
		return "", fmt.Errorf("failed to get account: %w", err)
	}

	var account authtypes.AccountI
	if err := b.clientCtx.InterfaceRegistry.UnpackAny(accountRes.Account, &account); err != nil {
		return "", fmt.Errorf("failed to unpack account: %w", err)
	}

	// Create match and execute orders message for this market
	msg := &clobtypes.MsgMatchAndExecuteOrders{
		Sender:    b.address.String(),
		MarketIds: []uint64{b.marketID},
		Limit:     100, // Process up to 100 orders
	}

	// Build transaction
	txBuilder := b.clientCtx.TxConfig.NewTxBuilder()
	if err := txBuilder.SetMsgs(msg); err != nil {
		return "", fmt.Errorf("failed to set messages: %w", err)
	}

	// Set gas and fees
	txBuilder.SetGasLimit(500000)
	fee := sdk.NewCoins(sdk.NewCoin("uelys", sdk.NewInt(5000)))
	txBuilder.SetFeeAmount(fee)

	// Sign transaction
	signerData := authclient.GetSignerData(b.clientCtx, account)
	sigV2, err := tx.SignWithPrivKey(
		signing.SignMode_SIGN_MODE_DIRECT,
		signerData,
		txBuilder,
		b.privKey,
		b.clientCtx.TxConfig,
		account.GetSequence(),
	)
	if err != nil {
		return "", fmt.Errorf("failed to sign transaction: %w", err)
	}

	if err := txBuilder.SetSignatures(sigV2); err != nil {
		return "", fmt.Errorf("failed to set signatures: %w", err)
	}

	// Encode transaction
	txBytes, err := b.clientCtx.TxConfig.TxEncoder()(txBuilder.GetTx())
	if err != nil {
		return "", fmt.Errorf("failed to encode transaction: %w", err)
	}

	// Broadcast transaction
	res, err := b.txClient.BroadcastTx(ctx, &tx.BroadcastTxRequest{
		TxBytes: txBytes,
		Mode:    tx.BroadcastMode_BROADCAST_MODE_SYNC,
	})
	if err != nil {
		return "", fmt.Errorf("failed to broadcast transaction: %w", err)
	}

	if res.TxResponse.Code != 0 {
		return "", fmt.Errorf("transaction failed: %s", res.TxResponse.RawLog)
	}

	log.Printf("Match transaction submitted: %s", res.TxResponse.TxHash)
	return res.TxResponse.TxHash, nil
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

func (b *PostgresBot) executeMatch(ctx context.Context, match *Match, txHash string) error {
	// Use transaction for consistency
	tx, err := b.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert match record with tx hash
	_, err = tx.ExecContext(ctx, `
		INSERT INTO clob_matches (buy_order_id, sell_order_id, market_id, price, amount, executed_at, tx_hash)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		match.BuyOrderID, match.SellOrderID, b.marketID,
		match.Price.String(), match.Amount.String(), match.ExecutedAt, txHash)
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
	if b.grpcConn != nil {
		b.grpcConn.Close()
	}
	return b.db.Close()
}

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://indexer:123123123@localhost:5433/elys_indexer?sslmode=disable"
	}

	grpcEndpoint := os.Getenv("GRPC_ENDPOINT")
	if grpcEndpoint == "" {
		grpcEndpoint = "localhost:9090"
	}

	chainID := os.Getenv("CHAIN_ID")
	if chainID == "" {
		chainID = "elys-1"
	}

	mnemonic := "choose isolate cruise nominee image peanut winter vacant enemy improve practice verb moon satisfy food fuel damage sugar load vendor mirror galaxy subject laptop"
	marketID := uint64(1)

	bot, err := NewPostgresBot(dbURL, mnemonic, marketID, grpcEndpoint, chainID)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}
	defer bot.Close()

	ctx := context.Background()
	bot.Start(ctx)
}
