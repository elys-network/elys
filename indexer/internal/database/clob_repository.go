package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/elys-network/elys/indexer/internal/models"
	"github.com/shopspring/decimal"
)

// CLOB Markets
func (r *Repository) CreateCLOBMarket(ctx context.Context, market *models.CLOBMarket) error {
	query := `
		INSERT INTO clob_markets (
			market_id, ticker, base_asset, quote_asset, tick_size, lot_size,
			min_order_size, max_order_size, max_leverage, initial_margin_fraction,
			maintenance_margin_fraction, funding_interval, next_funding_time,
			is_active, created_at, updated_at, block_height
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
		ON CONFLICT (market_id) DO UPDATE SET
			ticker = EXCLUDED.ticker,
			is_active = EXCLUDED.is_active,
			updated_at = EXCLUDED.updated_at`

	_, err := r.db.ExecContext(ctx, query,
		market.MarketID, market.Ticker, market.BaseAsset, market.QuoteAsset,
		market.TickSize.String(), market.LotSize.String(), market.MinOrderSize.String(),
		market.MaxOrderSize.String(), market.MaxLeverage.String(),
		market.InitialMarginFraction.String(), market.MaintenanceMarginFraction.String(),
		market.FundingInterval, market.NextFundingTime, market.IsActive,
		market.CreatedAt, market.UpdatedAt, market.BlockHeight,
	)
	return err
}

func (r *Repository) GetCLOBMarket(ctx context.Context, marketID uint64) (*models.CLOBMarket, error) {
	query := `SELECT * FROM clob_markets WHERE market_id = $1`

	market := &models.CLOBMarket{}
	err := r.db.QueryRowContext(ctx, query, marketID).Scan(
		&market.ID, &market.MarketID, &market.Ticker, &market.BaseAsset, &market.QuoteAsset,
		&market.TickSize, &market.LotSize, &market.MinOrderSize, &market.MaxOrderSize,
		&market.MaxLeverage, &market.InitialMarginFraction, &market.MaintenanceMarginFraction,
		&market.FundingInterval, &market.NextFundingTime, &market.IsActive,
		&market.CreatedAt, &market.UpdatedAt, &market.BlockHeight,
	)
	if err != nil {
		return nil, err
	}
	return market, nil
}

func (r *Repository) GetActiveCLOBMarkets(ctx context.Context) ([]*models.CLOBMarket, error) {
	query := `SELECT * FROM clob_markets WHERE is_active = true ORDER BY market_id`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var markets []*models.CLOBMarket
	for rows.Next() {
		market := &models.CLOBMarket{}
		err := rows.Scan(
			&market.ID, &market.MarketID, &market.Ticker, &market.BaseAsset, &market.QuoteAsset,
			&market.TickSize, &market.LotSize, &market.MinOrderSize, &market.MaxOrderSize,
			&market.MaxLeverage, &market.InitialMarginFraction, &market.MaintenanceMarginFraction,
			&market.FundingInterval, &market.NextFundingTime, &market.IsActive,
			&market.CreatedAt, &market.UpdatedAt, &market.BlockHeight,
		)
		if err != nil {
			return nil, err
		}
		markets = append(markets, market)
	}
	return markets, rows.Err()
}

// CLOB Orders
func (r *Repository) CreateCLOBOrder(ctx context.Context, order *models.CLOBOrder) error {
	query := `
		INSERT INTO clob_orders (
			order_id, market_id, counter, owner, sub_account_id, order_type,
			price, amount, filled_amount, remaining_amount, status, time_in_force,
			post_only, reduce_only, created_at, updated_at, block_height, tx_hash
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)
		ON CONFLICT (market_id, counter) DO UPDATE SET
			filled_amount = EXCLUDED.filled_amount,
			remaining_amount = EXCLUDED.remaining_amount,
			status = EXCLUDED.status,
			updated_at = EXCLUDED.updated_at`

	_, err := r.db.ExecContext(ctx, query,
		order.OrderID, order.MarketID, order.Counter, order.Owner, order.SubAccountID,
		order.OrderType, order.Price.String(), order.Amount.String(),
		order.FilledAmount.String(), order.RemainingAmount.String(), order.Status,
		order.TimeInForce, order.PostOnly, order.ReduceOnly,
		order.CreatedAt, order.UpdatedAt, order.BlockHeight, order.TxHash,
	)
	return err
}

func (r *Repository) UpdateCLOBOrderStatus(ctx context.Context, orderID uint64, status models.CLOBOrderStatus, filledAmount, remainingAmount decimal.Decimal) error {
	query := `UPDATE clob_orders SET 
		status = $1, filled_amount = $2, remaining_amount = $3,
		updated_at = NOW(),
		executed_at = CASE WHEN $1 IN ('FILLED', 'PARTIALLY_FILLED') THEN NOW() ELSE executed_at END,
		cancelled_at = CASE WHEN $1 = 'CANCELLED' THEN NOW() ELSE cancelled_at END
		WHERE order_id = $4`

	_, err := r.db.ExecContext(ctx, query, status, filledAmount.String(), remainingAmount.String(), orderID)
	return err
}

func (r *Repository) GetCLOBOrdersByOwner(ctx context.Context, owner string, status []models.CLOBOrderStatus, limit int) ([]*models.CLOBOrder, error) {
	query := `SELECT * FROM clob_orders WHERE owner = $1`
	args := []interface{}{owner}

	if len(status) > 0 {
		query += " AND status = ANY($2)"
		args = append(args, status)
	}

	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d", len(args)+1)
	args = append(args, limit)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanCLOBOrders(rows)
}

func (r *Repository) GetActiveOrdersForMarket(ctx context.Context, marketID uint64) ([]*models.CLOBOrder, error) {
	query := `SELECT * FROM clob_orders 
		WHERE market_id = $1 AND status IN ('PENDING', 'PARTIALLY_FILLED')
		ORDER BY price DESC, counter ASC`

	rows, err := r.db.QueryContext(ctx, query, marketID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanCLOBOrders(rows)
}

func (r *Repository) GetCLOBOrder(ctx context.Context, orderID uint64) (*models.CLOBOrder, error) {
	query := `SELECT * FROM clob_orders WHERE order_id = $1`

	rows, err := r.db.QueryContext(ctx, query, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders, err := scanCLOBOrders(rows)
	if err != nil {
		return nil, err
	}
	if len(orders) == 0 {
		return nil, sql.ErrNoRows
	}
	return orders[0], nil
}

// CLOB Trades
func (r *Repository) CreateCLOBTrade(ctx context.Context, trade *models.CLOBTrade) error {
	query := `
		INSERT INTO clob_trades (
			market_id, buyer, buyer_sub_account_id, seller, seller_sub_account_id,
			buyer_order_id, seller_order_id, price, quantity, trade_value,
			buyer_fee, seller_fee, is_buyer_taker, is_buyer_liquidation,
			is_seller_liquidation, executed_at, block_height, tx_hash
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)`

	_, err := r.db.ExecContext(ctx, query,
		trade.MarketID, trade.Buyer, trade.BuyerSubAccountID, trade.Seller,
		trade.SellerSubAccountID, trade.BuyerOrderID, trade.SellerOrderID,
		trade.Price.String(), trade.Quantity.String(), trade.TradeValue.String(),
		trade.BuyerFee.String(), trade.SellerFee.String(), trade.IsBuyerTaker,
		trade.IsBuyerLiquidation, trade.IsSellerLiquidation, trade.ExecutedAt,
		trade.BlockHeight, trade.TxHash,
	)
	return err
}

func (r *Repository) GetCLOBTradeHistory(ctx context.Context, marketID uint64, owner string, limit int) ([]*models.CLOBTrade, error) {
	query := `SELECT * FROM clob_trades WHERE 1=1`
	args := []interface{}{}
	argCount := 0

	if marketID > 0 {
		argCount++
		query += fmt.Sprintf(" AND market_id = $%d", argCount)
		args = append(args, marketID)
	}

	if owner != "" {
		argCount++
		query += fmt.Sprintf(" AND (buyer = $%d OR seller = $%d)", argCount, argCount)
		args = append(args, owner)
	}

	argCount++
	query += fmt.Sprintf(" ORDER BY executed_at DESC LIMIT $%d", argCount)
	args = append(args, limit)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanCLOBTrades(rows)
}

// CLOB Positions
func (r *Repository) CreateCLOBPosition(ctx context.Context, pos *models.CLOBPosition) error {
	query := `
		INSERT INTO clob_positions (
			market_id, owner, sub_account_id, side, size, notional, entry_price,
			mark_price, liquidation_price, margin, margin_ratio, unrealized_pnl,
			realized_pnl, cumulative_funding, opened_at, updated_at, block_height, tx_hash
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)`

	_, err := r.db.ExecContext(ctx, query,
		pos.MarketID, pos.Owner, pos.SubAccountID, pos.Side, pos.Size.String(),
		pos.Notional.String(), pos.EntryPrice.String(), pos.MarkPrice.String(),
		pos.LiquidationPrice.String(), pos.Margin.String(), pos.MarginRatio.String(),
		pos.UnrealizedPnL.String(), pos.RealizedPnL.String(), pos.CumulativeFunding.String(),
		pos.OpenedAt, pos.UpdatedAt, pos.BlockHeight, pos.TxHash,
	)
	return err
}

func (r *Repository) UpdateCLOBPosition(ctx context.Context, positionID uint64, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return nil
	}

	query := "UPDATE clob_positions SET updated_at = NOW()"
	args := []interface{}{}
	argCount := 0

	for field, value := range updates {
		argCount++
		query += fmt.Sprintf(", %s = $%d", field, argCount)
		args = append(args, value)
	}

	argCount++
	query += fmt.Sprintf(" WHERE position_id = $%d", argCount)
	args = append(args, positionID)

	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}

func (r *Repository) CloseCLOBPosition(ctx context.Context, positionID uint64, closePrice, realizedPnL decimal.Decimal) error {
	query := `UPDATE clob_positions SET 
		closed_at = NOW(), close_price = $1, realized_pnl = $2, updated_at = NOW()
		WHERE position_id = $3`

	_, err := r.db.ExecContext(ctx, query, closePrice.String(), realizedPnL.String(), positionID)
	return err
}

func (r *Repository) GetOpenCLOBPositions(ctx context.Context, owner string, marketID uint64) ([]*models.CLOBPosition, error) {
	query := `SELECT * FROM clob_positions WHERE closed_at IS NULL`
	args := []interface{}{}
	argCount := 0

	if owner != "" {
		argCount++
		query += fmt.Sprintf(" AND owner = $%d", argCount)
		args = append(args, owner)
	}

	if marketID > 0 {
		argCount++
		query += fmt.Sprintf(" AND market_id = $%d", argCount)
		args = append(args, marketID)
	}

	query += " ORDER BY opened_at DESC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanCLOBPositions(rows)
}

func (r *Repository) GetCLOBPosition(ctx context.Context, positionID uint64) (*models.CLOBPosition, error) {
	query := `SELECT * FROM clob_positions WHERE position_id = $1`

	rows, err := r.db.QueryContext(ctx, query, positionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	positions, err := scanCLOBPositions(rows)
	if err != nil {
		return nil, err
	}
	if len(positions) == 0 {
		return nil, sql.ErrNoRows
	}
	return positions[0], nil
}

// CLOB Order Book Snapshots
func (r *Repository) SaveCLOBOrderBookSnapshot(ctx context.Context, snapshot *models.CLOBOrderBookSnapshot) error {
	query := `
		INSERT INTO clob_order_book_snapshots (
			market_id, bids, asks, best_bid, best_ask, mid_price, spread,
			total_bid_volume, total_ask_volume, snapshot_time, block_height
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	_, err := r.db.ExecContext(ctx, query,
		snapshot.MarketID, snapshot.Bids, snapshot.Asks,
		decimalPtrToString(snapshot.BestBid), decimalPtrToString(snapshot.BestAsk),
		decimalPtrToString(snapshot.MidPrice), decimalPtrToString(snapshot.Spread),
		snapshot.TotalBidVolume.String(), snapshot.TotalAskVolume.String(),
		snapshot.SnapshotTime, snapshot.BlockHeight,
	)
	return err
}

func (r *Repository) GetLatestCLOBOrderBookSnapshot(ctx context.Context, marketID uint64) (*models.CLOBOrderBookSnapshot, error) {
	query := `SELECT * FROM clob_order_book_snapshots 
		WHERE market_id = $1 
		ORDER BY snapshot_time DESC 
		LIMIT 1`

	snapshot := &models.CLOBOrderBookSnapshot{}
	var bestBid, bestAsk, midPrice, spread sql.NullString

	err := r.db.QueryRowContext(ctx, query, marketID).Scan(
		&snapshot.ID, &snapshot.MarketID, &snapshot.Bids, &snapshot.Asks,
		&bestBid, &bestAsk, &midPrice, &spread,
		&snapshot.TotalBidVolume, &snapshot.TotalAskVolume,
		&snapshot.SnapshotTime, &snapshot.BlockHeight,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// Convert nullable strings to decimal pointers
	snapshot.BestBid = stringToDecimalPtr(bestBid)
	snapshot.BestAsk = stringToDecimalPtr(bestAsk)
	snapshot.MidPrice = stringToDecimalPtr(midPrice)
	snapshot.Spread = stringToDecimalPtr(spread)

	return snapshot, nil
}

// CLOB Funding Rates
func (r *Repository) CreateCLOBFundingRate(ctx context.Context, rate *models.CLOBFundingRate) error {
	query := `
		INSERT INTO clob_funding_rates (
			market_id, funding_rate, premium_rate, mark_price, index_price,
			timestamp, next_funding_time, block_height
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := r.db.ExecContext(ctx, query,
		rate.MarketID, rate.FundingRate.String(), rate.PremiumRate.String(),
		rate.MarkPrice.String(), rate.IndexPrice.String(),
		rate.Timestamp, rate.NextFundingTime, rate.BlockHeight,
	)
	return err
}

// CLOB Liquidations
func (r *Repository) CreateCLOBLiquidation(ctx context.Context, liq *models.CLOBLiquidation) error {
	query := `
		INSERT INTO clob_liquidations (
			market_id, position_id, owner, sub_account_id, liquidator, side,
			size, price, liquidation_fee, insurance_fund_contribution, is_adl,
			liquidated_at, block_height, tx_hash
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`

	_, err := r.db.ExecContext(ctx, query,
		liq.MarketID, liq.PositionID, liq.Owner, liq.SubAccountID, liq.Liquidator,
		liq.Side, liq.Size.String(), liq.Price.String(), liq.LiquidationFee.String(),
		liq.InsuranceFundContribution.String(), liq.IsADL, liq.LiquidatedAt,
		liq.BlockHeight, liq.TxHash,
	)
	return err
}

// CLOB Market Stats
func (r *Repository) UpdateCLOBMarketStats(ctx context.Context, marketID uint64) error {
	// This is called by the trigger, but we can also call it manually
	query := `
		INSERT INTO clob_market_stats (market_id, updated_at)
		VALUES ($1, NOW())
		ON CONFLICT (market_id) DO UPDATE SET updated_at = NOW()`

	_, err := r.db.ExecContext(ctx, query, marketID)
	return err
}

func (r *Repository) GetCLOBMarketStats(ctx context.Context, marketID uint64) (*models.CLOBMarketStats, error) {
	query := `SELECT * FROM clob_market_stats WHERE market_id = $1`

	stats := &models.CLOBMarketStats{}
	var high24h, low24h, lastPrice sql.NullString
	var lastTradeTime sql.NullTime

	err := r.db.QueryRowContext(ctx, query, marketID).Scan(
		&stats.ID, &stats.MarketID, &stats.Volume24h, &stats.Trades24h,
		&high24h, &low24h, &stats.OpenInterest, &stats.OpenInterestNotional,
		&lastPrice, &lastTradeTime, &stats.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	stats.High24h = stringToDecimalPtr(high24h)
	stats.Low24h = stringToDecimalPtr(low24h)
	stats.LastPrice = stringToDecimalPtr(lastPrice)
	if lastTradeTime.Valid {
		stats.LastTradeTime = &lastTradeTime.Time
	}

	return stats, nil
}

// Helper functions for scanning CLOB data
func scanCLOBOrders(rows *sql.Rows) ([]*models.CLOBOrder, error) {
	var orders []*models.CLOBOrder
	for rows.Next() {
		order := &models.CLOBOrder{}
		var executedAt, cancelledAt, expiresAt sql.NullTime

		err := rows.Scan(
			&order.ID, &order.OrderID, &order.MarketID, &order.Counter,
			&order.Owner, &order.SubAccountID, &order.OrderType, &order.Price,
			&order.Amount, &order.FilledAmount, &order.RemainingAmount,
			&order.Status, &order.TimeInForce, &order.PostOnly, &order.ReduceOnly,
			&order.CreatedAt, &order.UpdatedAt, &executedAt, &cancelledAt,
			&expiresAt, &order.BlockHeight, &order.TxHash,
		)
		if err != nil {
			return nil, err
		}

		if executedAt.Valid {
			order.ExecutedAt = &executedAt.Time
		}
		if cancelledAt.Valid {
			order.CancelledAt = &cancelledAt.Time
		}
		if expiresAt.Valid {
			order.ExpiresAt = &expiresAt.Time
		}

		orders = append(orders, order)
	}
	return orders, rows.Err()
}

func scanCLOBTrades(rows *sql.Rows) ([]*models.CLOBTrade, error) {
	var trades []*models.CLOBTrade
	for rows.Next() {
		trade := &models.CLOBTrade{}
		err := rows.Scan(
			&trade.ID, &trade.TradeID, &trade.MarketID, &trade.Buyer,
			&trade.BuyerSubAccountID, &trade.Seller, &trade.SellerSubAccountID,
			&trade.BuyerOrderID, &trade.SellerOrderID, &trade.Price,
			&trade.Quantity, &trade.TradeValue, &trade.BuyerFee, &trade.SellerFee,
			&trade.IsBuyerTaker, &trade.IsBuyerLiquidation, &trade.IsSellerLiquidation,
			&trade.ExecutedAt, &trade.BlockHeight, &trade.TxHash,
		)
		if err != nil {
			return nil, err
		}
		trades = append(trades, trade)
	}
	return trades, rows.Err()
}

func scanCLOBPositions(rows *sql.Rows) ([]*models.CLOBPosition, error) {
	var positions []*models.CLOBPosition
	for rows.Next() {
		pos := &models.CLOBPosition{}
		var lastFundingTime, closedAt sql.NullTime
		var closePrice sql.NullString

		err := rows.Scan(
			&pos.ID, &pos.PositionID, &pos.MarketID, &pos.Owner, &pos.SubAccountID,
			&pos.Side, &pos.Size, &pos.Notional, &pos.EntryPrice, &pos.MarkPrice,
			&pos.LiquidationPrice, &pos.Margin, &pos.MarginRatio, &pos.UnrealizedPnL,
			&pos.RealizedPnL, &pos.CumulativeFunding, &pos.LastFundingPayment,
			&lastFundingTime, &pos.OpenedAt, &pos.UpdatedAt, &closedAt,
			&closePrice, &pos.IsLiquidated, &pos.BlockHeight, &pos.TxHash,
		)
		if err != nil {
			return nil, err
		}

		if lastFundingTime.Valid {
			pos.LastFundingTime = &lastFundingTime.Time
		}
		if closedAt.Valid {
			pos.ClosedAt = &closedAt.Time
		}
		if closePrice.Valid {
			if d, err := decimal.NewFromString(closePrice.String); err == nil {
				pos.ClosePrice = &d
			}
		}

		positions = append(positions, pos)
	}
	return positions, rows.Err()
}

func decimalPtrToString(d *decimal.Decimal) interface{} {
	if d == nil {
		return nil
	}
	return d.String()
}

func stringToDecimalPtr(s sql.NullString) *decimal.Decimal {
	if !s.Valid {
		return nil
	}
	d, err := decimal.NewFromString(s.String)
	if err != nil {
		return nil
	}
	return &d
}
