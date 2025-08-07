package indexer

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/rpc/client/http"
	tmtypes "github.com/cometbft/cometbft/types"
	"github.com/elys-network/elys/indexer/internal/cache"
	"github.com/elys-network/elys/indexer/internal/config"
	"github.com/elys-network/elys/indexer/internal/database"
	"github.com/elys-network/elys/indexer/internal/events"
	"github.com/elys-network/elys/indexer/internal/models"
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
	// This would aggregate open orders into order book snapshots
	// For now, this is a placeholder - actual implementation would:
	// 1. Query open spot orders grouped by asset pair
	// 2. Build bid/ask arrays
	// 3. Save snapshots to database
	// 4. Publish updates via WebSocket
	return nil
}

func getTxHash(tx tmtypes.Tx) string {
	return fmt.Sprintf("%X", tx.Hash())
}
