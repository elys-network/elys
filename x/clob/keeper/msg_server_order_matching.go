package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v7/x/clob/types"
)

func (k Keeper) MatchAndExecuteOrders(goCtx context.Context, msg *types.MsgMatchAndExecuteOrders) (*types.MsgMatchAndExecuteOrdersResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	marketsProcessed := 0
	totalVolume := math.LegacyZeroDec()
	totalVolumeValue := math.LegacyZeroDec()
	totalOrdersProcessed := uint64(0)

	for _, marketId := range msg.MarketIds {
		volumeValue, ordersProcessed, volume, err := k.ExecuteMarket(ctx, marketId, int64(msg.Limit))
		if err != nil {
			return nil, err
		}
		if ordersProcessed > 0 {
			marketsProcessed++
			totalVolume = totalVolume.Add(volume)
			totalVolumeValue = totalVolumeValue.Add(volumeValue)
			totalOrdersProcessed += ordersProcessed
		}
	}

	// Emit event for the batch operation
	if marketsProcessed > 0 {
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventMatchAndExecuteOrders,
				sdk.NewAttribute(types.AttributeSender, msg.Sender),
				sdk.NewAttribute("markets_processed", fmt.Sprintf("%d", marketsProcessed)),
				sdk.NewAttribute("market_ids", fmt.Sprintf("%v", msg.MarketIds)),
				sdk.NewAttribute("total_orders_processed", fmt.Sprintf("%d", totalOrdersProcessed)),
				sdk.NewAttribute("total_volume", totalVolume.String()),
				sdk.NewAttribute("total_volume_value", totalVolumeValue.String()),
			),
		)
	}

	return &types.MsgMatchAndExecuteOrdersResponse{}, nil
}
