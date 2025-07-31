package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authclient "github.com/cosmos/cosmos-sdk/x/auth/client"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/go-bip39"
	clobtypes "github.com/elys-network/elys/v7/x/clob/types"
	"github.com/shopspring/decimal"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gopkg.in/yaml.v3"
)

// Config structure for the bot
type Config struct {
	Bot struct {
		EngineID         string        `yaml:"engine_id"`
		Mnemonic         string        `yaml:"mnemonic"`
		AccountNumber    uint64        `yaml:"account_number"`
		AccountSequence  uint64        `yaml:"account_sequence"`
		MatchingInterval time.Duration `yaml:"matching_interval"`
		MaxBatchSize     int           `yaml:"max_batch_size"`
		MinProfitBPS     int           `yaml:"min_profit_bps"` // Minimum profit in basis points
	} `yaml:"bot"`

	Chain struct {
		ChainID      string `yaml:"chain_id"`
		RPCEndpoint  string `yaml:"rpc_endpoint"`
		GRPCEndpoint string `yaml:"grpc_endpoint"`
		GasPrice     string `yaml:"gas_price"`
		GasLimit     uint64 `yaml:"gas_limit"`
		Denom        string `yaml:"denom"`
	} `yaml:"chain"`

	Indexer struct {
		URL string `yaml:"url"`
	} `yaml:"indexer"`

	Markets []uint64 `yaml:"markets"` // List of market IDs to monitor
}

// Enhanced Matching Bot with transaction capabilities
type MatchingBotV2 struct {
	config     *Config
	httpClient *http.Client
	grpcConn   *grpc.ClientConn
	txClient   tx.ServiceClient
	authClient authtypes.QueryClient
	clobClient clobtypes.QueryClient
	privKey    cryptotypes.PrivKey
	address    sdk.AccAddress
	clientCtx  client.Context
}

type Order struct {
	OrderID         uint64 `json:"order_id"`
	MarketID        uint64 `json:"market_id"`
	Counter         uint64 `json:"counter"`
	OrderType       string `json:"order_type"`
	Price           string `json:"price"`
	RemainingAmount string `json:"remaining_amount"`
	Owner           string `json:"owner"`
	SubAccountID    uint64 `json:"sub_account_id"`
}

type OrderBookResponse struct {
	MarketID   uint64   `json:"market_id"`
	BuyOrders  []*Order `json:"buy_orders"`
	SellOrders []*Order `json:"sell_orders"`
	Timestamp  string   `json:"timestamp"`
}

type Match struct {
	BuyOrder    *Order
	SellOrder   *Order
	MatchPrice  decimal.Decimal
	MatchAmount decimal.Decimal
	ProfitBPS   int // Profit in basis points
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &config, nil
}

func NewMatchingBotV2(config *Config) (*MatchingBotV2, error) {
	// Derive private key from mnemonic
	if !bip39.IsMnemonic(config.Bot.Mnemonic) {
		return nil, fmt.Errorf("invalid mnemonic")
	}

	seed, err := bip39.NewSeedWithErrorChecking(config.Bot.Mnemonic, "")
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
		config.Chain.GRPCEndpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to gRPC: %w", err)
	}

	// Create client context
	clientCtx := client.Context{}.
		WithChainID(config.Chain.ChainID).
		WithCodec(makeCodec()).
		WithFromAddress(address).
		WithFromName("bot").
		WithSkipConfirmation(true).
		WithTxConfig(makeTxConfig())

	return &MatchingBotV2{
		config:     config,
		httpClient: &http.Client{Timeout: 10 * time.Second},
		grpcConn:   grpcConn,
		txClient:   tx.NewServiceClient(grpcConn),
		authClient: authtypes.NewQueryClient(grpcConn),
		clobClient: clobtypes.NewQueryClient(grpcConn),
		privKey:    privKey,
		address:    address,
		clientCtx:  clientCtx,
	}, nil
}

func (b *MatchingBotV2) Start(ctx context.Context) {
	log.Printf("Starting matching bot v2 with address: %s", b.address.String())
	log.Printf("Monitoring markets: %v", b.config.Markets)

	ticker := time.NewTicker(b.config.Bot.MatchingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Matching bot stopped")
			return
		case <-ticker.C:
			for _, marketID := range b.config.Markets {
				if err := b.checkAndMatch(ctx, marketID); err != nil {
					log.Printf("Error in matching cycle for market %d: %v", marketID, err)
				}
			}
		}
	}
}

func (b *MatchingBotV2) checkAndMatch(ctx context.Context, marketID uint64) error {
	// Get current order book
	orderBook, err := b.getOrderBook(marketID)
	if err != nil {
		return fmt.Errorf("failed to get order book: %w", err)
	}

	// Find profitable matches
	matches := b.findProfitableMatches(orderBook.BuyOrders, orderBook.SellOrders)

	if len(matches) == 0 {
		return nil
	}

	// Batch matches up to max batch size
	batches := b.createBatches(matches)

	// Process each batch
	for i, batch := range batches {
		log.Printf("Processing batch %d/%d with %d matches", i+1, len(batches), len(batch))
		if err := b.processBatch(ctx, batch); err != nil {
			log.Printf("Failed to process batch: %v", err)
		}
	}

	return nil
}

func (b *MatchingBotV2) getOrderBook(marketID uint64) (*OrderBookResponse, error) {
	url := fmt.Sprintf("%s/api/v1/clob/matching/orders?market_id=%d", b.config.Indexer.URL, marketID)

	resp, err := b.httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: %s", string(body))
	}

	var orderBook OrderBookResponse
	if err := json.NewDecoder(resp.Body).Decode(&orderBook); err != nil {
		return nil, err
	}

	return &orderBook, nil
}

func (b *MatchingBotV2) findProfitableMatches(buyOrders, sellOrders []*Order) []Match {
	var matches []Match

	for _, buyOrder := range buyOrders {
		buyPrice, err := decimal.NewFromString(buyOrder.Price)
		if err != nil {
			continue
		}
		buyAmount, err := decimal.NewFromString(buyOrder.RemainingAmount)
		if err != nil {
			continue
		}

		for _, sellOrder := range sellOrders {
			sellPrice, err := decimal.NewFromString(sellOrder.Price)
			if err != nil {
				continue
			}
			sellAmount, err := decimal.NewFromString(sellOrder.RemainingAmount)
			if err != nil {
				continue
			}

			// Check if prices cross with profit margin
			if buyPrice.GreaterThan(sellPrice) {
				// Calculate profit in basis points
				spread := buyPrice.Sub(sellPrice)
				profitBPS := spread.Mul(decimal.NewFromInt(10000)).Div(sellPrice).IntPart()

				if profitBPS >= int64(b.config.Bot.MinProfitBPS) {
					// Calculate match amount
					matchAmount := buyAmount
					if sellAmount.LessThan(buyAmount) {
						matchAmount = sellAmount
					}

					matches = append(matches, Match{
						BuyOrder:    buyOrder,
						SellOrder:   sellOrder,
						MatchPrice:  sellPrice, // Use sell price for taker
						MatchAmount: matchAmount,
						ProfitBPS:   int(profitBPS),
					})

					// Update remaining amounts
					buyAmount = buyAmount.Sub(matchAmount)
					sellAmount = sellAmount.Sub(matchAmount)

					if buyAmount.IsZero() {
						break
					}
				}
			}
		}
	}

	return matches
}

func (b *MatchingBotV2) createBatches(matches []Match) [][]Match {
	var batches [][]Match

	for i := 0; i < len(matches); i += b.config.Bot.MaxBatchSize {
		end := i + b.config.Bot.MaxBatchSize
		if end > len(matches) {
			end = len(matches)
		}
		batches = append(batches, matches[i:end])
	}

	return batches
}

func (b *MatchingBotV2) processBatch(ctx context.Context, matches []Match) error {
	// Lock all orders in the batch
	lockedOrders := make(map[uint64]bool)

	for _, match := range matches {
		// Lock buy order
		if !lockedOrders[match.BuyOrder.OrderID] {
			locked, err := b.lockOrder(match.BuyOrder.OrderID)
			if err != nil || !locked {
				log.Printf("Failed to lock buy order %d", match.BuyOrder.OrderID)
				continue
			}
			lockedOrders[match.BuyOrder.OrderID] = true
			defer b.unlockOrder(match.BuyOrder.OrderID)
		}

		// Lock sell order
		if !lockedOrders[match.SellOrder.OrderID] {
			locked, err := b.lockOrder(match.SellOrder.OrderID)
			if err != nil || !locked {
				log.Printf("Failed to lock sell order %d", match.SellOrder.OrderID)
				continue
			}
			lockedOrders[match.SellOrder.OrderID] = true
			defer b.unlockOrder(match.SellOrder.OrderID)
		}
	}

	// Create batch match message
	msg, err := b.createBatchMatchMsg(matches)
	if err != nil {
		return fmt.Errorf("failed to create batch match message: %w", err)
	}

	// Sign and broadcast transaction
	txHash, err := b.broadcastTx(ctx, msg)
	if err != nil {
		return fmt.Errorf("failed to broadcast transaction: %w", err)
	}

	log.Printf("Batch match transaction submitted: %s", txHash)

	// Log match details
	for _, match := range matches {
		log.Printf("Matched: Buy #%d @ %s <-> Sell #%d @ %s, Amount: %s, Profit: %d bps",
			match.BuyOrder.OrderID, match.BuyOrder.Price,
			match.SellOrder.OrderID, match.SellOrder.Price,
			match.MatchAmount.String(), match.ProfitBPS)
	}

	return nil
}

func (b *MatchingBotV2) createBatchMatchMsg(matches []Match) (sdk.Msg, error) {
	// Create match operations
	var operations []*clobtypes.MatchOperation

	for _, match := range matches {
		buyOrderKey := &clobtypes.OrderKey{
			MarketId:  match.BuyOrder.MarketID,
			OrderType: clobtypes.OrderType_ORDER_TYPE_LIMIT_BUY,
			Price:     sdk.MustNewDecFromStr(match.BuyOrder.Price),
			Counter:   match.BuyOrder.Counter,
		}

		sellOrderKey := &clobtypes.OrderKey{
			MarketId:  match.SellOrder.MarketID,
			OrderType: clobtypes.OrderType_ORDER_TYPE_LIMIT_SELL,
			Price:     sdk.MustNewDecFromStr(match.SellOrder.Price),
			Counter:   match.SellOrder.Counter,
		}

		operation := &clobtypes.MatchOperation{
			BuyOrderKey:  buyOrderKey,
			SellOrderKey: sellOrderKey,
			Quantity:     sdk.MustNewDecFromStr(match.MatchAmount.String()),
			Price:        sdk.MustNewDecFromStr(match.MatchPrice.String()),
		}

		operations = append(operations, operation)
	}

	// Create batch match message
	msg := &clobtypes.MsgBatchMatch{
		Creator:    b.address.String(),
		Operations: operations,
	}

	return msg, nil
}

func (b *MatchingBotV2) broadcastTx(ctx context.Context, msgs ...sdk.Msg) (string, error) {
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

	// Build transaction
	txBuilder := b.clientCtx.TxConfig.NewTxBuilder()
	if err := txBuilder.SetMsgs(msgs...); err != nil {
		return "", fmt.Errorf("failed to set messages: %w", err)
	}

	// Set gas and fees
	txBuilder.SetGasLimit(b.config.Chain.GasLimit)
	gasPrice, err := sdk.ParseDecCoin(b.config.Chain.GasPrice)
	if err != nil {
		return "", fmt.Errorf("failed to parse gas price: %w", err)
	}
	fee := sdk.NewCoins(sdk.NewCoin(gasPrice.Denom, gasPrice.Amount.MulInt64(int64(b.config.Chain.GasLimit)).TruncateInt()))
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

	return res.TxResponse.TxHash, nil
}

func (b *MatchingBotV2) lockOrder(orderID uint64) (bool, error) {
	url := fmt.Sprintf("%s/api/v1/clob/matching/lock", b.config.Indexer.URL)

	reqBody := map[string]interface{}{
		"order_id":           orderID,
		"matching_engine_id": b.config.Bot.EngineID,
		"ttl_seconds":        10,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return false, err
	}

	resp, err := b.httpClient.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("failed to lock order: status %d", resp.StatusCode)
	}

	var result struct {
		Locked bool `json:"locked"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, err
	}

	return result.Locked, nil
}

func (b *MatchingBotV2) unlockOrder(orderID uint64) error {
	url := fmt.Sprintf("%s/api/v1/clob/matching/unlock", b.config.Indexer.URL)

	reqBody := map[string]interface{}{
		"order_id":           orderID,
		"matching_engine_id": b.config.Bot.EngineID,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	resp, err := b.httpClient.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to unlock order: status %d", resp.StatusCode)
	}

	return nil
}

// Helper functions to create codec and tx config
func makeCodec() *codec.ProtoCodec {
	// This should include all the necessary types
	// In production, use the app's codec
	return codec.NewProtoCodec(makeInterfaceRegistry())
}

func makeInterfaceRegistry() types.InterfaceRegistry {
	interfaceRegistry := types.NewInterfaceRegistry()
	// Register interfaces - in production, use app's interface registry
	authtypes.RegisterInterfaces(interfaceRegistry)
	clobtypes.RegisterInterfaces(interfaceRegistry)
	return interfaceRegistry
}

func makeTxConfig() client.TxConfig {
	return tx.NewTxConfig(makeCodec(), tx.DefaultSignModes)
}

func main() {
	// Load configuration
	configPath := "bot_config.yaml"
	if len(os.Args) > 1 {
		configPath = os.Args[1]
	}

	config, err := LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Create bot
	bot, err := NewMatchingBotV2(config)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}
	defer bot.grpcConn.Close()

	// Start bot
	ctx := context.Background()
	bot.Start(ctx)
}
