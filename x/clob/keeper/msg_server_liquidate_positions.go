package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v7/x/clob/types"
)

func (k Keeper) LiquidatePositions(goCtx context.Context, msg *types.MsgLiquidatePositions) (*types.MsgLiquidatePositionsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	liquidator, err := sdk.AccAddressFromBech32(msg.Liquidator)
	if err != nil {
		return nil, err
	}
	liquidatorReward := sdk.Coins{}
	successCount := 0
	failureCount := 0

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
			successCount++
			if !liquidatorRewardAmount.IsZero() {
				liquidatorRewardCoin := sdk.NewCoin(market.QuoteDenom, liquidatorRewardAmount)
				liquidatorReward = liquidatorReward.Add(liquidatorRewardCoin)
			}

			// Emit individual liquidation event
			ctx.EventManager().EmitEvent(
				sdk.NewEvent(
					types.EventLiquidation,
					sdk.NewAttribute(types.AttributeLiquidator, msg.Liquidator),
					sdk.NewAttribute(types.AttributeOwner, perpetual.GetOwner()),
					sdk.NewAttribute(types.AttributeMarketId, fmt.Sprintf("%d", perpetual.MarketId)),
					sdk.NewAttribute(types.AttributePositionId, fmt.Sprintf("%d", perpetual.Id)),
					sdk.NewAttribute(types.AttributeReward, liquidatorRewardAmount.String()),
				),
			)
		} else {
			failureCount++
			ctx.Logger().Error(fmt.Sprintf("Error liquidating position: Address:%s Id:%d cannot be liquidated due to err: %s", perpetual.GetOwner(), perpetual.Id, err.Error()))
		}
	}

	return &types.MsgLiquidatePositionsResponse{
		LiquidatorReward: liquidatorReward,
	}, nil
}
