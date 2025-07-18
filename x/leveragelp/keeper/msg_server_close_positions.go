package keeper

import (
	"context"
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v6/x/leveragelp/types"
)

func (k msgServer) ClosePositions(goCtx context.Context, msg *types.MsgClosePositions) (*types.MsgClosePositionsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Handle liquidations
	liqLog := []uint64{}
	for _, val := range msg.Liquidate {
		position, err := k.GetPosition(ctx, val.GetAccountAddress(), val.Id)
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
			liqLog = append(liqLog, position.Id)
			ctx.Logger().Error(fmt.Sprintf("Unhealthy Position: Address:%s Id:%d cannot be liquidated due to err: %s", position.Address, position.Id, err.Error()))
		}
	}

	// Handle stop loss
	closeLog := []uint64{}
	for _, val := range msg.StopLoss {
		position, err := k.GetPosition(ctx, val.GetAccountAddress(), val.Id)
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
			closeLog = append(closeLog, position.Id)
			ctx.Logger().Error(fmt.Sprintf("Stop Loss Position: Address:%s Id:%d cannot be liquidated due to err: %s", position.Address, position.Id, err.Error()))
		}
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(types.EventFailedClosePositions,
		sdk.NewAttribute("liquidations", strings.Trim(strings.Replace(fmt.Sprint(liqLog), " ", ",", -1), "[]")),
		sdk.NewAttribute("stop_loss", strings.Trim(strings.Replace(fmt.Sprint(closeLog), " ", ",", -1), "[]")),
	))

	return &types.MsgClosePositionsResponse{}, nil
}
