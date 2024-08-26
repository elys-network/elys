package keeper

import (
	"context"
	"fmt"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/leveragelp/types"
)

func (k msgServer) ClosePositions(goCtx context.Context, msg *types.MsgClosePositions) (*types.MsgClosePositionsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	poolMap := make(map[uint64]types.Pool)
	ammPoolMap := make(map[uint64]ammtypes.Pool)
	// Handle liquidations
	liqLog := []string{}
	for _, val := range msg.Liquidate {
		position, err := k.GetPosition(ctx, val.GetAccountAddress(), val.Id)
		if err != nil {
			continue
		}

		pool, found := poolMap[position.AmmPoolId]
		if !found {
			leveragePool, poolFound := k.GetPool(ctx, position.AmmPoolId)
			if !poolFound {
				continue
			}
			poolMap[position.AmmPoolId] = leveragePool

			ammPool, poolErr := k.GetAmmPool(ctx, position.AmmPoolId)
			if poolErr != nil {
				continue
			}
			ammPoolMap[position.AmmPoolId] = ammPool
		}
		pool = poolMap[position.AmmPoolId]
		ammPool := ammPoolMap[position.AmmPoolId]

		_, _, _, err = k.CheckAndLiquidateUnhealthyPosition(ctx, &position, pool, ammPool)
		if err != nil {
			// Add log about error or not liquidated
			liqLog = append(liqLog, fmt.Sprintf("Position: Address:%s Id:%d cannot be liquidated due to err: %s", position.Address, position.Id, err.Error()))
		}
	}

	// Handle stop loss
	closeLog := []string{}
	for _, val := range msg.StopLoss {
		position, err := k.GetPosition(ctx, val.GetAccountAddress(), val.Id)
		if err != nil {
			continue
		}

		pool, found := poolMap[position.AmmPoolId]
		if !found {
			leveragePool, poolFound := k.GetPool(ctx, position.AmmPoolId)
			if !poolFound {
				continue
			}
			poolMap[position.AmmPoolId] = leveragePool

			ammPool, poolErr := k.GetAmmPool(ctx, position.AmmPoolId)
			if poolErr != nil {
				continue
			}
			ammPoolMap[position.AmmPoolId] = ammPool
		}
		pool = poolMap[position.AmmPoolId]
		ammPool := ammPoolMap[position.AmmPoolId]

		_, _, err = k.CheckAndCloseAtStopLoss(ctx, &position, pool, ammPool)
		if err != nil {
			// Add log about error or not closed
			closeLog = append(closeLog, fmt.Sprintf("Position: Address:%s Id:%d cannot be liquidated due to err: %s", position.Address, position.Id, err.Error()))
		}
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(types.EventClosePositions,
		sdk.NewAttribute("liquidations", strings.Join(liqLog, "\n")),
		sdk.NewAttribute("stop_loss", strings.Join(closeLog, "\n")),
	))

	return &types.MsgClosePositionsResponse{}, nil
}
