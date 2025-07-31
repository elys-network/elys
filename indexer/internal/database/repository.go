package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/elys-network/elys/indexer/internal/models"
	"go.uber.org/zap"
)

type Repository struct {
	db     *DB
	logger *zap.Logger
}

func NewRepository(db *DB, logger *zap.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}

// Spot Orders
func (r *Repository) CreateSpotOrder(ctx context.Context, order *models.SpotOrder) error {
	query := `
		INSERT INTO spot_orders (
			order_id, order_type, owner_address, order_target_denom,
			order_price, order_amount, status, created_at, block_height, tx_hash
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (order_id) DO UPDATE SET
			status = EXCLUDED.status,
			executed_at = CASE WHEN EXCLUDED.status = 'EXECUTED' THEN NOW() ELSE spot_orders.executed_at END,
			closed_at = CASE WHEN EXCLUDED.status = 'CLOSED' THEN NOW() ELSE spot_orders.closed_at END`

	_, err := r.db.ExecContext(ctx, query,
		order.OrderID, order.OrderType, order.OwnerAddress, order.OrderTargetDenom,
		order.OrderPrice, order.OrderAmount.String(), order.Status, order.CreatedAt,
		order.BlockHeight, order.TxHash,
	)
	return err
}

func (r *Repository) UpdateSpotOrderStatus(ctx context.Context, orderID uint64, status models.OrderStatus) error {
	query := `UPDATE spot_orders SET status = $1, 
		executed_at = CASE WHEN $1 = 'EXECUTED' THEN NOW() ELSE executed_at END,
		closed_at = CASE WHEN $1 = 'CLOSED' THEN NOW() ELSE closed_at END
		WHERE order_id = $2`

	_, err := r.db.ExecContext(ctx, query, status, orderID)
	return err
}

func (r *Repository) GetSpotOrdersByOwner(ctx context.Context, owner string, limit int) ([]*models.SpotOrder, error) {
	query := `SELECT * FROM spot_orders WHERE owner_address = $1 ORDER BY created_at DESC LIMIT $2`

	rows, err := r.db.QueryContext(ctx, query, owner, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*models.SpotOrder
	for rows.Next() {
		order := &models.SpotOrder{}
		if err := rows.Scan(
			&order.ID, &order.OrderID, &order.OrderType, &order.OwnerAddress,
			&order.OrderTargetDenom, &order.OrderPrice, &order.OrderAmount,
			&order.Status, &order.CreatedAt, &order.ExecutedAt, &order.ClosedAt,
			&order.BlockHeight, &order.TxHash,
		); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, rows.Err()
}

// Perpetual Orders
func (r *Repository) CreatePerpetualOrder(ctx context.Context, order *models.PerpetualOrder) error {
	query := `
		INSERT INTO perpetual_orders (
			order_id, owner_address, perpetual_order_type, position,
			trigger_price, collateral, leverage, take_profit_price,
			stop_loss_price, position_id, status, created_at, block_height, tx_hash
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		ON CONFLICT (order_id) DO UPDATE SET
			status = EXCLUDED.status,
			position_id = COALESCE(EXCLUDED.position_id, perpetual_orders.position_id),
			executed_at = CASE WHEN EXCLUDED.status = 'EXECUTED' THEN NOW() ELSE perpetual_orders.executed_at END,
			cancelled_at = CASE WHEN EXCLUDED.status = 'CANCELLED' THEN NOW() ELSE perpetual_orders.cancelled_at END`

	_, err := r.db.ExecContext(ctx, query,
		order.OrderID, order.OwnerAddress, order.PerpetualOrderType, order.Position,
		order.TriggerPrice, order.Collateral.String(), order.Leverage.String(),
		order.TakeProfitPrice.String(), order.StopLossPrice.String(), order.PositionID,
		order.Status, order.CreatedAt, order.BlockHeight, order.TxHash,
	)
	return err
}

func (r *Repository) UpdatePerpetualOrderStatus(ctx context.Context, orderID uint64, status models.OrderStatus, positionID *uint64) error {
	query := `UPDATE perpetual_orders SET status = $1, position_id = COALESCE($2, position_id),
		executed_at = CASE WHEN $1 = 'EXECUTED' THEN NOW() ELSE executed_at END,
		cancelled_at = CASE WHEN $1 = 'CANCELLED' THEN NOW() ELSE cancelled_at END
		WHERE order_id = $3`

	_, err := r.db.ExecContext(ctx, query, status, positionID, orderID)
	return err
}

// Perpetual Positions
func (r *Repository) CreatePerpetualPosition(ctx context.Context, pos *models.PerpetualPosition) error {
	query := `
		INSERT INTO perpetual_positions (
			mtp_id, owner_address, amm_pool_id, position, collateral_asset,
			collateral, liabilities, custody, mtp_health, open_price,
			stop_loss_price, take_profit_price, opened_at, block_height, tx_hash
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		ON CONFLICT (mtp_id) DO UPDATE SET
			collateral = EXCLUDED.collateral,
			liabilities = EXCLUDED.liabilities,
			custody = EXCLUDED.custody,
			mtp_health = EXCLUDED.mtp_health,
			stop_loss_price = EXCLUDED.stop_loss_price,
			take_profit_price = EXCLUDED.take_profit_price`

	_, err := r.db.ExecContext(ctx, query,
		pos.MtpID, pos.OwnerAddress, pos.AmmPoolID, pos.Position, pos.CollateralAsset,
		pos.Collateral.String(), pos.Liabilities.String(), pos.Custody.String(),
		pos.MtpHealth.String(), pos.OpenPrice.String(), pos.StopLossPrice.String(),
		pos.TakeProfitPrice.String(), pos.OpenedAt, pos.BlockHeight, pos.TxHash,
	)
	return err
}

func (r *Repository) ClosePerpetualPosition(ctx context.Context, mtpID uint64, closingPrice, netPnL string, closedBy, trigger string) error {
	query := `UPDATE perpetual_positions SET 
		closing_price = $1, net_pnl = $2, closed_at = NOW(), 
		closed_by = $3, close_trigger = $4
		WHERE mtp_id = $5`

	_, err := r.db.ExecContext(ctx, query, closingPrice, netPnL, closedBy, trigger, mtpID)
	return err
}

func (r *Repository) GetOpenPositions(ctx context.Context, owner string) ([]*models.PerpetualPosition, error) {
	query := `SELECT * FROM perpetual_positions 
		WHERE owner_address = $1 AND closed_at IS NULL 
		ORDER BY opened_at DESC`

	rows, err := r.db.QueryContext(ctx, query, owner)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var positions []*models.PerpetualPosition
	for rows.Next() {
		pos := &models.PerpetualPosition{}
		if err := scanPerpetualPosition(rows, pos); err != nil {
			return nil, err
		}
		positions = append(positions, pos)
	}

	return positions, rows.Err()
}

// Trades
func (r *Repository) CreateTrade(ctx context.Context, trade *models.Trade) error {
	query := `
		INSERT INTO trades (
			trade_type, reference_id, owner_address, asset, amount,
			price, fees, executed_at, block_height, tx_hash, event_type
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	_, err := r.db.ExecContext(ctx, query,
		trade.TradeType, trade.ReferenceID, trade.OwnerAddress, trade.Asset,
		trade.Amount.String(), trade.Price.String(), trade.Fees, trade.ExecutedAt,
		trade.BlockHeight, trade.TxHash, trade.EventType,
	)
	return err
}

func (r *Repository) GetTradeHistory(ctx context.Context, owner string, limit int) ([]*models.Trade, error) {
	query := `SELECT * FROM trades WHERE owner_address = $1 ORDER BY executed_at DESC LIMIT $2`

	rows, err := r.db.QueryContext(ctx, query, owner, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trades []*models.Trade
	for rows.Next() {
		trade := &models.Trade{}
		if err := scanTrade(rows, trade); err != nil {
			return nil, err
		}
		trades = append(trades, trade)
	}

	return trades, rows.Err()
}

// Order Book
func (r *Repository) SaveOrderBookSnapshot(ctx context.Context, snapshot *models.OrderBookSnapshot) error {
	bidsJSON, err := json.Marshal(snapshot.Bids)
	if err != nil {
		return err
	}

	asksJSON, err := json.Marshal(snapshot.Asks)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO order_book_snapshots (
			asset_pair, bids, asks, snapshot_time, block_height
		) VALUES ($1, $2, $3, $4, $5)`

	_, err = r.db.ExecContext(ctx, query,
		snapshot.AssetPair, bidsJSON, asksJSON,
		snapshot.SnapshotTime, snapshot.BlockHeight,
	)
	return err
}

func (r *Repository) GetLatestOrderBookSnapshot(ctx context.Context, assetPair string) (*models.OrderBookSnapshot, error) {
	query := `SELECT * FROM order_book_snapshots 
		WHERE asset_pair = $1 
		ORDER BY snapshot_time DESC 
		LIMIT 1`

	snapshot := &models.OrderBookSnapshot{}
	row := r.db.QueryRowContext(ctx, query, assetPair)

	var bidsJSON, asksJSON []byte
	err := row.Scan(
		&snapshot.ID, &snapshot.AssetPair, &bidsJSON, &asksJSON,
		&snapshot.SnapshotTime, &snapshot.BlockHeight,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if err := json.Unmarshal(bidsJSON, &snapshot.Bids); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(asksJSON, &snapshot.Asks); err != nil {
		return nil, err
	}

	return snapshot, nil
}

// Indexer State
func (r *Repository) GetIndexerState(ctx context.Context) (*models.IndexerState, error) {
	query := `SELECT * FROM indexer_state LIMIT 1`

	state := &models.IndexerState{}
	err := r.db.QueryRowContext(ctx, query).Scan(
		&state.ID, &state.LastProcessedHeight,
		&state.LastProcessedTime, &state.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return &models.IndexerState{
				LastProcessedHeight: 0,
				LastProcessedTime:   time.Now(),
			}, nil
		}
		return nil, err
	}

	return state, nil
}

func (r *Repository) UpdateIndexerState(ctx context.Context, height int64) error {
	query := `UPDATE indexer_state SET 
		last_processed_height = $1, 
		last_processed_time = NOW(), 
		updated_at = NOW()`

	_, err := r.db.ExecContext(ctx, query, height)
	return err
}

// WebSocket Subscriptions
func (r *Repository) CreateSubscription(ctx context.Context, sub *models.WebSocketSubscription) error {
	query := `
		INSERT INTO websocket_subscriptions (
			id, client_id, subscription_type, filters, created_at, last_ping
		) VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := r.db.ExecContext(ctx, query,
		sub.ID, sub.ClientID, sub.SubscriptionType,
		sub.Filters, sub.CreatedAt, sub.LastPing,
	)
	return err
}

func (r *Repository) UpdateSubscriptionPing(ctx context.Context, clientID string) error {
	query := `UPDATE websocket_subscriptions SET last_ping = NOW() WHERE client_id = $1`
	_, err := r.db.ExecContext(ctx, query, clientID)
	return err
}

func (r *Repository) DeleteSubscription(ctx context.Context, subID string) error {
	query := `DELETE FROM websocket_subscriptions WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, subID)
	return err
}

func (r *Repository) CleanupStaleSubscriptions(ctx context.Context, timeout time.Duration) error {
	query := `DELETE FROM websocket_subscriptions WHERE last_ping < $1`
	_, err := r.db.ExecContext(ctx, query, time.Now().Add(-timeout))
	return err
}

// Helper functions for scanning
func scanPerpetualPosition(rows *sql.Rows, pos *models.PerpetualPosition) error {
	var collateral, liabilities, custody, mtpHealth, openPrice, stopLoss, takeProfit string
	var closingPrice, netPnL sql.NullString

	err := rows.Scan(
		&pos.ID, &pos.MtpID, &pos.OwnerAddress, &pos.AmmPoolID, &pos.Position,
		&pos.CollateralAsset, &collateral, &liabilities, &custody, &mtpHealth,
		&openPrice, &closingPrice, &netPnL, &stopLoss, &takeProfit,
		&pos.OpenedAt, &pos.ClosedAt, &pos.ClosedBy, &pos.CloseTrigger,
		&pos.BlockHeight, &pos.TxHash,
	)
	if err != nil {
		return err
	}

	// Convert string values back to math types
	// This would need proper conversion logic
	return nil
}

func scanTrade(rows *sql.Rows, trade *models.Trade) error {
	var amount, price string

	err := rows.Scan(
		&trade.ID, &trade.TradeType, &trade.ReferenceID, &trade.OwnerAddress,
		&trade.Asset, &amount, &price, &trade.Fees, &trade.ExecutedAt,
		&trade.BlockHeight, &trade.TxHash, &trade.EventType,
	)
	if err != nil {
		return err
	}

	// Convert string values back to math types
	// This would need proper conversion logic
	return nil
}
