package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/leveragelp/types"
)

func (k msgServer) ClosePositions(goCtx context.Context, msg *types.MsgClosePositions) (*types.MsgClosePositionsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Handle liquidations
	for _, val := range msg.Liquidate {
		position, err := k.GetPosition(ctx, val.Address, val.Id)
		if err != nil {
			continue
		}
		pool, found := k.GetPool(ctx, position.AmmPoolId)
		ammPool, err := k.GetAmmPool(ctx, position.AmmPoolId)
		if !found || err != nil {
			continue
		}
		_, _ = k.LiquidatePositionIfUnhealthy(ctx, &position, pool, ammPool)
	}

	// Handle stop loss
	for _, val := range msg.Stoploss {
		position, err := k.GetPosition(ctx, val.Address, val.Id)
		if err != nil {
			continue
		}
		pool, found := k.GetPool(ctx, position.AmmPoolId)
		ammPool, err := k.GetAmmPool(ctx, position.AmmPoolId)
		if !found || err != nil {
			continue
		}
		_, _ = k.ClosePositionIfUnderStopLossPrice(ctx, &position, pool, ammPool)
	}

	return &types.MsgClosePositionsResponse{}, nil
}
