package keeper

import (
	"context"
	"strconv"
	"strings"

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
		pool, found := k.GetPool(ctx, position.AmmPoolId)
		ammPool, err := k.GetAmmPool(ctx, position.AmmPoolId)
		if !found || err != nil {
			continue
		}
		isHealthy, isEarly, health, err := k.LiquidatePositionIfUnhealthy(ctx, &position, pool, ammPool)
		// position is liquidated
		if err != nil {
			// Add log about error or not liquidated
			liqLog = append(liqLog, "Position: Address:%s Id:%d cannot be liquidated due to err: %s", position.Address, strconv.FormatUint(position.Id, 10), err.Error())
		} else if !(!isHealthy && !isEarly) {
			liqLog = append(liqLog, "Position: Address:%s Id:%s is healthy: %s", position.Address, strconv.FormatUint(position.Id, 10), health.String())
		}
	}

	// Handle stop loss
	closeLog := []string{}
	for _, val := range msg.StopLoss {
		position, err := k.GetPosition(ctx, val.GetAccountAddress(), val.Id)
		if err != nil {
			continue
		}
		pool, found := k.GetPool(ctx, position.AmmPoolId)
		ammPool, err := k.GetAmmPool(ctx, position.AmmPoolId)
		if !found || err != nil {
			continue
		}
		under, early, err := k.ClosePositionIfUnderStopLossPrice(ctx, &position, pool, ammPool)
		if err != nil {
			// Add log about error or not closed
			closeLog = append(closeLog, "Position: Address:%s Id:%s cannot be liquidated due to err: %s", position.Address, strconv.FormatUint(position.Id, 10), err.Error())
		} else if !(under && !early) {
			closeLog = append(closeLog, "Position: Address:%s Id:%s is not under stop loss", position.Address, strconv.FormatUint(position.Id, 10))
		}
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(types.EventClosePositions,
		sdk.NewAttribute("liquidations", strings.Join(liqLog, "\n")),
		sdk.NewAttribute("stop_loss", strings.Join(closeLog, "\n")),
	))

	return &types.MsgClosePositionsResponse{}, nil
}
