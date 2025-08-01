package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/elys-network/elys/indexer/internal/models"
	"github.com/redis/go-redis/v9"
	"github.com/shopspring/decimal"
)

// CLOB Order Book Management
func (c *Cache) AddCLOBOrderToBook(ctx context.Context, order *models.CLOBOrder) error {
	// Add to order book sorted set
	var key string
	var score float64

	if order.OrderType == models.CLOBOrderTypeLimitBuy || order.OrderType == models.CLOBOrderTypeMarketBuy {
		key = fmt.Sprintf("clob:orderbook:%d:bids", order.MarketID)
		// For bids, higher price has higher priority (use negative score for reverse order)
		score = order.Price.InexactFloat64()*1e18 - float64(order.Counter)
	} else {
		key = fmt.Sprintf("clob:orderbook:%d:asks", order.MarketID)
		// For asks, lower price has higher priority
		score = order.Price.InexactFloat64()*1e18 + float64(order.Counter)
	}

	// Add to sorted set
	if err := c.client.ZAdd(ctx, key, redis.Z{
		Score:  score,
		Member: order.OrderID,
	}).Err(); err != nil {
		return err
	}

	// Store order details
	if err := c.SetCLOBOrder(ctx, order); err != nil {
		return err
	}

	// Add to user's active orders
	userOrdersKey := fmt.Sprintf("clob:user:%s:orders", order.Owner)
	if err := c.client.SAdd(ctx, userOrdersKey, order.OrderID).Err(); err != nil {
		return err
	}

	// Add to user's market orders
	marketOrdersKey := fmt.Sprintf("clob:user:%s:market:%d:orders", order.Owner, order.MarketID)
	if err := c.client.SAdd(ctx, marketOrdersKey, order.OrderID).Err(); err != nil {
		return err
	}

	// Update best bid/ask if necessary
	return c.updateBestPrices(ctx, order.MarketID)
}

func (c *Cache) RemoveCLOBOrderFromBook(ctx context.Context, order *models.CLOBOrder) error {
	// Remove from order book
	var key string
	if order.OrderType == models.CLOBOrderTypeLimitBuy || order.OrderType == models.CLOBOrderTypeMarketBuy {
		key = fmt.Sprintf("clob:orderbook:%d:bids", order.MarketID)
	} else {
		key = fmt.Sprintf("clob:orderbook:%d:asks", order.MarketID)
	}

	if err := c.client.ZRem(ctx, key, order.OrderID).Err(); err != nil {
		return err
	}

	// Remove from user's active orders
	userOrdersKey := fmt.Sprintf("clob:user:%s:orders", order.Owner)
	if err := c.client.SRem(ctx, userOrdersKey, order.OrderID).Err(); err != nil {
		return err
	}

	// Remove from user's market orders
	marketOrdersKey := fmt.Sprintf("clob:user:%s:market:%d:orders", order.Owner, order.MarketID)
	if err := c.client.SRem(ctx, marketOrdersKey, order.OrderID).Err(); err != nil {
		return err
	}

	// Update best prices
	return c.updateBestPrices(ctx, order.MarketID)
}

func (c *Cache) SetCLOBOrder(ctx context.Context, order *models.CLOBOrder) error {
	key := fmt.Sprintf("clob:order:%d", order.OrderID)

	data := map[string]interface{}{
		"market_id":        order.MarketID,
		"owner":            order.Owner,
		"sub_account_id":   order.SubAccountID,
		"order_type":       order.OrderType,
		"price":            order.Price.String(),
		"amount":           order.Amount.String(),
		"filled_amount":    order.FilledAmount.String(),
		"remaining_amount": order.RemainingAmount.String(),
		"status":           order.Status,
		"created_at":       order.CreatedAt.Unix(),
		"block_height":     order.BlockHeight,
	}

	return c.client.HMSet(ctx, key, data).Err()
}

func (c *Cache) GetCLOBOrder(ctx context.Context, orderID uint64) (*models.CLOBOrder, error) {
	key := fmt.Sprintf("clob:order:%d", orderID)

	data, err := c.client.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return nil, nil
	}

	// Parse the order from Redis hash
	order := &models.CLOBOrder{
		OrderID: orderID,
	}

	// Parse fields (simplified - full implementation would handle all fields)
	if marketID, ok := data["market_id"]; ok {
		fmt.Sscanf(marketID, "%d", &order.MarketID)
	}
	order.Owner = data["owner"]
	order.OrderType = models.CLOBOrderType(data["order_type"])
	order.Status = models.CLOBOrderStatus(data["status"])

	if price, ok := data["price"]; ok {
		order.Price, _ = decimal.NewFromString(price)
	}
	if amount, ok := data["amount"]; ok {
		order.Amount, _ = decimal.NewFromString(amount)
	}

	return order, nil
}

// CLOB Positions
func (c *Cache) SetCLOBPosition(ctx context.Context, position *models.CLOBPosition) error {
	// Add to positions by market
	marketPosKey := fmt.Sprintf("clob:positions:market:%d", position.MarketID)
	if err := c.client.SAdd(ctx, marketPosKey, position.PositionID).Err(); err != nil {
		return err
	}

	// Add to user's positions
	userPosKey := fmt.Sprintf("clob:user:%s:positions", position.Owner)
	if err := c.client.SAdd(ctx, userPosKey, position.PositionID).Err(); err != nil {
		return err
	}

	// Store position details
	posKey := fmt.Sprintf("clob:position:%d", position.PositionID)
	data := map[string]interface{}{
		"market_id":         position.MarketID,
		"owner":             position.Owner,
		"sub_account_id":    position.SubAccountID,
		"side":              position.Side,
		"size":              position.Size.String(),
		"entry_price":       position.EntryPrice.String(),
		"mark_price":        position.MarkPrice.String(),
		"liquidation_price": position.LiquidationPrice.String(),
		"margin":            position.Margin.String(),
		"unrealized_pnl":    position.UnrealizedPnL.String(),
		"updated_at":        position.UpdatedAt.Unix(),
	}

	return c.client.HMSet(ctx, posKey, data).Err()
}

func (c *Cache) RemoveCLOBPosition(ctx context.Context, position *models.CLOBPosition) error {
	// Remove from market positions
	marketPosKey := fmt.Sprintf("clob:positions:market:%d", position.MarketID)
	if err := c.client.SRem(ctx, marketPosKey, position.PositionID).Err(); err != nil {
		return err
	}

	// Remove from user's positions
	userPosKey := fmt.Sprintf("clob:user:%s:positions", position.Owner)
	if err := c.client.SRem(ctx, userPosKey, position.PositionID).Err(); err != nil {
		return err
	}

	// Delete position data
	posKey := fmt.Sprintf("clob:position:%d", position.PositionID)
	return c.client.Del(ctx, posKey).Err()
}

// CLOB Market Data
func (c *Cache) SetCLOBMarketData(ctx context.Context, marketID uint64, data map[string]interface{}) error {
	key := fmt.Sprintf("clob:market:%d", marketID)
	return c.client.HMSet(ctx, key, data).Err()
}

func (c *Cache) GetCLOBMarketData(ctx context.Context, marketID uint64) (map[string]string, error) {
	key := fmt.Sprintf("clob:market:%d", marketID)
	return c.client.HGetAll(ctx, key).Result()
}

func (c *Cache) SetCLOBTicker(ctx context.Context, marketID uint64, ticker *models.CLOBMarketStats) error {
	tickerData := map[string]interface{}{
		"last_price":    ticker.LastPrice.String(),
		"volume_24h":    ticker.Volume24h.String(),
		"trades_24h":    ticker.Trades24h,
		"high_24h":      ticker.High24h.String(),
		"low_24h":       ticker.Low24h.String(),
		"open_interest": ticker.OpenInterest.String(),
		"updated_at":    time.Now().Unix(),
	}

	data, err := json.Marshal(tickerData)
	if err != nil {
		return err
	}

	return c.client.HSet(ctx, "clob:tickers", fmt.Sprintf("%d", marketID), data).Err()
}

// CLOB Trades
func (c *Cache) AddCLOBTrade(ctx context.Context, trade *models.CLOBTrade) error {
	// Add to market trades
	marketTradesKey := fmt.Sprintf("clob:trades:market:%d", trade.MarketID)
	tradeJSON, err := json.Marshal(trade)
	if err != nil {
		return err
	}

	// Use LPUSH and LTRIM to maintain a capped list
	pipe := c.client.Pipeline()
	pipe.LPush(ctx, marketTradesKey, tradeJSON)
	pipe.LTrim(ctx, marketTradesKey, 0, 999) // Keep last 1000 trades

	// Add to user trades
	buyerTradesKey := fmt.Sprintf("clob:user:%s:trades", trade.Buyer)
	pipe.LPush(ctx, buyerTradesKey, tradeJSON)
	pipe.LTrim(ctx, buyerTradesKey, 0, 999)

	if trade.Buyer != trade.Seller {
		sellerTradesKey := fmt.Sprintf("clob:user:%s:trades", trade.Seller)
		pipe.LPush(ctx, sellerTradesKey, tradeJSON)
		pipe.LTrim(ctx, sellerTradesKey, 0, 999)
	}

	_, err = pipe.Exec(ctx)
	return err
}

// Order Matching Support
func (c *Cache) LockCLOBOrder(ctx context.Context, orderID uint64, matchingEngineID string, ttl time.Duration) (bool, error) {
	key := fmt.Sprintf("clob:order:lock:%d", orderID)
	return c.client.SetNX(ctx, key, matchingEngineID, ttl).Result()
}

func (c *Cache) UnlockCLOBOrder(ctx context.Context, orderID uint64, matchingEngineID string) error {
	key := fmt.Sprintf("clob:order:lock:%d", orderID)

	// Use Lua script to ensure we only unlock if we hold the lock
	script := `
		if redis.call("get", KEYS[1]) == ARGV[1] then
			return redis.call("del", KEYS[1])
		else
			return 0
		end
	`

	return c.client.Eval(ctx, script, []string{key}, matchingEngineID).Err()
}

func (c *Cache) GetBestBidAsk(ctx context.Context, marketID uint64) (bestBid, bestAsk *decimal.Decimal, err error) {
	bidKey := fmt.Sprintf("clob:market:%d:best_bid", marketID)
	askKey := fmt.Sprintf("clob:market:%d:best_ask", marketID)

	pipe := c.client.Pipeline()
	bidCmd := pipe.Get(ctx, bidKey)
	askCmd := pipe.Get(ctx, askKey)
	_, err = pipe.Exec(ctx)

	if err != nil && err != redis.Nil {
		return nil, nil, err
	}

	if bidStr, err := bidCmd.Result(); err == nil && bidStr != "" {
		if bid, err := decimal.NewFromString(bidStr); err == nil {
			bestBid = &bid
		}
	}

	if askStr, err := askCmd.Result(); err == nil && askStr != "" {
		if ask, err := decimal.NewFromString(askStr); err == nil {
			bestAsk = &ask
		}
	}

	return bestBid, bestAsk, nil
}

func (c *Cache) updateBestPrices(ctx context.Context, marketID uint64) error {
	// Get best bid
	bidKey := fmt.Sprintf("clob:orderbook:%d:bids", marketID)
	bestBids, err := c.client.ZRevRangeWithScores(ctx, bidKey, 0, 0).Result()
	if err != nil {
		return err
	}

	// Get best ask
	askKey := fmt.Sprintf("clob:orderbook:%d:asks", marketID)
	bestAsks, err := c.client.ZRangeWithScores(ctx, askKey, 0, 0).Result()
	if err != nil {
		return err
	}

	pipe := c.client.Pipeline()

	// Update best bid
	if len(bestBids) > 0 {
		orderID := bestBids[0].Member.(string)
		order, err := c.GetCLOBOrder(ctx, parseUint64(orderID))
		if err == nil && order != nil {
			bestBidKey := fmt.Sprintf("clob:market:%d:best_bid", marketID)
			pipe.Set(ctx, bestBidKey, order.Price.String(), 0)
		}
	} else {
		// No bids, delete the key
		bestBidKey := fmt.Sprintf("clob:market:%d:best_bid", marketID)
		pipe.Del(ctx, bestBidKey)
	}

	// Update best ask
	if len(bestAsks) > 0 {
		orderID := bestAsks[0].Member.(string)
		order, err := c.GetCLOBOrder(ctx, parseUint64(orderID))
		if err == nil && order != nil {
			bestAskKey := fmt.Sprintf("clob:market:%d:best_ask", marketID)
			pipe.Set(ctx, bestAskKey, order.Price.String(), 0)
		}
	} else {
		// No asks, delete the key
		bestAskKey := fmt.Sprintf("clob:market:%d:best_ask", marketID)
		pipe.Del(ctx, bestAskKey)
	}

	_, err = pipe.Exec(ctx)
	return err
}

// WebSocket Updates
func (c *Cache) PublishCLOBOrderUpdate(ctx context.Context, owner string, update *models.WSCLOBOrderUpdate) error {
	channel := fmt.Sprintf("clob:updates:orders:%s", owner)
	data, err := json.Marshal(update)
	if err != nil {
		return err
	}

	// Publish to user channel
	if err := c.client.Publish(ctx, channel, data).Err(); err != nil {
		return err
	}

	// Also add to stream for persistence
	return c.client.XAdd(ctx, &redis.XAddArgs{
		Stream: "clob:stream:orders",
		Values: map[string]interface{}{
			"owner":     owner,
			"action":    update.Action,
			"order_id":  update.Order.OrderID,
			"market_id": update.Order.MarketID,
			"data":      string(data),
		},
		MaxLen: 10000, // Keep last 10k events
		Approx: true,
	}).Err()
}

func (c *Cache) PublishCLOBTradeUpdate(ctx context.Context, marketID uint64, update *models.WSCLOBTradeUpdate) error {
	channel := fmt.Sprintf("clob:updates:trades:%d", marketID)
	data, err := json.Marshal(update)
	if err != nil {
		return err
	}

	// Publish to market channel
	if err := c.client.Publish(ctx, channel, data).Err(); err != nil {
		return err
	}

	// Also add to stream
	return c.client.XAdd(ctx, &redis.XAddArgs{
		Stream: "clob:stream:trades",
		Values: map[string]interface{}{
			"market_id": marketID,
			"trade_id":  update.Trade.TradeID,
			"price":     update.Trade.Price.String(),
			"quantity":  update.Trade.Quantity.String(),
			"data":      string(data),
		},
		MaxLen: 10000,
		Approx: true,
	}).Err()
}

func (c *Cache) PublishCLOBPositionUpdate(ctx context.Context, owner string, update *models.WSCLOBPositionUpdate) error {
	channel := fmt.Sprintf("clob:updates:positions:%s", owner)
	data, err := json.Marshal(update)
	if err != nil {
		return err
	}

	return c.client.Publish(ctx, channel, data).Err()
}

func (c *Cache) PublishCLOBOrderBookUpdate(ctx context.Context, marketID uint64, update *models.WSCLOBOrderBookUpdate) error {
	channel := fmt.Sprintf("clob:updates:orderbook:%d", marketID)
	data, err := json.Marshal(update)
	if err != nil {
		return err
	}

	return c.client.Publish(ctx, channel, data).Err()
}

// Helper function
func parseUint64(s string) uint64 {
	var v uint64
	fmt.Sscanf(s, "%d", &v)
	return v
}
