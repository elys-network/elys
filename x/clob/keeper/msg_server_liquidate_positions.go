package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/clob/types"
)

func (k Keeper) LiquidatePositions(goCtx context.Context, msg *types.MsgLiquidatePositions) (*types.MsgLiquidatePositionsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	liquidator := sdk.MustAccAddressFromBech32(msg.Liquidator)
	liquidatorReward := sdk.Coins{}
	for _, position := range msg.Positions {
		perpetual, err := k.GetPerpetual(ctx, position.MarketId, position.PerpetualId)
		if err != nil {
			return nil, err
		}

		market, err := k.GetPerpetualMarket(ctx, position.MarketId)
		if err != nil {
			return nil, err
		}

		cachedCtx, write := ctx.CacheContext()
		liquidatorRewardAmount, err := k.ForcedLiquidation(cachedCtx, perpetual, market, liquidator)
		if err == nil {
			if !liquidatorRewardAmount.IsZero() {
				liquidatorRewardCoin := sdk.NewCoin(market.QuoteDenom, liquidatorRewardAmount)
				liquidatorReward = liquidatorReward.Add(liquidatorRewardCoin)
			}
			write()
		} else {
			ctx.Logger().Error(fmt.Sprintf("Error liquidating position: Address:%s Id:%d cannot be liquidated due to err: %s", perpetual.GetOwner(), perpetual.Id, err.Error()))
		}
	}

	return &types.MsgLiquidatePositionsResponse{
		LiquidatorReward: liquidatorReward,
	}, nil
}
