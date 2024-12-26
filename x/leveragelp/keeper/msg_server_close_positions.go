package keeper

import (
	"context"
	"fmt"
	"strings"

	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/leveragelp/types"
)

func (k msgServer) ClosePositions(goCtx context.Context, msg *types.MsgClosePositions) (*types.MsgClosePositionsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Handle liquidations
	liqLog := []string{}
	for _, val := range msg.Liquidate {
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
		_, _, _, err = k.CheckAndLiquidateUnhealthyPosition(cachedCtx, &position, pool, ammPool)
		if err == nil {
			write()
		}
		if err != nil {
			// Add log about error or not liquidated
			liqLog = append(liqLog, fmt.Sprintf("Position: Address:%s Id:%d cannot be liquidated due to err: %s", position.Address, position.Id, err.Error()))
		}

		if k.hooks != nil {
			// ammPool will have updated values for opening position
			found := false
			ammPool, found = k.amm.GetPool(ctx, position.AmmPoolId)
			if !found {
				return nil, errorsmod.Wrap(types.ErrPoolDoesNotExist, fmt.Sprintf("poolId: %d", position.AmmPoolId))
			}
			err = k.hooks.AfterLeverageLpPositionClose(ctx, position.GetOwnerAddress(), ammPool)
			if err != nil {
				return nil, err
			}
		}
	}

	// Handle stop loss
	closeLog := []string{}
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
			// Add log about error or not closed
			closeLog = append(closeLog, fmt.Sprintf("Position: Address:%s Id:%d cannot be liquidated due to err: %s", position.Address, position.Id, err.Error()))
		}

		if k.hooks != nil {
			// ammPool will have updated values for opening position
			found := false
			ammPool, found = k.amm.GetPool(ctx, position.AmmPoolId)
			if !found {
				return nil, errorsmod.Wrap(types.ErrPoolDoesNotExist, fmt.Sprintf("poolId: %d", position.AmmPoolId))
			}
			err = k.hooks.AfterLeverageLpPositionClose(ctx, position.GetOwnerAddress(), ammPool)
			if err != nil {
				return nil, err
			}
		}
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(types.EventClosePositions,
		sdk.NewAttribute("liquidations", strings.Join(liqLog, "\n")),
		sdk.NewAttribute("stop_loss", strings.Join(closeLog, "\n")),
	))

	return &types.MsgClosePositionsResponse{}, nil
}
