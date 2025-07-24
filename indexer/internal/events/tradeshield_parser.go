package events

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"cosmossdk.io/math"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/elys-network/elys/indexer/internal/models"
	tradeshieldtypes "github.com/elys-network/elys/v6/x/tradeshield/types"
	"go.uber.org/zap"
)

type TradeShieldParser struct {
	logger *zap.Logger
}

func NewTradeShieldParser(logger *zap.Logger) *TradeShieldParser {
	return &TradeShieldParser{
		logger: logger,
	}
}

func (p *TradeShieldParser) ParseEvents(ctx context.Context, events []abci.Event, blockHeight int64, txHash string) ([]interface{}, error) {
	var results []interface{}

	for _, event := range events {
		switch event.Type {
		case tradeshieldtypes.TypeEvtCreatePerpetualOpenOrder:
			order, err := p.parseCreatePerpetualOpenOrder(event, blockHeight, txHash)
			if err != nil {
				p.logger.Error("failed to parse create perpetual open order", zap.Error(err))
				continue
			}
			results = append(results, order)

		case tradeshieldtypes.TypeEvtCreatePerpetualCloseOrder:
			order, err := p.parseCreatePerpetualCloseOrder(event, blockHeight, txHash)
			if err != nil {
				p.logger.Error("failed to parse create perpetual close order", zap.Error(err))
				continue
			}
			results = append(results, order)

		case tradeshieldtypes.TypeEvtExecuteLimitOpenPerpetualOrder:
			update, err := p.parseExecuteLimitOpenPerpetualOrder(event, blockHeight, txHash)
			if err != nil {
				p.logger.Error("failed to parse execute limit open perpetual order", zap.Error(err))
				continue
			}
			results = append(results, update)

		case tradeshieldtypes.TypeEvtExecuteLimitClosePerpetualOrder:
			update, err := p.parseExecuteLimitClosePerpetualOrder(event, blockHeight, txHash)
			if err != nil {
				p.logger.Error("failed to parse execute limit close perpetual order", zap.Error(err))
				continue
			}
			results = append(results, update)

		case tradeshieldtypes.TypeEvtCancelPerpetualOrder:
			update, err := p.parseCancelPerpetualOrder(event, blockHeight, txHash)
			if err != nil {
				p.logger.Error("failed to parse cancel perpetual order", zap.Error(err))
				continue
			}
			results = append(results, update)

		case tradeshieldtypes.TypeEvtExecuteLimitBuySpotOrder:
			trade, err := p.parseExecuteLimitBuySpotOrder(event, blockHeight, txHash)
			if err != nil {
				p.logger.Error("failed to parse execute limit buy spot order", zap.Error(err))
				continue
			}
			results = append(results, trade)

		case tradeshieldtypes.TypeEvtExecuteLimitSellSpotOrder:
			trade, err := p.parseExecuteLimitSellSpotOrder(event, blockHeight, txHash)
			if err != nil {
				p.logger.Error("failed to parse execute limit sell spot order", zap.Error(err))
				continue
			}
			results = append(results, trade)

		case tradeshieldtypes.TypeEvtExecuteMarketBuySpotOrder:
			trade, err := p.parseExecuteMarketBuySpotOrder(event, blockHeight, txHash)
			if err != nil {
				p.logger.Error("failed to parse execute market buy spot order", zap.Error(err))
				continue
			}
			results = append(results, trade)

		case tradeshieldtypes.TypeEvtExecuteStopLossSpotOrder:
			trade, err := p.parseExecuteStopLossSpotOrder(event, blockHeight, txHash)
			if err != nil {
				p.logger.Error("failed to parse execute stop loss spot order", zap.Error(err))
				continue
			}
			results = append(results, trade)

		case tradeshieldtypes.TypeEvtCloseSpotOrder:
			update, err := p.parseCloseSpotOrder(event, blockHeight, txHash)
			if err != nil {
				p.logger.Error("failed to parse close spot order", zap.Error(err))
				continue
			}
			results = append(results, update)
		}
	}

	return results, nil
}

func (p *TradeShieldParser) parseCreatePerpetualOpenOrder(event abci.Event, blockHeight int64, txHash string) (*models.PerpetualOrder, error) {
	attrs := parseEventAttributes(event.Attributes)

	orderID, err := strconv.ParseUint(attrs["order_id"], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse order_id: %w", err)
	}

	order := &models.PerpetualOrder{
		OrderID:            orderID,
		OwnerAddress:       attrs["owner_address"],
		PerpetualOrderType: models.PerpetualOrderTypeLimitOpen,
		Position:           models.PositionType(attrs["position"]),
		Status:             models.OrderStatusPending,
		CreatedAt:          time.Now(),
		BlockHeight:        blockHeight,
		TxHash:             txHash,
	}

	// Parse order type
	if orderType := attrs["order_type"]; orderType != "" {
		switch orderType {
		case "LIMIT_OPEN":
			order.PerpetualOrderType = models.PerpetualOrderTypeLimitOpen
		case "LIMIT_CLOSE":
			order.PerpetualOrderType = models.PerpetualOrderTypeLimitClose
		case "MARKET":
			order.PerpetualOrderType = models.PerpetualOrderTypeMarket
		}
	}

	return order, nil
}

func (p *TradeShieldParser) parseCreatePerpetualCloseOrder(event abci.Event, blockHeight int64, txHash string) (*models.PerpetualOrder, error) {
	attrs := parseEventAttributes(event.Attributes)

	orderID, err := strconv.ParseUint(attrs["order_id"], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse order_id: %w", err)
	}

	order := &models.PerpetualOrder{
		OrderID:            orderID,
		OwnerAddress:       attrs["owner_address"],
		PerpetualOrderType: models.PerpetualOrderTypeLimitClose,
		Status:             models.OrderStatusPending,
		CreatedAt:          time.Now(),
		BlockHeight:        blockHeight,
		TxHash:             txHash,
	}

	return order, nil
}

func (p *TradeShieldParser) parseExecuteLimitOpenPerpetualOrder(event abci.Event, blockHeight int64, txHash string) (*PerpetualOrderUpdate, error) {
	attrs := parseEventAttributes(event.Attributes)

	orderID, err := strconv.ParseUint(attrs["order_id"], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse order_id: %w", err)
	}

	positionID, err := strconv.ParseUint(attrs["position_id"], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse position_id: %w", err)
	}

	return &PerpetualOrderUpdate{
		OrderID:    orderID,
		Status:     models.OrderStatusExecuted,
		PositionID: &positionID,
	}, nil
}

func (p *TradeShieldParser) parseExecuteLimitClosePerpetualOrder(event abci.Event, blockHeight int64, txHash string) (*PerpetualOrderUpdate, error) {
	attrs := parseEventAttributes(event.Attributes)

	orderID, err := strconv.ParseUint(attrs["order_id"], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse order_id: %w", err)
	}

	return &PerpetualOrderUpdate{
		OrderID: orderID,
		Status:  models.OrderStatusExecuted,
	}, nil
}

func (p *TradeShieldParser) parseCancelPerpetualOrder(event abci.Event, blockHeight int64, txHash string) (*PerpetualOrderUpdate, error) {
	attrs := parseEventAttributes(event.Attributes)

	orderID, err := strconv.ParseUint(attrs["id"], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse order_id: %w", err)
	}

	return &PerpetualOrderUpdate{
		OrderID: orderID,
		Status:  models.OrderStatusCancelled,
	}, nil
}

func (p *TradeShieldParser) parseExecuteLimitBuySpotOrder(event abci.Event, blockHeight int64, txHash string) (*SpotOrderExecution, error) {
	attrs := parseEventAttributes(event.Attributes)

	orderID, err := strconv.ParseUint(attrs["order_id"], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse order_id: %w", err)
	}

	amount, ok := math.NewIntFromString(attrs["amount"])
	if !ok {
		return nil, fmt.Errorf("failed to parse amount")
	}

	spotPrice, err := math.LegacyNewDecFromStr(attrs["spot_price"])
	if err != nil {
		return nil, fmt.Errorf("failed to parse spot_price: %w", err)
	}

	swapFee, err := math.LegacyNewDecFromStr(attrs["swap_fee"])
	if err != nil {
		return nil, fmt.Errorf("failed to parse swap_fee: %w", err)
	}

	// Parse order price JSON
	var orderPrice models.JSONB
	if err := json.Unmarshal([]byte(attrs["order_price"]), &orderPrice); err != nil {
		p.logger.Warn("failed to parse order_price JSON", zap.Error(err))
		orderPrice = make(models.JSONB)
	}

	return &SpotOrderExecution{
		Order: &models.SpotOrder{
			OrderID:          orderID,
			OrderType:        models.OrderTypeLimitBuy,
			OwnerAddress:     attrs["owner_address"],
			OrderTargetDenom: attrs["order_target_denom"],
			OrderPrice:       orderPrice,
			OrderAmount:      amount,
			Status:           models.OrderStatusExecuted,
			ExecutedAt:       nil,
			BlockHeight:      blockHeight,
			TxHash:           txHash,
		},
		Trade: &models.Trade{
			TradeType:    "spot",
			ReferenceID:  orderID,
			OwnerAddress: attrs["owner_address"],
			Asset:        attrs["order_target_denom"],
			Amount:       amount,
			Price:        spotPrice,
			Fees: models.JSONB{
				"swap_fee": swapFee.String(),
				"discount": attrs["discount"],
			},
			ExecutedAt:  time.Now(),
			BlockHeight: blockHeight,
			TxHash:      txHash,
			EventType:   tradeshieldtypes.TypeEvtExecuteLimitBuySpotOrder,
		},
	}, nil
}

func (p *TradeShieldParser) parseExecuteLimitSellSpotOrder(event abci.Event, blockHeight int64, txHash string) (*SpotOrderExecution, error) {
	// Similar to buy order but with sell type
	attrs := parseEventAttributes(event.Attributes)

	orderID, err := strconv.ParseUint(attrs["order_id"], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse order_id: %w", err)
	}

	amount, ok := math.NewIntFromString(attrs["amount"])
	if !ok {
		return nil, fmt.Errorf("failed to parse amount")
	}

	spotPrice, err := math.LegacyNewDecFromStr(attrs["spot_price"])
	if err != nil {
		return nil, fmt.Errorf("failed to parse spot_price: %w", err)
	}

	swapFee, err := math.LegacyNewDecFromStr(attrs["swap_fee"])
	if err != nil {
		return nil, fmt.Errorf("failed to parse swap_fee: %w", err)
	}

	var orderPrice models.JSONB
	if err := json.Unmarshal([]byte(attrs["order_price"]), &orderPrice); err != nil {
		p.logger.Warn("failed to parse order_price JSON", zap.Error(err))
		orderPrice = make(models.JSONB)
	}

	return &SpotOrderExecution{
		Order: &models.SpotOrder{
			OrderID:          orderID,
			OrderType:        models.OrderTypeLimitSell,
			OwnerAddress:     attrs["owner_address"],
			OrderTargetDenom: attrs["order_target_denom"],
			OrderPrice:       orderPrice,
			OrderAmount:      amount,
			Status:           models.OrderStatusExecuted,
			ExecutedAt:       nil,
			BlockHeight:      blockHeight,
			TxHash:           txHash,
		},
		Trade: &models.Trade{
			TradeType:    "spot",
			ReferenceID:  orderID,
			OwnerAddress: attrs["owner_address"],
			Asset:        attrs["order_target_denom"],
			Amount:       amount,
			Price:        spotPrice,
			Fees: models.JSONB{
				"swap_fee": swapFee.String(),
				"discount": attrs["discount"],
			},
			ExecutedAt:  time.Now(),
			BlockHeight: blockHeight,
			TxHash:      txHash,
			EventType:   tradeshieldtypes.TypeEvtExecuteLimitSellSpotOrder,
		},
	}, nil
}

func (p *TradeShieldParser) parseExecuteMarketBuySpotOrder(event abci.Event, blockHeight int64, txHash string) (*SpotOrderExecution, error) {
	attrs := parseEventAttributes(event.Attributes)

	orderID, err := strconv.ParseUint(attrs["order_id"], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse order_id: %w", err)
	}

	amount, ok := math.NewIntFromString(attrs["amount"])
	if !ok {
		return nil, fmt.Errorf("failed to parse amount")
	}

	spotPrice, err := math.LegacyNewDecFromStr(attrs["spot_price"])
	if err != nil {
		return nil, fmt.Errorf("failed to parse spot_price: %w", err)
	}

	swapFee, err := math.LegacyNewDecFromStr(attrs["swap_fee"])
	if err != nil {
		return nil, fmt.Errorf("failed to parse swap_fee: %w", err)
	}

	return &SpotOrderExecution{
		Order: &models.SpotOrder{
			OrderID:          orderID,
			OrderType:        models.OrderTypeMarketBuy,
			OwnerAddress:     attrs["owner_address"],
			OrderTargetDenom: attrs["order_target_denom"],
			OrderAmount:      amount,
			Status:           models.OrderStatusExecuted,
			ExecutedAt:       nil,
			BlockHeight:      blockHeight,
			TxHash:           txHash,
		},
		Trade: &models.Trade{
			TradeType:    "spot",
			ReferenceID:  orderID,
			OwnerAddress: attrs["owner_address"],
			Asset:        attrs["order_target_denom"],
			Amount:       amount,
			Price:        spotPrice,
			Fees: models.JSONB{
				"swap_fee": swapFee.String(),
				"discount": attrs["discount"],
			},
			ExecutedAt:  time.Now(),
			BlockHeight: blockHeight,
			TxHash:      txHash,
			EventType:   tradeshieldtypes.TypeEvtExecuteMarketBuySpotOrder,
		},
	}, nil
}

func (p *TradeShieldParser) parseExecuteStopLossSpotOrder(event abci.Event, blockHeight int64, txHash string) (*SpotOrderExecution, error) {
	attrs := parseEventAttributes(event.Attributes)

	orderID, err := strconv.ParseUint(attrs["order_id"], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse order_id: %w", err)
	}

	amount, ok := math.NewIntFromString(attrs["amount"])
	if !ok {
		return nil, fmt.Errorf("failed to parse amount")
	}

	spotPrice, err := math.LegacyNewDecFromStr(attrs["spot_price"])
	if err != nil {
		return nil, fmt.Errorf("failed to parse spot_price: %w", err)
	}

	swapFee, err := math.LegacyNewDecFromStr(attrs["swap_fee"])
	if err != nil {
		return nil, fmt.Errorf("failed to parse swap_fee: %w", err)
	}

	var orderPrice models.JSONB
	if err := json.Unmarshal([]byte(attrs["order_price"]), &orderPrice); err != nil {
		p.logger.Warn("failed to parse order_price JSON", zap.Error(err))
		orderPrice = make(models.JSONB)
	}

	return &SpotOrderExecution{
		Order: &models.SpotOrder{
			OrderID:          orderID,
			OrderType:        models.OrderTypeStopLoss,
			OwnerAddress:     attrs["owner_address"],
			OrderTargetDenom: attrs["order_target_denom"],
			OrderPrice:       orderPrice,
			OrderAmount:      amount,
			Status:           models.OrderStatusExecuted,
			ExecutedAt:       nil,
			BlockHeight:      blockHeight,
			TxHash:           txHash,
		},
		Trade: &models.Trade{
			TradeType:    "spot",
			ReferenceID:  orderID,
			OwnerAddress: attrs["owner_address"],
			Asset:        attrs["order_target_denom"],
			Amount:       amount,
			Price:        spotPrice,
			Fees: models.JSONB{
				"swap_fee": swapFee.String(),
				"discount": attrs["discount"],
			},
			ExecutedAt:  time.Now(),
			BlockHeight: blockHeight,
			TxHash:      txHash,
			EventType:   tradeshieldtypes.TypeEvtExecuteStopLossSpotOrder,
		},
	}, nil
}

func (p *TradeShieldParser) parseCloseSpotOrder(event abci.Event, blockHeight int64, txHash string) (*SpotOrderUpdate, error) {
	attrs := parseEventAttributes(event.Attributes)

	orderID, err := strconv.ParseUint(attrs["id"], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse order_id: %w", err)
	}

	return &SpotOrderUpdate{
		OrderID: orderID,
		Status:  models.OrderStatusClosed,
	}, nil
}

// Helper function to parse event attributes
func parseEventAttributes(attrs []abci.EventAttribute) map[string]string {
	result := make(map[string]string)
	for _, attr := range attrs {
		result[string(attr.Key)] = string(attr.Value)
	}
	return result
}

// Update types
type PerpetualOrderUpdate struct {
	OrderID    uint64
	Status     models.OrderStatus
	PositionID *uint64
}

type SpotOrderUpdate struct {
	OrderID uint64
	Status  models.OrderStatus
}

type SpotOrderExecution struct {
	Order *models.SpotOrder
	Trade *models.Trade
}
