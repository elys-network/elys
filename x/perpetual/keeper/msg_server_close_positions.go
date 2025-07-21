package keeper

import (
	"context"
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v6/x/perpetual/types"
)

func (k msgServer) ClosePositions(goCtx context.Context, msg *types.MsgClosePositions) (*types.MsgClosePositionsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Log for errors in liquidations, stop losses, and take profits
	liqLog := []string{}
	stopLossLog := []string{}
	takeProfitLog := []string{}

	// Handle liquidations
	for _, val := range msg.Liquidate {
		owner := sdk.MustAccAddressFromBech32(val.Address)
		position, err := k.GetMTP(ctx, val.PoolId, owner, val.Id)
		if err != nil {
			continue
		}

		// We fetch the amm pool again as there can be changes in amm pool when position is closed
		pool, poolFound := k.GetPool(ctx, position.AmmPoolId)
		if !poolFound {
			continue
		}
		ammPool, poolErr := k.GetAmmPool(ctx, position.AmmPoolId)
		if poolErr != nil {
			continue
		}

		cachedCtx, write := ctx.CacheContext()
		err = k.CheckAndLiquidatePosition(cachedCtx, &position, pool, &ammPool, msg.Creator)
		if err == nil {
			write()
		}
		if err != nil {
			liqLog = append(liqLog, fmt.Sprintf("MTP Unhealthy Close Position: Address:%s Id:%d cannot be liquidated due to err: %s", position.Address, position.Id, err.Error()))
		}
	}

	//Handle StopLoss
	for _, val := range msg.StopLoss {
		owner := sdk.MustAccAddressFromBech32(val.Address)
		position, err := k.GetMTP(ctx, val.PoolId, owner, val.Id)
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
		err = k.CheckAndLiquidatePosition(cachedCtx, &position, pool, &ammPool, msg.Creator)
		if err == nil {
			write()
		}
		if err != nil {
			stopLossLog = append(stopLossLog, fmt.Sprintf("MTP StopLossPrice Close Position: Address:%s Id:%d cannot be liquidated due to err: %s", position.Address, position.Id, err.Error()))
		}
	}

	//Handle take profit
	for _, val := range msg.TakeProfit {
		owner := sdk.MustAccAddressFromBech32(val.Address)
		position, err := k.GetMTP(ctx, val.PoolId, owner, val.Id)
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
		err = k.CheckAndLiquidatePosition(cachedCtx, &position, pool, &ammPool, msg.Creator)
		if err == nil {
			write()
		}
		if err != nil {
			takeProfitLog = append(takeProfitLog, fmt.Sprintf("MTP TakeProfitPrice Close Position: Address:%s Id:%d cannot be liquidated due to err: %s", position.Address, position.Id, err.Error()))
		}
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(types.EventClosePositions,
		sdk.NewAttribute("liquidations", strings.Join(liqLog, "\n")),
		sdk.NewAttribute("stop_losses", strings.Join(stopLossLog, "\n")),
		sdk.NewAttribute("take_profits", strings.Join(takeProfitLog, "\n")),
	))

	return &types.MsgClosePositionsResponse{}, nil
}
