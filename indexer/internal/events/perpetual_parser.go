package events

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"cosmossdk.io/math"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/elys-network/elys/indexer/internal/models"
	perpetualtypes "github.com/elys-network/elys/v7/x/perpetual/types"
	"go.uber.org/zap"
)

type PerpetualParser struct {
	logger *zap.Logger
}

func NewPerpetualParser(logger *zap.Logger) *PerpetualParser {
	return &PerpetualParser{
		logger: logger,
	}
}

func (p *PerpetualParser) ParseEvents(ctx context.Context, events []abci.Event, blockHeight int64, txHash string) ([]interface{}, error) {
	var results []interface{}

	for _, event := range events {
		switch event.Type {
		case perpetualtypes.EventOpen:
			position, err := p.parseOpenPosition(event, blockHeight, txHash)
			if err != nil {
				p.logger.Error("failed to parse open position", zap.Error(err))
				continue
			}
			results = append(results, position)

		case perpetualtypes.EventClose:
			closeData, err := p.parseClosePosition(event, blockHeight, txHash)
			if err != nil {
				p.logger.Error("failed to parse close position", zap.Error(err))
				continue
			}
			results = append(results, closeData)

		case perpetualtypes.EventForceClosed:
			closeData, err := p.parseForceClosePosition(event, blockHeight, txHash)
			if err != nil {
				p.logger.Error("failed to parse force close position", zap.Error(err))
				continue
			}
			results = append(results, closeData)

		case perpetualtypes.EventUpdateStopLoss:
			update, err := p.parseUpdateStopLoss(event, blockHeight, txHash)
			if err != nil {
				p.logger.Error("failed to parse update stop loss", zap.Error(err))
				continue
			}
			results = append(results, update)

		case perpetualtypes.EventUpdateTakeProfitPrice:
			update, err := p.parseUpdateTakeProfit(event, blockHeight, txHash)
			if err != nil {
				p.logger.Error("failed to parse update take profit", zap.Error(err))
				continue
			}
			results = append(results, update)

		case perpetualtypes.EventAddCollateral:
			update, err := p.parseAddCollateral(event, blockHeight, txHash)
			if err != nil {
				p.logger.Error("failed to parse add collateral", zap.Error(err))
				continue
			}
			results = append(results, update)

		case perpetualtypes.EventOpenConsolidate:
			consolidate, err := p.parseOpenConsolidate(event, blockHeight, txHash)
			if err != nil {
				p.logger.Error("failed to parse open consolidate", zap.Error(err))
				continue
			}
			results = append(results, consolidate)
		}
	}

	return results, nil
}

func (p *PerpetualParser) parseOpenPosition(event abci.Event, blockHeight int64, txHash string) (*models.PerpetualPosition, error) {
	attrs := parseEventAttributes(event.Attributes)

	mtpID, err := strconv.ParseUint(attrs["mtp_id"], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse mtp_id: %w", err)
	}

	ammPoolID, err := strconv.ParseUint(attrs["amm_pool_id"], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse amm_pool_id: %w", err)
	}

	collateral, ok := math.NewIntFromString(attrs["collateral"])
	if !ok {
		return nil, fmt.Errorf("failed to parse collateral")
	}

	liabilities, ok := math.NewIntFromString(attrs["liabilities"])
	if !ok {
		return nil, fmt.Errorf("failed to parse liabilities")
	}

	custody, ok := math.NewIntFromString(attrs["custody"])
	if !ok {
		return nil, fmt.Errorf("failed to parse custody")
	}

	mtpHealth, err := math.LegacyNewDecFromStr(attrs["mtp_health"])
	if err != nil {
		return nil, fmt.Errorf("failed to parse mtp_health: %w", err)
	}

	openPrice, err := math.LegacyNewDecFromStr(attrs["open_price"])
	if err != nil {
		return nil, fmt.Errorf("failed to parse open_price: %w", err)
	}

	stopLossPrice, err := math.LegacyNewDecFromStr(attrs["stop_loss_price"])
	if err != nil {
		return nil, fmt.Errorf("failed to parse stop_loss_price: %w", err)
	}

	takeProfitPrice, err := math.LegacyNewDecFromStr(attrs["take_profit_price"])
	if err != nil {
		return nil, fmt.Errorf("failed to parse take_profit_price: %w", err)
	}

	position := &models.PerpetualPosition{
		MtpID:           mtpID,
		OwnerAddress:    attrs["owner"],
		AmmPoolID:       ammPoolID,
		Position:        models.PositionType(attrs["position"]),
		CollateralAsset: attrs["collateral_asset"],
		Collateral:      collateral,
		Liabilities:     liabilities,
		Custody:         custody,
		MtpHealth:       mtpHealth,
		OpenPrice:       openPrice,
		StopLossPrice:   stopLossPrice,
		TakeProfitPrice: takeProfitPrice,
		OpenedAt:        time.Now(),
		BlockHeight:     blockHeight,
		TxHash:          txHash,
	}

	// Also create a trade record for the opening - this will be handled separately
	// in the indexer when processing the position opening event

	return position, nil
}

func (p *PerpetualParser) parseClosePosition(event abci.Event, blockHeight int64, txHash string) (*PositionCloseData, error) {
	attrs := parseEventAttributes(event.Attributes)

	mtpID, err := strconv.ParseUint(attrs["mtp_id"], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse mtp_id: %w", err)
	}

	closingPrice, err := math.LegacyNewDecFromStr(attrs["closing_price"])
	if err != nil {
		return nil, fmt.Errorf("failed to parse closing_price: %w", err)
	}

	netPnL, err := math.LegacyNewDecFromStr(attrs["net_pnl"])
	if err != nil {
		return nil, fmt.Errorf("failed to parse net_pnl: %w", err)
	}

	returnAmount, ok := math.NewIntFromString(attrs["return_amount"])
	if !ok {
		return nil, fmt.Errorf("failed to parse return_amount")
	}

	closeData := &PositionCloseData{
		MtpID:        mtpID,
		ClosingPrice: closingPrice,
		NetPnL:       netPnL,
		ClosedBy:     attrs["owner"],
		CloseTrigger: "manual_close",
	}

	// Create trade record for the close
	closeData.Trade = &models.Trade{
		TradeType:    "perpetual",
		ReferenceID:  mtpID,
		OwnerAddress: attrs["owner"],
		Asset:        attrs["collateral_asset"],
		Amount:       returnAmount,
		Price:        closingPrice,
		Fees: models.JSONB{
			"perp_fee":            attrs[perpetualtypes.AttributeKeyPerpFee],
			"slippage":            attrs[perpetualtypes.AttributeKeySlippage],
			"weight_breaking_fee": attrs[perpetualtypes.AttributeKeyWeightBreakingFee],
			"taker_fees":          attrs[perpetualtypes.AttributeTakerFees],
		},
		ExecutedAt:  time.Now(),
		BlockHeight: blockHeight,
		TxHash:      txHash,
		EventType:   perpetualtypes.EventClose,
	}

	return closeData, nil
}

func (p *PerpetualParser) parseForceClosePosition(event abci.Event, blockHeight int64, txHash string) (*PositionCloseData, error) {
	attrs := parseEventAttributes(event.Attributes)

	mtpID, err := strconv.ParseUint(attrs["mtp_id"], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse mtp_id: %w", err)
	}

	closingPrice, err := math.LegacyNewDecFromStr(attrs["closing_price"])
	if err != nil {
		return nil, fmt.Errorf("failed to parse closing_price: %w", err)
	}

	netPnL, err := math.LegacyNewDecFromStr(attrs["net_pnl"])
	if err != nil {
		return nil, fmt.Errorf("failed to parse net_pnl: %w", err)
	}

	returnAmount, ok := math.NewIntFromString(attrs["return_amount"])
	if !ok {
		return nil, fmt.Errorf("failed to parse return_amount")
	}

	closeData := &PositionCloseData{
		MtpID:        mtpID,
		ClosingPrice: closingPrice,
		NetPnL:       netPnL,
		ClosedBy:     attrs["closer"],
		CloseTrigger: attrs["trigger"],
	}

	// Create trade record for the force close
	closeData.Trade = &models.Trade{
		TradeType:    "perpetual",
		ReferenceID:  mtpID,
		OwnerAddress: attrs["owner"],
		Asset:        attrs["collateral_asset"],
		Amount:       returnAmount,
		Price:        closingPrice,
		Fees: models.JSONB{
			"perp_fee":            attrs[perpetualtypes.AttributeKeyPerpFee],
			"slippage":            attrs[perpetualtypes.AttributeKeySlippage],
			"weight_breaking_fee": attrs[perpetualtypes.AttributeKeyWeightBreakingFee],
			"taker_fees":          attrs[perpetualtypes.AttributeTakerFees],
			"interest_amount":     attrs["interest_amount"],
			"funding_fee_amount":  attrs["funding_fee_amount"],
		},
		ExecutedAt:  time.Now(),
		BlockHeight: blockHeight,
		TxHash:      txHash,
		EventType:   perpetualtypes.EventForceClosed,
	}

	return closeData, nil
}

func (p *PerpetualParser) parseUpdateStopLoss(event abci.Event, blockHeight int64, txHash string) (*PositionUpdate, error) {
	attrs := parseEventAttributes(event.Attributes)

	mtpID, err := strconv.ParseUint(attrs["mtp_id"], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse mtp_id: %w", err)
	}

	stopLossPrice, err := math.LegacyNewDecFromStr(attrs["stop_loss_price"])
	if err != nil {
		return nil, fmt.Errorf("failed to parse stop_loss_price: %w", err)
	}

	return &PositionUpdate{
		MtpID:         mtpID,
		UpdateType:    "stop_loss",
		StopLossPrice: &stopLossPrice,
	}, nil
}

func (p *PerpetualParser) parseUpdateTakeProfit(event abci.Event, blockHeight int64, txHash string) (*PositionUpdate, error) {
	attrs := parseEventAttributes(event.Attributes)

	mtpID, err := strconv.ParseUint(attrs["mtp_id"], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse mtp_id: %w", err)
	}

	takeProfitPrice, err := math.LegacyNewDecFromStr(attrs["take_profit_price"])
	if err != nil {
		return nil, fmt.Errorf("failed to parse take_profit_price: %w", err)
	}

	return &PositionUpdate{
		MtpID:           mtpID,
		UpdateType:      "take_profit",
		TakeProfitPrice: &takeProfitPrice,
	}, nil
}

func (p *PerpetualParser) parseAddCollateral(event abci.Event, blockHeight int64, txHash string) (*PositionUpdate, error) {
	attrs := parseEventAttributes(event.Attributes)

	mtpID, err := strconv.ParseUint(attrs["mtp_id"], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse mtp_id: %w", err)
	}

	collateral, ok := math.NewIntFromString(attrs["collateral"])
	if !ok {
		return nil, fmt.Errorf("failed to parse collateral")
	}

	mtpHealth, err := math.LegacyNewDecFromStr(attrs["mtp_health"])
	if err != nil {
		return nil, fmt.Errorf("failed to parse mtp_health: %w", err)
	}

	return &PositionUpdate{
		MtpID:      mtpID,
		UpdateType: "add_collateral",
		Collateral: &collateral,
		MtpHealth:  &mtpHealth,
	}, nil
}

func (p *PerpetualParser) parseOpenConsolidate(event abci.Event, blockHeight int64, txHash string) (*PositionConsolidate, error) {
	attrs := parseEventAttributes(event.Attributes)

	oldMtpID, err := strconv.ParseUint(attrs["old_mtp_id"], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse old_mtp_id: %w", err)
	}

	newMtpID, err := strconv.ParseUint(attrs["new_mtp_id"], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse new_mtp_id: %w", err)
	}

	collateral, ok := math.NewIntFromString(attrs["collateral"])
	if !ok {
		return nil, fmt.Errorf("failed to parse collateral")
	}

	custody, ok := math.NewIntFromString(attrs["custody"])
	if !ok {
		return nil, fmt.Errorf("failed to parse custody")
	}

	liabilities, ok := math.NewIntFromString(attrs["liabilities"])
	if !ok {
		return nil, fmt.Errorf("failed to parse liabilities")
	}

	return &PositionConsolidate{
		OldMtpID:    oldMtpID,
		NewMtpID:    newMtpID,
		Collateral:  collateral,
		Custody:     custody,
		Liabilities: liabilities,
	}, nil
}

// Update types specific to perpetual module
type PositionCloseData struct {
	MtpID        uint64
	ClosingPrice math.LegacyDec
	NetPnL       math.LegacyDec
	ClosedBy     string
	CloseTrigger string
	Trade        *models.Trade
}

type PositionUpdate struct {
	MtpID           uint64
	UpdateType      string
	Collateral      *math.Int
	MtpHealth       *math.LegacyDec
	StopLossPrice   *math.LegacyDec
	TakeProfitPrice *math.LegacyDec
}

type PositionConsolidate struct {
	OldMtpID    uint64
	NewMtpID    uint64
	Collateral  math.Int
	Custody     math.Int
	Liabilities math.Int
}
