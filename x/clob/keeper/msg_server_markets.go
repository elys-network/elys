package keeper

import (
	"context"
	"cosmossdk.io/math"
	"errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/clob/types"
)

func (k Keeper) CreatePerpetualMarket(goCtx context.Context, msg *types.MsgCreatPerpetualMarket) (*types.MsgCreatPerpetualMarketResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	marketExists := k.CheckPerpetualMarketAlreadyExists(ctx, msg.BaseDenom, msg.QuoteDenom)
	if marketExists {
		return nil, errors.New("perpetual market already exists")
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
		AllowedCollateral:       msg.AllowedCollateral,
		TotalOpen:               math.ZeroInt(),
		MaxFundingRateChange:    math.NewDecWithExp(1, 3),
		MaxFundingRate:          math.NewDecWithExp(2, 2),
	}

	k.SetPerpetualMarket(ctx, perpetualMarket)

	return &types.MsgCreatPerpetualMarketResponse{}, nil
}
