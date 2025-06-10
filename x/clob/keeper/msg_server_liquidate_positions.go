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

		// Fetch the market after each close since the exchange updates the total open interest (total open) of the market.
		market, err := k.GetPerpetualMarket(ctx, position.MarketId)
		if err != nil {
			return nil, err
		}

		liquidatorRewardAmount, err := k.ForcedLiquidation(ctx, perpetual, market, liquidator)
		if err == nil {
			if !liquidatorRewardAmount.IsZero() {
				liquidatorRewardCoin := sdk.NewCoin(market.QuoteDenom, liquidatorRewardAmount)
				liquidatorReward = liquidatorReward.Add(liquidatorRewardCoin)
			}
		} else {
			ctx.Logger().Error(fmt.Sprintf("Error liquidating position: Address:%s Id:%d cannot be liquidated due to err: %s", perpetual.GetOwner(), perpetual.Id, err.Error()))
		}
	}

	return &types.MsgLiquidatePositionsResponse{
		LiquidatorReward: liquidatorReward,
	}, nil
}
