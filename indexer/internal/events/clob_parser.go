package events

import (
	"context"
	"fmt"
	"strconv"
	"time"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/elys-network/elys/indexer/internal/models"
	clobtypes "github.com/elys-network/elys/v7/x/clob/types"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

type CLOBParser struct {
	logger *zap.Logger
}

func NewCLOBParser(logger *zap.Logger) *CLOBParser {
	return &CLOBParser{
		logger: logger,
	}
}

func (p *CLOBParser) ParseEvents(ctx context.Context, events []abci.Event, height int64, txHash string) ([]interface{}, error) {
	var results []interface{}

	for _, event := range events {
		switch event.Type {
		case clobtypes.EventPlaceLimitOrder:
			order, err := p.parsePlaceLimitOrder(event, height, txHash)
			if err != nil {
				p.logger.Error("Failed to parse place limit order event", zap.Error(err))
				continue
			}
			results = append(results, order)

		case clobtypes.EventPlaceMarketOrder:
			order, err := p.parsePlaceMarketOrder(event, height, txHash)
			if err != nil {
				p.logger.Error("Failed to parse place market order event", zap.Error(err))
				continue
			}
			results = append(results, order)

		case clobtypes.EventCancelOrder:
			update, err := p.parseCancelOrder(event, height, txHash)
			if err != nil {
				p.logger.Error("Failed to parse cancel order event", zap.Error(err))
				continue
			}
			results = append(results, update)

		case clobtypes.EventOrderExecuted:
			execution, err := p.parseOrderExecuted(event, height, txHash)
			if err != nil {
				p.logger.Error("Failed to parse order executed event", zap.Error(err))
				continue
			}
			results = append(results, execution)

		case clobtypes.EventTrade:
			trade, err := p.parseTrade(event, height, txHash)
			if err != nil {
				p.logger.Error("Failed to parse trade event", zap.Error(err))
				continue
			}
			results = append(results, trade)

		case clobtypes.EventPositionOpened:
			position, err := p.parsePositionOpened(event, height, txHash)
			if err != nil {
				p.logger.Error("Failed to parse position opened event", zap.Error(err))
				continue
			}
			results = append(results, position)

		case clobtypes.EventPositionClosed:
			update, err := p.parsePositionClosed(event, height, txHash)
			if err != nil {
				p.logger.Error("Failed to parse position closed event", zap.Error(err))
				continue
			}
			results = append(results, update)

		case clobtypes.EventPositionModified:
			update, err := p.parsePositionModified(event, height, txHash)
			if err != nil {
				p.logger.Error("Failed to parse position modified event", zap.Error(err))
				continue
			}
			results = append(results, update)

		case clobtypes.EventLiquidation:
			liquidation, err := p.parseLiquidation(event, height, txHash)
			if err != nil {
				p.logger.Error("Failed to parse liquidation event", zap.Error(err))
				continue
			}
			results = append(results, liquidation)

		case clobtypes.EventFundingRateUpdate:
			fundingRate, err := p.parseFundingRateUpdate(event, height, txHash)
			if err != nil {
				p.logger.Error("Failed to parse funding rate update event", zap.Error(err))
				continue
			}
			results = append(results, fundingRate)

		case clobtypes.EventCreateMarket:
			market, err := p.parseCreateMarket(event, height, txHash)
			if err != nil {
				p.logger.Error("Failed to parse create market event", zap.Error(err))
				continue
			}
			results = append(results, market)
		}
	}

	return results, nil
}

func (p *CLOBParser) parsePlaceLimitOrder(event abci.Event, height int64, txHash string) (*models.CLOBOrder, error) {
	attrs := getAttributeMap(event.Attributes)

	marketID, err := strconv.ParseUint(attrs[clobtypes.AttributeMarketId], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse market id: %w", err)
	}

	orderID, err := strconv.ParseUint(attrs[clobtypes.AttributeOrderId], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse order id: %w", err)
	}

	counter, err := strconv.ParseUint(attrs["counter"], 10, 64)
	if err != nil {
		counter = orderID // Fallback to orderID if counter not present
	}

	subAccountID, _ := strconv.ParseUint(attrs["sub_account_id"], 10, 64)

	price, err := decimal.NewFromString(attrs[clobtypes.AttributePrice])
	if err != nil {
		return nil, fmt.Errorf("failed to parse price: %w", err)
	}

	amount, err := decimal.NewFromString(attrs[clobtypes.AttributeQuantity])
	if err != nil {
		return nil, fmt.Errorf("failed to parse quantity: %w", err)
	}

	orderType := models.CLOBOrderTypeLimitBuy
	if attrs[clobtypes.AttributeOrderType] == "LIMIT_SELL" {
		orderType = models.CLOBOrderTypeLimitSell
	}

	return &models.CLOBOrder{
		OrderID:         orderID,
		MarketID:        marketID,
		Counter:         counter,
		Owner:           attrs[clobtypes.AttributeOwner],
		SubAccountID:    subAccountID,
		OrderType:       orderType,
		Price:           price,
		Amount:          amount,
		FilledAmount:    decimal.Zero,
		RemainingAmount: amount,
		Status:          models.CLOBOrderStatusPending,
		TimeInForce:     attrs["time_in_force"],
		PostOnly:        attrs["post_only"] == "true",
		ReduceOnly:      attrs["reduce_only"] == "true",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		BlockHeight:     height,
		TxHash:          txHash,
	}, nil
}

func (p *CLOBParser) parsePlaceMarketOrder(event abci.Event, height int64, txHash string) (*models.CLOBOrder, error) {
	attrs := getAttributeMap(event.Attributes)

	marketID, err := strconv.ParseUint(attrs[clobtypes.AttributeMarketId], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse market id: %w", err)
	}

	orderID, err := strconv.ParseUint(attrs[clobtypes.AttributeOrderId], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse order id: %w", err)
	}

	counter, err := strconv.ParseUint(attrs["counter"], 10, 64)
	if err != nil {
		counter = orderID
	}

	subAccountID, _ := strconv.ParseUint(attrs["sub_account_id"], 10, 64)

	amount, err := decimal.NewFromString(attrs[clobtypes.AttributeQuantity])
	if err != nil {
		return nil, fmt.Errorf("failed to parse quantity: %w", err)
	}

	orderType := models.CLOBOrderTypeMarketBuy
	if attrs[clobtypes.AttributeOrderType] == "MARKET_SELL" {
		orderType = models.CLOBOrderTypeMarketSell
	}

	return &models.CLOBOrder{
		OrderID:         orderID,
		MarketID:        marketID,
		Counter:         counter,
		Owner:           attrs[clobtypes.AttributeOwner],
		SubAccountID:    subAccountID,
		OrderType:       orderType,
		Price:           decimal.Zero, // Market orders don't have a fixed price
		Amount:          amount,
		FilledAmount:    decimal.Zero,
		RemainingAmount: amount,
		Status:          models.CLOBOrderStatusPending,
		TimeInForce:     "IOC", // Market orders are immediate or cancel
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		BlockHeight:     height,
		TxHash:          txHash,
	}, nil
}

func (p *CLOBParser) parseCancelOrder(event abci.Event, height int64, txHash string) (*CLOBOrderUpdate, error) {
	attrs := getAttributeMap(event.Attributes)

	orderID, err := strconv.ParseUint(attrs[clobtypes.AttributeOrderId], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse order id: %w", err)
	}

	return &CLOBOrderUpdate{
		OrderID:     orderID,
		Status:      models.CLOBOrderStatusCancelled,
		CancelledAt: timePtr(time.Now()),
		BlockHeight: height,
		TxHash:      txHash,
	}, nil
}

func (p *CLOBParser) parseOrderExecuted(event abci.Event, height int64, txHash string) (*CLOBOrderExecution, error) {
	attrs := getAttributeMap(event.Attributes)

	orderID, err := strconv.ParseUint(attrs[clobtypes.AttributeOrderId], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse order id: %w", err)
	}

	filledQty, err := decimal.NewFromString(attrs[clobtypes.AttributeFilledQuantity])
	if err != nil {
		return nil, fmt.Errorf("failed to parse filled quantity: %w", err)
	}

	totalFilled, err := decimal.NewFromString(attrs[clobtypes.AttributeTotalFilled])
	if err != nil {
		totalFilled = filledQty
	}

	remainingQty, err := decimal.NewFromString(attrs[clobtypes.AttributeRemainingQuantity])
	if err != nil {
		remainingQty = decimal.Zero
	}

	status := models.CLOBOrderStatusPartiallyFilled
	if remainingQty.IsZero() {
		status = models.CLOBOrderStatusFilled
	}

	executionPrice, _ := decimal.NewFromString(attrs[clobtypes.AttributeExecutionPrice])

	return &CLOBOrderExecution{
		OrderID:        orderID,
		FilledQuantity: filledQty,
		TotalFilled:    totalFilled,
		RemainingQty:   remainingQty,
		ExecutionPrice: executionPrice,
		Status:         status,
		ExecutedAt:     time.Now(),
		BlockHeight:    height,
		TxHash:         txHash,
	}, nil
}

func (p *CLOBParser) parseTrade(event abci.Event, height int64, txHash string) (*models.CLOBTrade, error) {
	attrs := getAttributeMap(event.Attributes)

	marketID, err := strconv.ParseUint(attrs[clobtypes.AttributeMarketId], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse market id: %w", err)
	}

	tradeID, _ := strconv.ParseUint(attrs[clobtypes.AttributeTradeId], 10, 64)

	buyerSubAccountID, _ := strconv.ParseUint(attrs["buyer_sub_account_id"], 10, 64)
	sellerSubAccountID, _ := strconv.ParseUint(attrs["seller_sub_account_id"], 10, 64)

	buyerOrderID, _ := strconv.ParseUint(attrs["buyer_order_id"], 10, 64)
	sellerOrderID, _ := strconv.ParseUint(attrs["seller_order_id"], 10, 64)

	price, err := decimal.NewFromString(attrs[clobtypes.AttributeTradePrice])
	if err != nil {
		return nil, fmt.Errorf("failed to parse trade price: %w", err)
	}

	quantity, err := decimal.NewFromString(attrs[clobtypes.AttributeTradeQuantity])
	if err != nil {
		return nil, fmt.Errorf("failed to parse trade quantity: %w", err)
	}

	tradeValue := price.Mul(quantity)

	buyerFee, _ := decimal.NewFromString(attrs["buyer_fee"])
	sellerFee, _ := decimal.NewFromString(attrs["seller_fee"])

	return &models.CLOBTrade{
		TradeID:             tradeID,
		MarketID:            marketID,
		Buyer:               attrs[clobtypes.AttributeBuyer],
		BuyerSubAccountID:   buyerSubAccountID,
		Seller:              attrs[clobtypes.AttributeSeller],
		SellerSubAccountID:  sellerSubAccountID,
		BuyerOrderID:        buyerOrderID,
		SellerOrderID:       sellerOrderID,
		Price:               price,
		Quantity:            quantity,
		TradeValue:          tradeValue,
		BuyerFee:            buyerFee,
		SellerFee:           sellerFee,
		IsBuyerTaker:        attrs[clobtypes.AttributeIsTaker] == "true",
		IsBuyerLiquidation:  attrs["is_buyer_liquidation"] == "true",
		IsSellerLiquidation: attrs["is_seller_liquidation"] == "true",
		ExecutedAt:          time.Now(),
		BlockHeight:         height,
		TxHash:              txHash,
	}, nil
}

func (p *CLOBParser) parsePositionOpened(event abci.Event, height int64, txHash string) (*models.CLOBPosition, error) {
	attrs := getAttributeMap(event.Attributes)

	marketID, err := strconv.ParseUint(attrs[clobtypes.AttributeMarketId], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse market id: %w", err)
	}

	positionID, err := strconv.ParseUint(attrs[clobtypes.AttributePositionId], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse position id: %w", err)
	}

	subAccountID, _ := strconv.ParseUint(attrs["sub_account_id"], 10, 64)

	size, err := decimal.NewFromString(attrs["size"])
	if err != nil {
		return nil, fmt.Errorf("failed to parse size: %w", err)
	}

	entryPrice, err := decimal.NewFromString(attrs["entry_price"])
	if err != nil {
		return nil, fmt.Errorf("failed to parse entry price: %w", err)
	}

	margin, err := decimal.NewFromString(attrs["margin"])
	if err != nil {
		return nil, fmt.Errorf("failed to parse margin: %w", err)
	}

	side := models.CLOBPositionSideLong
	if attrs[clobtypes.AttributeSide] == "SHORT" {
		side = models.CLOBPositionSideShort
	}

	notional := size.Mul(entryPrice)
	markPrice := entryPrice // Initially, mark price equals entry price
	liquidationPrice, _ := decimal.NewFromString(attrs[clobtypes.AttributeLiquidationPrice])
	marginRatio, _ := decimal.NewFromString(attrs[clobtypes.AttributeMarginRatio])

	return &models.CLOBPosition{
		PositionID:       positionID,
		MarketID:         marketID,
		Owner:            attrs[clobtypes.AttributeOwner],
		SubAccountID:     subAccountID,
		Side:             side,
		Size:             size,
		Notional:         notional,
		EntryPrice:       entryPrice,
		MarkPrice:        markPrice,
		LiquidationPrice: liquidationPrice,
		Margin:           margin,
		MarginRatio:      marginRatio,
		UnrealizedPnL:    decimal.Zero,
		RealizedPnL:      decimal.Zero,
		OpenedAt:         time.Now(),
		UpdatedAt:        time.Now(),
		BlockHeight:      height,
		TxHash:           txHash,
	}, nil
}

func (p *CLOBParser) parsePositionClosed(event abci.Event, height int64, txHash string) (*CLOBPositionUpdate, error) {
	attrs := getAttributeMap(event.Attributes)

	positionID, err := strconv.ParseUint(attrs[clobtypes.AttributePositionId], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse position id: %w", err)
	}

	closePrice, _ := decimal.NewFromString(attrs["close_price"])
	realizedPnL, _ := decimal.NewFromString(attrs[clobtypes.AttributePnL])

	return &CLOBPositionUpdate{
		PositionID:  positionID,
		Action:      "closed",
		ClosedAt:    timePtr(time.Now()),
		ClosePrice:  &closePrice,
		RealizedPnL: &realizedPnL,
		BlockHeight: height,
		TxHash:      txHash,
	}, nil
}

func (p *CLOBParser) parsePositionModified(event abci.Event, height int64, txHash string) (*CLOBPositionUpdate, error) {
	attrs := getAttributeMap(event.Attributes)

	positionID, err := strconv.ParseUint(attrs[clobtypes.AttributePositionId], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse position id: %w", err)
	}

	update := &CLOBPositionUpdate{
		PositionID:  positionID,
		Action:      "modified",
		UpdatedAt:   time.Now(),
		BlockHeight: height,
		TxHash:      txHash,
	}

	// Parse optional fields that might be updated
	if newSize, err := decimal.NewFromString(attrs["new_size"]); err == nil {
		update.NewSize = &newSize
	}
	if newMargin, err := decimal.NewFromString(attrs["new_margin"]); err == nil {
		update.NewMargin = &newMargin
	}
	if newMarkPrice, err := decimal.NewFromString(attrs[clobtypes.AttributeMarkPrice]); err == nil {
		update.NewMarkPrice = &newMarkPrice
	}
	if newLiquidationPrice, err := decimal.NewFromString(attrs[clobtypes.AttributeLiquidationPrice]); err == nil {
		update.NewLiquidationPrice = &newLiquidationPrice
	}
	if newUnrealizedPnL, err := decimal.NewFromString(attrs["unrealized_pnl"]); err == nil {
		update.NewUnrealizedPnL = &newUnrealizedPnL
	}

	return update, nil
}

func (p *CLOBParser) parseLiquidation(event abci.Event, height int64, txHash string) (*models.CLOBLiquidation, error) {
	attrs := getAttributeMap(event.Attributes)

	marketID, err := strconv.ParseUint(attrs[clobtypes.AttributeMarketId], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse market id: %w", err)
	}

	positionID, err := strconv.ParseUint(attrs[clobtypes.AttributePositionId], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse position id: %w", err)
	}

	subAccountID, _ := strconv.ParseUint(attrs["sub_account_id"], 10, 64)

	size, err := decimal.NewFromString(attrs["size"])
	if err != nil {
		return nil, fmt.Errorf("failed to parse size: %w", err)
	}

	price, err := decimal.NewFromString(attrs["liquidation_price"])
	if err != nil {
		return nil, fmt.Errorf("failed to parse liquidation price: %w", err)
	}

	liquidationFee, _ := decimal.NewFromString(attrs["liquidation_fee"])
	insuranceFund, _ := decimal.NewFromString(attrs["insurance_fund_contribution"])

	side := models.CLOBPositionSideLong
	if attrs[clobtypes.AttributeSide] == "SHORT" {
		side = models.CLOBPositionSideShort
	}

	var liquidator *string
	if liquidatorAddr := attrs[clobtypes.AttributeLiquidator]; liquidatorAddr != "" {
		liquidator = &liquidatorAddr
	}

	return &models.CLOBLiquidation{
		MarketID:                  marketID,
		PositionID:                positionID,
		Owner:                     attrs[clobtypes.AttributeOwner],
		SubAccountID:              subAccountID,
		Liquidator:                liquidator,
		Side:                      side,
		Size:                      size,
		Price:                     price,
		LiquidationFee:            liquidationFee,
		InsuranceFundContribution: insuranceFund,
		IsADL:                     attrs["is_adl"] == "true",
		LiquidatedAt:              time.Now(),
		BlockHeight:               height,
		TxHash:                    txHash,
	}, nil
}

func (p *CLOBParser) parseFundingRateUpdate(event abci.Event, height int64, txHash string) (*models.CLOBFundingRate, error) {
	attrs := getAttributeMap(event.Attributes)

	marketID, err := strconv.ParseUint(attrs[clobtypes.AttributeMarketId], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse market id: %w", err)
	}

	fundingRate, err := decimal.NewFromString(attrs[clobtypes.AttributeFundingRate])
	if err != nil {
		return nil, fmt.Errorf("failed to parse funding rate: %w", err)
	}

	premiumRate, _ := decimal.NewFromString(attrs["premium_rate"])
	markPrice, _ := decimal.NewFromString(attrs[clobtypes.AttributeMarkPrice])
	indexPrice, _ := decimal.NewFromString(attrs[clobtypes.AttributeIndexPrice])

	nextFundingTime, _ := time.Parse(time.RFC3339, attrs["next_funding_time"])
	if nextFundingTime.IsZero() {
		nextFundingTime = time.Now().Add(8 * time.Hour) // Default 8-hour funding interval
	}

	return &models.CLOBFundingRate{
		MarketID:        marketID,
		FundingRate:     fundingRate,
		PremiumRate:     premiumRate,
		MarkPrice:       markPrice,
		IndexPrice:      indexPrice,
		Timestamp:       time.Now(),
		NextFundingTime: nextFundingTime,
		BlockHeight:     height,
	}, nil
}

func (p *CLOBParser) parseCreateMarket(event abci.Event, height int64, txHash string) (*models.CLOBMarket, error) {
	attrs := getAttributeMap(event.Attributes)

	marketID, err := strconv.ParseUint(attrs[clobtypes.AttributeMarketId], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse market id: %w", err)
	}

	tickSize, _ := decimal.NewFromString(attrs["tick_size"])
	lotSize, _ := decimal.NewFromString(attrs["lot_size"])
	minOrderSize, _ := decimal.NewFromString(attrs["min_order_size"])
	maxOrderSize, _ := decimal.NewFromString(attrs["max_order_size"])
	maxLeverage, _ := decimal.NewFromString(attrs["max_leverage"])
	initialMarginFraction, _ := decimal.NewFromString(attrs["initial_margin_fraction"])
	maintenanceMarginFraction, _ := decimal.NewFromString(attrs["maintenance_margin_fraction"])
	fundingInterval, _ := strconv.ParseInt(attrs["funding_interval"], 10, 64)

	return &models.CLOBMarket{
		MarketID:                  marketID,
		Ticker:                    attrs["ticker"],
		BaseAsset:                 attrs["base_asset"],
		QuoteAsset:                attrs["quote_asset"],
		TickSize:                  tickSize,
		LotSize:                   lotSize,
		MinOrderSize:              minOrderSize,
		MaxOrderSize:              maxOrderSize,
		MaxLeverage:               maxLeverage,
		InitialMarginFraction:     initialMarginFraction,
		MaintenanceMarginFraction: maintenanceMarginFraction,
		FundingInterval:           fundingInterval,
		NextFundingTime:           time.Now().Add(time.Duration(fundingInterval) * time.Second),
		IsActive:                  true,
		CreatedAt:                 time.Now(),
		UpdatedAt:                 time.Now(),
		BlockHeight:               height,
	}, nil
}

// Helper types for updates
type CLOBOrderUpdate struct {
	OrderID     uint64
	Status      models.CLOBOrderStatus
	ExecutedAt  *time.Time
	CancelledAt *time.Time
	BlockHeight int64
	TxHash      string
}

type CLOBOrderExecution struct {
	OrderID        uint64
	FilledQuantity decimal.Decimal
	TotalFilled    decimal.Decimal
	RemainingQty   decimal.Decimal
	ExecutionPrice decimal.Decimal
	Status         models.CLOBOrderStatus
	ExecutedAt     time.Time
	BlockHeight    int64
	TxHash         string
}

type CLOBPositionUpdate struct {
	PositionID          uint64
	Action              string
	NewSize             *decimal.Decimal
	NewMargin           *decimal.Decimal
	NewMarkPrice        *decimal.Decimal
	NewLiquidationPrice *decimal.Decimal
	NewUnrealizedPnL    *decimal.Decimal
	RealizedPnL         *decimal.Decimal
	ClosedAt            *time.Time
	ClosePrice          *decimal.Decimal
	UpdatedAt           time.Time
	BlockHeight         int64
	TxHash              string
}

func getAttributeMap(attrs []abci.EventAttribute) map[string]string {
	m := make(map[string]string)
	for _, attr := range attrs {
		m[string(attr.Key)] = string(attr.Value)
	}
	return m
}

func timePtr(t time.Time) *time.Time {
	return &t
}
