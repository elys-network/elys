package main

import (
	"context"
	"encoding/json"
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
	"github.com/redis/go-redis/v9"
	"github.com/shopspring/decimal"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type RedisBot struct {
	client           *redis.Client
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
	OrderID         uint64          `json:"order_id"`
	MarketID        uint64          `json:"market_id"`
	Side            string          `json:"side"` // "buy" or "sell"
	Price           decimal.Decimal `json:"price"`
	RemainingAmount decimal.Decimal `json:"remaining_amount"`
	Owner           string          `json:"owner"`
	Timestamp       time.Time       `json:"timestamp"`
	Counter         uint64          `json:"counter"` // Order counter for OrderKey
}

type Match struct {
	BuyOrderID  uint64          `json:"buy_order_id"`
	SellOrderID uint64          `json:"sell_order_id"`
	Price       decimal.Decimal `json:"price"`
	Amount      decimal.Decimal `json:"amount"`
	Timestamp   time.Time       `json:"timestamp"`
}

func NewRedisBot(redisURL string, mnemonic string, marketID uint64, grpcEndpoint string, chainID string) (*RedisBot, error) {
	// Redis setup
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Redis URL: %w", err)
	}

	client := redis.NewClient(opt)
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
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

	return &RedisBot{
		client:           client,
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

func (b *RedisBot) Start(ctx context.Context) {
	log.Printf("Starting Redis bot for market %d", b.marketID)
	log.Printf("Bot address: %s", b.address.String())

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
	// Create some test orders with counters
	testOrders := []Order{
		{OrderID: 1, MarketID: b.marketID, Side: "buy", Price: decimal.NewFromInt(100), RemainingAmount: decimal.NewFromInt(10), Owner: "test1", Timestamp: time.Now(), Counter: 1},
		{OrderID: 2, MarketID: b.marketID, Side: "buy", Price: decimal.NewFromInt(99), RemainingAmount: decimal.NewFromInt(5), Owner: "test2", Timestamp: time.Now(), Counter: 2},
		{OrderID: 3, MarketID: b.marketID, Side: "sell", Price: decimal.NewFromInt(101), RemainingAmount: decimal.NewFromInt(8), Owner: "test3", Timestamp: time.Now(), Counter: 3},
		{OrderID: 4, MarketID: b.marketID, Side: "sell", Price: decimal.NewFromInt(102), RemainingAmount: decimal.NewFromInt(12), Owner: "test4", Timestamp: time.Now(), Counter: 4},
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
	// Get all crossing orders
	crossingBuyOrders, crossingSellOrders, err := b.findCrossingOrders(ctx)
	if err != nil {
		return err
	}

	if len(crossingBuyOrders) == 0 || len(crossingSellOrders) == 0 {
		return nil // No orders to match
	}

	// Execute match on blockchain by calling ExecuteMarket for this market
	if err := b.executeMarketOnChain(ctx); err != nil {
		log.Printf("Failed to execute market on chain: %v", err)
		// Continue with local recording even if blockchain tx fails
	}

	// Update local state based on matches
	for _, buyOrder := range crossingBuyOrders {
		for _, sellOrder := range crossingSellOrders {
			if buyOrder.Price.GreaterThanOrEqual(sellOrder.Price) {
				match := b.createMatch(&buyOrder, &sellOrder, sellOrder.Price)
				if match != nil {
					if err := b.executeMatch(ctx, match); err != nil {
						log.Printf("Failed to record match locally: %v", err)
					}

					log.Printf("Matched: Buy #%d @ %s with Sell #%d @ %s, Amount: %s",
						buyOrder.OrderID, buyOrder.Price.String(),
						sellOrder.OrderID, sellOrder.Price.String(),
						match.Amount.String())

					// Update remaining amounts
					buyOrder.RemainingAmount = buyOrder.RemainingAmount.Sub(match.Amount)
					sellOrder.RemainingAmount = sellOrder.RemainingAmount.Sub(match.Amount)

					if buyOrder.RemainingAmount.IsZero() {
						b.removeOrder(ctx, &buyOrder)
						break
					}
				}
			}
		}

		// Update or remove sell orders
		for _, sellOrder := range crossingSellOrders {
			if sellOrder.RemainingAmount.IsZero() {
				b.removeOrder(ctx, &sellOrder)
			} else {
				b.updateOrder(ctx, &sellOrder)
			}
		}

		// Update buy order if not fully filled
		if !buyOrder.RemainingAmount.IsZero() {
			b.updateOrder(ctx, &buyOrder)
		}
	}

	return nil
}

func (b *RedisBot) findCrossingOrders(ctx context.Context) ([]Order, []Order, error) {
	var crossingBuyOrders []Order
	var crossingSellOrders []Order

	// Get all buy orders sorted by price (highest first)
	buyBookKey := fmt.Sprintf("market:%d:buy", b.marketID)
	buyOrdersZ, err := b.client.ZRangeWithScores(ctx, buyBookKey, 0, -1).Result()
	if err != nil {
		return nil, nil, err
	}

	// Get all sell orders sorted by price (lowest first)
	sellBookKey := fmt.Sprintf("market:%d:sell", b.marketID)
	sellOrdersZ, err := b.client.ZRangeWithScores(ctx, sellBookKey, 0, -1).Result()
	if err != nil {
		return nil, nil, err
	}

	if len(buyOrdersZ) == 0 || len(sellOrdersZ) == 0 {
		return nil, nil, nil
	}

	// Get the best sell price (lowest)
	lowestSellOrderID := uint64(sellOrdersZ[0].Member.(int64))
	lowestSellOrder, err := b.getOrder(ctx, lowestSellOrderID)
	if err != nil {
		return nil, nil, err
	}
	lowestSellPrice := lowestSellOrder.Price

	// Find all buy orders that cross with the lowest sell price
	for _, buyZ := range buyOrdersZ {
		buyOrderID := uint64(buyZ.Member.(int64))
		buyOrder, err := b.getOrder(ctx, buyOrderID)
		if err != nil {
			continue
		}

		if buyOrder.Price.GreaterThanOrEqual(lowestSellPrice) {
			crossingBuyOrders = append(crossingBuyOrders, *buyOrder)
		} else {
			break // No more crossing buy orders
		}
	}

	if len(crossingBuyOrders) == 0 {
		return nil, nil, nil
	}

	// Get the highest buy price
	highestBuyPrice := crossingBuyOrders[0].Price

	// Find all sell orders that cross with the highest buy price
	for _, sellZ := range sellOrdersZ {
		sellOrderID := uint64(sellZ.Member.(int64))
		sellOrder, err := b.getOrder(ctx, sellOrderID)
		if err != nil {
			continue
		}

		if sellOrder.Price.LessThanOrEqual(highestBuyPrice) {
			crossingSellOrders = append(crossingSellOrders, *sellOrder)
		} else {
			break // No more crossing sell orders
		}
	}

	return crossingBuyOrders, crossingSellOrders, nil
}

func (b *RedisBot) executeMarketOnChain(ctx context.Context) error {
	// Get account info
	accountRes, err := b.authClient.Account(ctx, &authtypes.QueryAccountRequest{
		Address: b.address.String(),
	})
	if err != nil {
		return fmt.Errorf("failed to get account: %w", err)
	}

	var account authtypes.AccountI
	if err := b.clientCtx.InterfaceRegistry.UnpackAny(accountRes.Account, &account); err != nil {
		return fmt.Errorf("failed to unpack account: %w", err)
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
		return fmt.Errorf("failed to set messages: %w", err)
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
		return fmt.Errorf("failed to sign transaction: %w", err)
	}

	if err := txBuilder.SetSignatures(sigV2); err != nil {
		return fmt.Errorf("failed to set signatures: %w", err)
	}

	// Encode transaction
	txBytes, err := b.clientCtx.TxConfig.TxEncoder()(txBuilder.GetTx())
	if err != nil {
		return fmt.Errorf("failed to encode transaction: %w", err)
	}

	// Broadcast transaction
	res, err := b.txClient.BroadcastTx(ctx, &tx.BroadcastTxRequest{
		TxBytes: txBytes,
		Mode:    tx.BroadcastMode_BROADCAST_MODE_SYNC,
	})
	if err != nil {
		return fmt.Errorf("failed to broadcast transaction: %w", err)
	}

	if res.TxResponse.Code != 0 {
		return fmt.Errorf("transaction failed: %s", res.TxResponse.RawLog)
	}

	log.Printf("Match transaction submitted: %s", res.TxResponse.TxHash)
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
	if b.grpcConn != nil {
		b.grpcConn.Close()
	}
	return b.client.Close()
}

func main() {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "redis://localhost:6380"
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

	bot, err := NewRedisBot(redisURL, mnemonic, marketID, grpcEndpoint, chainID)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}
	defer bot.Close()

	ctx := context.Background()
	bot.Start(ctx)
}
