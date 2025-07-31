package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/shopspring/decimal"
)

// Example CLOB Matching Bot
// This demonstrates how external bots can use the indexer API to perform order matching

type MatchingBot struct {
	indexerURL       string
	engineID         string
	httpClient       *http.Client
	matchingInterval time.Duration
}

type Order struct {
	OrderID         uint64 `json:"order_id"`
	MarketID        uint64 `json:"market_id"`
	OrderType       string `json:"order_type"`
	Price           string `json:"price"`
	RemainingAmount string `json:"remaining_amount"`
	Owner           string `json:"owner"`
}

type OrderBookResponse struct {
	MarketID   uint64   `json:"market_id"`
	BuyOrders  []*Order `json:"buy_orders"`
	SellOrders []*Order `json:"sell_orders"`
	Timestamp  string   `json:"timestamp"`
}

func NewMatchingBot(indexerURL, engineID string) *MatchingBot {
	return &MatchingBot{
		indexerURL:       indexerURL,
		engineID:         engineID,
		httpClient:       &http.Client{Timeout: 10 * time.Second},
		matchingInterval: 100 * time.Millisecond, // Check for matches every 100ms
	}
}

func (b *MatchingBot) Start(ctx context.Context, marketID uint64) {
	ticker := time.NewTicker(b.matchingInterval)
	defer ticker.Stop()

	log.Printf("Starting matching bot for market %d", marketID)

	for {
		select {
		case <-ctx.Done():
			log.Println("Matching bot stopped")
			return
		case <-ticker.C:
			if err := b.checkAndMatch(ctx, marketID); err != nil {
				log.Printf("Error in matching cycle: %v", err)
			}
		}
	}
}

func (b *MatchingBot) checkAndMatch(ctx context.Context, marketID uint64) error {
	// Get current order book
	orderBook, err := b.getOrderBook(marketID)
	if err != nil {
		return fmt.Errorf("failed to get order book: %w", err)
	}

	// Find matching opportunities
	matches := b.findMatches(orderBook.BuyOrders, orderBook.SellOrders)

	// Process each match
	for _, match := range matches {
		if err := b.processMatch(ctx, match); err != nil {
			log.Printf("Failed to process match: %v", err)
		}
	}

	return nil
}

func (b *MatchingBot) getOrderBook(marketID uint64) (*OrderBookResponse, error) {
	url := fmt.Sprintf("%s/api/v1/clob/matching/orders?market_id=%d", b.indexerURL, marketID)

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

type Match struct {
	BuyOrder    *Order
	SellOrder   *Order
	MatchPrice  decimal.Decimal
	MatchAmount decimal.Decimal
}

func (b *MatchingBot) findMatches(buyOrders, sellOrders []*Order) []Match {
	var matches []Match

	// Simple matching algorithm: match highest buy with lowest sell
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

			// Check if prices cross
			if buyPrice.GreaterThanOrEqual(sellPrice) {
				// Calculate match amount (minimum of both order amounts)
				matchAmount := buyAmount
				if sellAmount.LessThan(buyAmount) {
					matchAmount = sellAmount
				}

				// Use sell price as match price (price-time priority)
				matches = append(matches, Match{
					BuyOrder:    buyOrder,
					SellOrder:   sellOrder,
					MatchPrice:  sellPrice,
					MatchAmount: matchAmount,
				})

				// Update remaining amounts for next iteration
				buyAmount = buyAmount.Sub(matchAmount)
				sellAmount = sellAmount.Sub(matchAmount)

				if buyAmount.IsZero() {
					break // Buy order fully matched
				}
			}
		}
	}

	return matches
}

func (b *MatchingBot) processMatch(ctx context.Context, match Match) error {
	log.Printf("Processing match: Buy #%d @ %s <-> Sell #%d @ %s, Amount: %s",
		match.BuyOrder.OrderID, match.BuyOrder.Price,
		match.SellOrder.OrderID, match.SellOrder.Price,
		match.MatchAmount.String())

	// Step 1: Lock both orders
	buyLocked, err := b.lockOrder(match.BuyOrder.OrderID)
	if err != nil {
		return fmt.Errorf("failed to lock buy order: %w", err)
	}
	if !buyLocked {
		return fmt.Errorf("buy order already locked")
	}
	defer b.unlockOrder(match.BuyOrder.OrderID)

	sellLocked, err := b.lockOrder(match.SellOrder.OrderID)
	if err != nil {
		return fmt.Errorf("failed to lock sell order: %w", err)
	}
	if !sellLocked {
		return fmt.Errorf("sell order already locked")
	}
	defer b.unlockOrder(match.SellOrder.OrderID)

	// Step 2: Submit match to blockchain
	// In a real implementation, this would:
	// 1. Create a transaction with the match details
	// 2. Sign and broadcast the transaction
	// 3. Wait for confirmation

	log.Printf("Match would be submitted to blockchain: Price=%s, Amount=%s",
		match.MatchPrice.String(), match.MatchAmount.String())

	// Simulate blockchain submission
	time.Sleep(100 * time.Millisecond)

	return nil
}

func (b *MatchingBot) lockOrder(orderID uint64) (bool, error) {
	url := fmt.Sprintf("%s/api/v1/clob/matching/lock", b.indexerURL)

	reqBody := map[string]interface{}{
		"order_id":           orderID,
		"matching_engine_id": b.engineID,
		"ttl_seconds":        5,
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

func (b *MatchingBot) unlockOrder(orderID uint64) error {
	url := fmt.Sprintf("%s/api/v1/clob/matching/unlock", b.indexerURL)

	reqBody := map[string]interface{}{
		"order_id":           orderID,
		"matching_engine_id": b.engineID,
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

func main() {
	// Example usage
	indexerURL := "http://localhost:8080"
	engineID := "matching-bot-1"
	marketID := uint64(1)

	bot := NewMatchingBot(indexerURL, engineID)

	ctx := context.Background()
	bot.Start(ctx, marketID)
}
