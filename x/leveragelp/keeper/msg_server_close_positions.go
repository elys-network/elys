package keeper

import (
	"context"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v7/x/leveragelp/types"
	"strings"
)

func (k msgServer) ClosePositions(goCtx context.Context, msg *types.MsgClosePositions) (*types.MsgClosePositionsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Handle liquidations
	liqLog := []string{}
	for _, val := range msg.Liquidate {
		position, err := k.GetPosition(ctx, val.PoolId, val.GetAccountAddress(), val.Id)
		if err != nil {
			continue
		}

		pool, poolFound := k.GetPool(ctx, position.AmmPoolId)
		if !poolFound {
			continue
		}

		cachedCtx, write := ctx.CacheContext()
		_, _, _, err = k.CheckAndLiquidateUnhealthyPosition(cachedCtx, &position, pool)
		if err == nil {
			write()
		}
		if err != nil {
			liqLog = append(liqLog, getFailedLiquidationLog(position.AmmPoolId, position.Id))
			ctx.Logger().Error(fmt.Sprintf("Unhealthy Position: PoolId: %d Address:%s Id:%d cannot be liquidated due to err: %s", position.AmmPoolId, position.Address, position.Id, err.Error()))
		}
	}

	// Handle stop loss
	closeLog := []string{}
	for _, val := range msg.StopLoss {
		position, err := k.GetPosition(ctx, val.PoolId, val.GetAccountAddress(), val.Id)
		if err != nil {
			continue
		}

		pool, poolFound := k.GetPool(ctx, position.AmmPoolId)
		if !poolFound {
			continue
		}
		ammPool, poolErr := k.GetAmmPool(ctx, position.AmmPoolId)
		if poolErr != nil {
			continue
		}
		cachedCtx, write := ctx.CacheContext()
		_, _, err = k.CheckAndCloseAtStopLoss(cachedCtx, &position, pool, ammPool)
		if err == nil {
			write()
		}
		if err != nil {
			closeLog = append(closeLog, getFailedLiquidationLog(position.AmmPoolId, position.Id))
			ctx.Logger().Error(fmt.Sprintf("Stop Loss Position: PoolId: %d Address:%s Id:%d cannot be liquidated due to err: %s", position.AmmPoolId, position.Address, position.Id, err.Error()))
		}
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(types.EventFailedClosePositions,
		sdk.NewAttribute("sender", msg.Creator),
		sdk.NewAttribute("liquidations", strings.Join(liqLog, ",")),
		sdk.NewAttribute("stop_loss", strings.Join(closeLog, ",")),
	))

	return &types.MsgClosePositionsResponse{}, nil
}

func getFailedLiquidationLog(poolId, positionId uint64) string {
	return fmt.Sprintf("PoolId_%d/Id_%d", poolId, positionId)
}
