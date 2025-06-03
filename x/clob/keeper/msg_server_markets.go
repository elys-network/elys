package keeper

import (
	"context"
	"cosmossdk.io/math"
	"errors"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/clob/types"
)

func (k Keeper) CreatePerpetualMarket(goCtx context.Context, msg *types.MsgCreatPerpetualMarket) (*types.MsgCreatPerpetualMarketResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	marketExists := k.CheckPerpetualMarketAlreadyExists(ctx, msg.BaseDenom, msg.QuoteDenom)
	if marketExists {
		return nil, errors.New("perpetual market already exists")
	}

	_, found := k.oracleKeeper.GetAssetInfo(ctx, msg.BaseDenom)
	if !found {
		return nil, fmt.Errorf("asset info for %s not found", msg.BaseDenom)
	}

	_, found = k.oracleKeeper.GetAssetInfo(ctx, msg.QuoteDenom)
	if !found {
		return nil, fmt.Errorf("asset info for %s not found", msg.QuoteDenom)
	}

	newMarketId := k.CountAllPerpetualMarket(ctx) + 1
	perpetualMarket := types.PerpetualMarket{
		Id:                      newMarketId,
		BaseDenom:               msg.BaseDenom,
		QuoteDenom:              msg.QuoteDenom,
		InitialMarginRatio:      msg.InitialMarginRatio,
		MaintenanceMarginRatio:  msg.MaintenanceMarginRatio,
		MakerFeeRate:            msg.MakerFeeRate,
		TakerFeeRate:            msg.TakerFeeRate,
		LiquidationFeeShareRate: msg.LiquidationFeeShareRate,
		Status:                  types.PerpetualMarketStatus_MARKET_STATUS_ACTIVE,
		MinPriceTickSize:        msg.MinPriceTickSize,
		MinQuantityTickSize:     msg.MinQuantityTickSize,
		MinNotional:             msg.MinNotional,
		Admin:                   msg.Creator,
		TotalOpen:               math.LegacyZeroDec(),
		MaxAbsFundingRate:       msg.MaxAbsFundingRate,
		MaxAbsFundingRateChange: msg.MaxAbsFundingRateChange,
		TwapPricesWindow:        msg.TwapPricesWindow,
	}
	k.SetPerpetualMarket(ctx, perpetualMarket)

	k.SetPerpetualMarketCounter(ctx, types.PerpetualMarketCounter{
		MarketId:          perpetualMarket.Id,
		OrderCounter:      0,
		PerpetualCounter:  0,
		TotalOpenPosition: 0,
		TotalOpenOrders:   0,
	})

	return &types.MsgCreatPerpetualMarketResponse{}, nil
}
