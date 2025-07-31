package indexer

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"sort"
	"time"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/rpc/client/http"
	tmtypes "github.com/cometbft/cometbft/types"
	"github.com/elys-network/elys/indexer/internal/cache"
	"github.com/elys-network/elys/indexer/internal/config"
	"github.com/elys-network/elys/indexer/internal/database"
	"github.com/elys-network/elys/indexer/internal/events"
	"github.com/elys-network/elys/indexer/internal/models"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Indexer struct {
	config            *config.Config
	db                *database.Repository
	cache             *cache.Cache
	rpcClient         *http.HTTP
	grpcConn          *grpc.ClientConn
	tradeshieldParser *events.TradeShieldParser
	perpetualParser   *events.PerpetualParser
	clobParser        *events.CLOBParser
	logger            *zap.Logger
	eventChan         chan *EventBatch
	stopChan          chan struct{}
}

type EventBatch struct {
	Height     int64
	Events     []interface{}
	OrderBooks []*models.OrderBookSnapshot
}

func New(cfg *config.Config, db *database.Repository, cache *cache.Cache, logger *zap.Logger) (*Indexer, error) {
	// Create RPC client
	rpcClient, err := http.New(cfg.Chain.RPCEndpoint, "/websocket")
	if err != nil {
		return nil, fmt.Errorf("failed to create RPC client: %w", err)
	}

	// Create gRPC connection
	grpcConn, err := grpc.Dial(cfg.Chain.GRPCEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection: %w", err)
	}

	return &Indexer{
		config:            cfg,
		db:                db,
		cache:             cache,
		rpcClient:         rpcClient,
		grpcConn:          grpcConn,
		tradeshieldParser: events.NewTradeShieldParser(logger),
		perpetualParser:   events.NewPerpetualParser(logger),
		clobParser:        events.NewCLOBParser(logger),
		logger:            logger,
		eventChan:         make(chan *EventBatch, cfg.Indexer.EventBufferSize),
		stopChan:          make(chan struct{}),
	}, nil
}

func (i *Indexer) Start(ctx context.Context) error {
	// Start RPC client
	if err := i.rpcClient.Start(); err != nil {
		return fmt.Errorf("failed to start RPC client: %w", err)
	}

	// Get current indexer state
	state, err := i.db.GetIndexerState(ctx)
	if err != nil {
		return fmt.Errorf("failed to get indexer state: %w", err)
	}

	startHeight := state.LastProcessedHeight + 1
	if i.config.Chain.StartHeight > 0 && state.LastProcessedHeight == 0 {
		startHeight = i.config.Chain.StartHeight
	}

	i.logger.Info("Starting indexer",
		zap.Int64("start_height", startHeight),
		zap.Int64("last_processed", state.LastProcessedHeight),
	)

	// Start worker pools
	for j := 0; j < i.config.Indexer.WorkerCount; j++ {
		go i.eventProcessor(ctx)
	}

	// Start order book aggregator
	go i.orderBookAggregator(ctx)

	// Start block processor
	go i.blockProcessor(ctx, startHeight)

	return nil
}

func (i *Indexer) Stop() error {
	close(i.stopChan)
	i.rpcClient.Stop()
	i.grpcConn.Close()
	return nil
}

func (i *Indexer) blockProcessor(ctx context.Context, startHeight int64) {
	height := startHeight
	ticker := time.NewTicker(i.config.Chain.PollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-i.stopChan:
			return
		case <-ticker.C:
			// Get latest block height
			status, err := i.rpcClient.Status(ctx)
			if err != nil {
				i.logger.Error("Failed to get chain status", zap.Error(err))
				continue
			}

			latestHeight := status.SyncInfo.LatestBlockHeight

			// Process blocks in batches
			for height <= latestHeight {
				batchEnd := height + int64(i.config.Chain.BatchSize)
				if batchEnd > latestHeight {
					batchEnd = latestHeight
				}

				if err := i.processBatch(ctx, height, batchEnd); err != nil {
					i.logger.Error("Failed to process batch",
						zap.Int64("start", height),
						zap.Int64("end", batchEnd),
						zap.Error(err),
					)
					// Retry after delay
					time.Sleep(i.config.Chain.RetryDelay)
					continue
				}

				height = batchEnd + 1
			}
		}
	}
}

func (i *Indexer) processBatch(ctx context.Context, startHeight, endHeight int64) error {
	for height := startHeight; height <= endHeight; height++ {
		if err := i.processBlock(ctx, height); err != nil {
			return fmt.Errorf("failed to process block %d: %w", height, err)
		}
	}
	return nil
}

func (i *Indexer) processBlock(ctx context.Context, height int64) error {
	// Get block results
	blockResults, err := i.rpcClient.BlockResults(ctx, &height)
	if err != nil {
		return fmt.Errorf("failed to get block results: %w", err)
	}

	// Get block for timestamp
	block, err := i.rpcClient.Block(ctx, &height)
	if err != nil {
		return fmt.Errorf("failed to get block: %w", err)
	}

	batch := &EventBatch{
		Height: height,
		Events: make([]interface{}, 0),
	}

	// Process transaction results
	for txIdx, txResult := range blockResults.TxsResults {
		if txResult.Code != 0 {
			// Skip failed transactions
			continue
		}

		txHash := getTxHash(block.Block.Txs[txIdx])

		// Parse TradeSheild events
		tradeshieldEvents, err := i.tradeshieldParser.ParseEvents(ctx, txResult.Events, height, txHash)
		if err != nil {
			i.logger.Error("Failed to parse tradeshield events",
				zap.Int64("height", height),
				zap.Int("tx_index", txIdx),
				zap.Error(err),
			)
		}
		batch.Events = append(batch.Events, tradeshieldEvents...)

		// Parse Perpetual events
		perpetualEvents, err := i.perpetualParser.ParseEvents(ctx, txResult.Events, height, txHash)
		if err != nil {
			i.logger.Error("Failed to parse perpetual events",
				zap.Int64("height", height),
				zap.Int("tx_index", txIdx),
				zap.Error(err),
			)
		}
		batch.Events = append(batch.Events, perpetualEvents...)

		// Parse CLOB events
		clobEvents, err := i.clobParser.ParseEvents(ctx, txResult.Events, height, txHash)
		if err != nil {
			i.logger.Error("Failed to parse CLOB events",
				zap.Int64("height", height),
				zap.Int("tx_index", txIdx),
				zap.Error(err),
			)
		}
		batch.Events = append(batch.Events, clobEvents...)
	}

	// Also process end block events (if available)
	if blockResults.FinalizeBlockEvents != nil {
		endBlockEvents, err := i.processEndBlockEvents(ctx, blockResults.FinalizeBlockEvents, height, block.Block.Time)
		if err != nil {
			i.logger.Error("Failed to process end block events", zap.Error(err))
		}
		batch.Events = append(batch.Events, endBlockEvents...)
	}

	// Send batch to event processor
	select {
	case i.eventChan <- batch:
	case <-ctx.Done():
		return ctx.Err()
	}

	return nil
}

func (i *Indexer) processEndBlockEvents(ctx context.Context, events []abci.Event, height int64, blockTime time.Time) ([]interface{}, error) {
	// Process any automated executions or liquidations that happen in end block
	var results []interface{}

	// Parse events
	tradeshieldEvents, err := i.tradeshieldParser.ParseEvents(ctx, events, height, "endblock")
	if err != nil {
		return nil, err
	}
	results = append(results, tradeshieldEvents...)

	perpetualEvents, err := i.perpetualParser.ParseEvents(ctx, events, height, "endblock")
	if err != nil {
		return nil, err
	}
	results = append(results, perpetualEvents...)

	clobEvents, err := i.clobParser.ParseEvents(ctx, events, height, "endblock")
	if err != nil {
		return nil, err
	}
	results = append(results, clobEvents...)

	return results, nil
}

func (i *Indexer) eventProcessor(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case batch := <-i.eventChan:
			if err := i.processBatchEvents(ctx, batch); err != nil {
				i.logger.Error("Failed to process batch events",
					zap.Int64("height", batch.Height),
					zap.Error(err),
				)
			}
		}
	}
}

func (i *Indexer) processBatchEvents(ctx context.Context, batch *EventBatch) error {
	// Process all events
	for _, event := range batch.Events {
		if err := i.processEvent(ctx, nil, event); err != nil {
			i.logger.Error("Failed to process event", zap.Error(err))
			// Continue processing other events
		}
	}

	// Update indexer state
	if err := i.db.UpdateIndexerState(ctx, batch.Height); err != nil {
		return fmt.Errorf("failed to update indexer state: %w", err)
	}

	return nil
}

func (i *Indexer) processEvent(ctx context.Context, tx *sql.Tx, event interface{}) error {
	switch e := event.(type) {
	case *models.SpotOrder:
		if err := i.db.CreateSpotOrder(ctx, e); err != nil {
			return err
		}
		// Cache and publish update
		i.cache.SetSpotOrder(ctx, e)
		i.publishOrderUpdate(ctx, e, "created")

	case *models.PerpetualOrder:
		if err := i.db.CreatePerpetualOrder(ctx, e); err != nil {
			return err
		}
		// Cache and publish update
		i.cache.SetPerpetualOrder(ctx, e)
		i.publishOrderUpdate(ctx, e, "created")

	case *models.PerpetualPosition:
		if err := i.db.CreatePerpetualPosition(ctx, e); err != nil {
			return err
		}
		// Cache and publish update
		i.cache.SetPosition(ctx, e)
		i.publishPositionUpdate(ctx, e, "opened")

	case *events.SpotOrderExecution:
		// Update order status
		if err := i.db.UpdateSpotOrderStatus(ctx, e.Order.OrderID, models.OrderStatusExecuted); err != nil {
			return err
		}
		// Create trade record
		if err := i.db.CreateTrade(ctx, e.Trade); err != nil {
			return err
		}
		// Publish updates
		i.publishOrderUpdate(ctx, e.Order, "executed")
		i.publishTradeUpdate(ctx, e.Trade)

	case *events.PerpetualOrderUpdate:
		if err := i.db.UpdatePerpetualOrderStatus(ctx, e.OrderID, e.Status, e.PositionID); err != nil {
			return err
		}

	case *events.PositionCloseData:
		if err := i.db.ClosePerpetualPosition(ctx, e.MtpID,
			e.ClosingPrice.String(), e.NetPnL.String(),
			e.ClosedBy, e.CloseTrigger); err != nil {
			return err
		}
		if e.Trade != nil {
			if err := i.db.CreateTrade(ctx, e.Trade); err != nil {
				return err
			}
			i.publishTradeUpdate(ctx, e.Trade)
		}

	case *events.PositionUpdate:
		// Handle position updates (stop loss, take profit, add collateral)
		// This would need specific update methods in the repository

	// CLOB Events
	case *models.CLOBMarket:
		if err := i.db.CreateCLOBMarket(ctx, e); err != nil {
			return err
		}
		// Cache market data
		i.cache.SetCLOBMarketData(ctx, e.MarketID, map[string]interface{}{
			"ticker":      e.Ticker,
			"base_asset":  e.BaseAsset,
			"quote_asset": e.QuoteAsset,
			"is_active":   e.IsActive,
		})

	case *models.CLOBOrder:
		if err := i.db.CreateCLOBOrder(ctx, e); err != nil {
			return err
		}
		// Add to order book and cache
		i.cache.AddCLOBOrderToBook(ctx, e)
		i.publishCLOBOrderUpdate(ctx, e, "created")

	case *models.CLOBTrade:
		if err := i.db.CreateCLOBTrade(ctx, e); err != nil {
			return err
		}
		// Add to cache and publish
		i.cache.AddCLOBTrade(ctx, e)
		i.publishCLOBTradeUpdate(ctx, e)
		// Update market stats
		i.db.UpdateCLOBMarketStats(ctx, e.MarketID)

	case *models.CLOBPosition:
		if err := i.db.CreateCLOBPosition(ctx, e); err != nil {
			return err
		}
		// Cache and publish update
		i.cache.SetCLOBPosition(ctx, e)
		i.publishCLOBPositionUpdate(ctx, e, "opened")

	case *events.CLOBOrderUpdate:
		if err := i.db.UpdateCLOBOrderStatus(ctx, e.OrderID, e.Status, e.FilledQuantity, e.RemainingQty); err != nil {
			return err
		}
		// Update cache
		if order, err := i.db.GetCLOBOrder(ctx, e.OrderID); err == nil && order != nil {
			if e.Status == models.CLOBOrderStatusCancelled || e.Status == models.CLOBOrderStatusFilled {
				i.cache.RemoveCLOBOrderFromBook(ctx, order)
			}
			i.publishCLOBOrderUpdate(ctx, order, string(e.Status))
		}

	case *events.CLOBOrderExecution:
		if err := i.db.UpdateCLOBOrderStatus(ctx, e.OrderID, e.Status, e.TotalFilled, e.RemainingQty); err != nil {
			return err
		}

	case *events.CLOBPositionUpdate:
		updates := make(map[string]interface{})
		if e.NewSize != nil {
			updates["size"] = e.NewSize.String()
		}
		if e.NewMargin != nil {
			updates["margin"] = e.NewMargin.String()
		}
		if e.NewMarkPrice != nil {
			updates["mark_price"] = e.NewMarkPrice.String()
		}
		if e.NewLiquidationPrice != nil {
			updates["liquidation_price"] = e.NewLiquidationPrice.String()
		}
		if e.NewUnrealizedPnL != nil {
			updates["unrealized_pnl"] = e.NewUnrealizedPnL.String()
		}
		
		if e.Action == "closed" && e.ClosePrice != nil && e.RealizedPnL != nil {
			if err := i.db.CloseCLOBPosition(ctx, e.PositionID, *e.ClosePrice, *e.RealizedPnL); err != nil {
				return err
			}
			// Remove from cache
			if pos, err := i.db.GetCLOBPosition(ctx, e.PositionID); err == nil && pos != nil {
				i.cache.RemoveCLOBPosition(ctx, pos)
				i.publishCLOBPositionUpdate(ctx, pos, "closed")
			}
		} else if len(updates) > 0 {
			if err := i.db.UpdateCLOBPosition(ctx, e.PositionID, updates); err != nil {
				return err
			}
			// Update cache
			if pos, err := i.db.GetCLOBPosition(ctx, e.PositionID); err == nil && pos != nil {
				i.cache.SetCLOBPosition(ctx, pos)
				i.publishCLOBPositionUpdate(ctx, pos, "modified")
			}
		}

	case *models.CLOBLiquidation:
		if err := i.db.CreateCLOBLiquidation(ctx, e); err != nil {
			return err
		}
		// The position should already be closed by a position update event

	case *models.CLOBFundingRate:
		if err := i.db.CreateCLOBFundingRate(ctx, e); err != nil {
			return err
		}
	}

	return nil
}

func (i *Indexer) publishOrderUpdate(ctx context.Context, order interface{}, action string) {
	update := &models.WSOrderUpdate{
		Action:    action,
		Timestamp: time.Now(),
	}

	switch o := order.(type) {
	case *models.SpotOrder:
		update.OrderType = "spot"
		update.Order = o
		i.cache.PublishOrderUpdate(ctx, o.OwnerAddress, update)
	case *models.PerpetualOrder:
		update.OrderType = "perpetual"
		update.Order = o
		i.cache.PublishOrderUpdate(ctx, o.OwnerAddress, update)
	}
}

func (i *Indexer) publishPositionUpdate(ctx context.Context, position *models.PerpetualPosition, action string) {
	update := &models.WSPositionUpdate{
		Position:  position,
		Action:    action,
		Timestamp: time.Now(),
	}
	i.cache.PublishPositionUpdate(ctx, position.OwnerAddress, update)
}

func (i *Indexer) publishTradeUpdate(ctx context.Context, trade *models.Trade) {
	update := &models.WSTradeUpdate{
		Trade:     trade,
		Timestamp: time.Now(),
	}
	i.cache.PublishTradeUpdate(ctx, trade.Asset, update)
}

// CLOB publish methods
func (i *Indexer) publishCLOBOrderUpdate(ctx context.Context, order *models.CLOBOrder, action string) {
	update := &models.WSCLOBOrderUpdate{
		Action:    action,
		Order:     order,
		Timestamp: time.Now(),
	}
	i.cache.PublishCLOBOrderUpdate(ctx, order.Owner, update)
}

func (i *Indexer) publishCLOBPositionUpdate(ctx context.Context, position *models.CLOBPosition, action string) {
	update := &models.WSCLOBPositionUpdate{
		Action:    action,
		Position:  position,
		Timestamp: time.Now(),
	}
	i.cache.PublishCLOBPositionUpdate(ctx, position.Owner, update)
}

func (i *Indexer) publishCLOBTradeUpdate(ctx context.Context, trade *models.CLOBTrade) {
	update := &models.WSCLOBTradeUpdate{
		Trade:     trade,
		Timestamp: time.Now(),
	}
	i.cache.PublishCLOBTradeUpdate(ctx, trade.MarketID, update)
}

func (i *Indexer) orderBookAggregator(ctx context.Context) {
	ticker := time.NewTicker(i.config.Indexer.OrderBookUpdateInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := i.aggregateOrderBooks(ctx); err != nil {
				i.logger.Error("Failed to aggregate order books", zap.Error(err))
			}
		}
	}
}

func (i *Indexer) aggregateOrderBooks(ctx context.Context) error {
	// Aggregate CLOB order books
	markets, err := i.db.GetActiveCLOBMarkets(ctx)
	if err != nil {
		return fmt.Errorf("failed to get active CLOB markets: %w", err)
	}

	for _, market := range markets {
		if err := i.aggregateCLOBOrderBook(ctx, market.MarketID); err != nil {
			i.logger.Error("Failed to aggregate CLOB order book",
				zap.Uint64("market_id", market.MarketID),
				zap.Error(err),
			)
		}
	}

	// TODO: Also aggregate spot order books when implemented
	return nil
}

func (i *Indexer) aggregateCLOBOrderBook(ctx context.Context, marketID uint64) error {
	// Get active orders for the market
	orders, err := i.db.GetActiveOrdersForMarket(ctx, marketID)
	if err != nil {
		return err
	}

	// Build order book levels
	bids := make([]models.OrderBookLevel, 0)
	asks := make([]models.OrderBookLevel, 0)
	bidsByPrice := make(map[string]*models.OrderBookLevel)
	asksByPrice := make(map[string]*models.OrderBookLevel)

	var totalBidVolume, totalAskVolume decimal.Decimal

	for _, order := range orders {
		priceKey := order.Price.String()
		
		if order.OrderType == models.CLOBOrderTypeLimitBuy {
			if level, exists := bidsByPrice[priceKey]; exists {
				level.Quantity = level.Quantity.Add(order.RemainingAmount)
				level.Orders++
			} else {
				bidsByPrice[priceKey] = &models.OrderBookLevel{
					Price:    order.Price,
					Quantity: order.RemainingAmount,
					Orders:   1,
				}
			}
			totalBidVolume = totalBidVolume.Add(order.RemainingAmount)
		} else if order.OrderType == models.CLOBOrderTypeLimitSell {
			if level, exists := asksByPrice[priceKey]; exists {
				level.Quantity = level.Quantity.Add(order.RemainingAmount)
				level.Orders++
			} else {
				asksByPrice[priceKey] = &models.OrderBookLevel{
					Price:    order.Price,
					Quantity: order.RemainingAmount,
					Orders:   1,
				}
			}
			totalAskVolume = totalAskVolume.Add(order.RemainingAmount)
		}
	}

	// Convert maps to sorted slices
	for _, level := range bidsByPrice {
		bids = append(bids, *level)
	}
	for _, level := range asksByPrice {
		asks = append(asks, *level)
	}

	// Sort bids descending, asks ascending
	sort.Slice(bids, func(i, j int) bool {
		return bids[i].Price.GreaterThan(bids[j].Price)
	})
	sort.Slice(asks, func(i, j int) bool {
		return asks[i].Price.LessThan(asks[j].Price)
	})

	// Calculate best bid/ask and spread
	var bestBid, bestAsk, midPrice, spread *decimal.Decimal
	if len(bids) > 0 {
		bestBid = &bids[0].Price
	}
	if len(asks) > 0 {
		bestAsk = &asks[0].Price
	}
	if bestBid != nil && bestAsk != nil {
		mid := bestBid.Add(*bestAsk).Div(decimal.NewFromInt(2))
		midPrice = &mid
		spr := bestAsk.Sub(*bestBid)
		spread = &spr
	}

	// Create snapshot
	bidsJSON, _ := json.Marshal(bids)
	asksJSON, _ := json.Marshal(asks)

	snapshot := &models.CLOBOrderBookSnapshot{
		MarketID:       marketID,
		Bids:           bidsJSON,
		Asks:           asksJSON,
		BestBid:        bestBid,
		BestAsk:        bestAsk,
		MidPrice:       midPrice,
		Spread:         spread,
		TotalBidVolume: totalBidVolume,
		TotalAskVolume: totalAskVolume,
		SnapshotTime:   time.Now(),
		BlockHeight:    0, // Will be updated on next block
	}

	// Save snapshot
	if err := i.db.SaveCLOBOrderBookSnapshot(ctx, snapshot); err != nil {
		return err
	}

	// Publish WebSocket update
	update := &models.WSCLOBOrderBookUpdate{
		MarketID:  marketID,
		Bids:      bids[:min(20, len(bids))], // Top 20 levels
		Asks:      asks[:min(20, len(asks))], // Top 20 levels
		Timestamp: time.Now(),
	}
	i.cache.PublishCLOBOrderBookUpdate(ctx, marketID, update)

	return nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func getTxHash(tx tmtypes.Tx) string {
	return fmt.Sprintf("%X", tx.Hash())
}
