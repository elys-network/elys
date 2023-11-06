package keeper

import (
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/types/errors"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/leveragelp/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) BeginBlocker(ctx sdk.Context) {
	// check if epoch has passed then execute
	epochLength := k.GetEpochLength(ctx)
	epochPosition := k.GetEpochPosition(ctx, epochLength)

	if epochPosition == 0 { // if epoch has passed
		pools := k.GetAllPools(ctx)
		for _, pool := range pools {
			ammPool, err := k.GetAmmPool(ctx, pool.AmmPoolId)
			if err != nil {
				ctx.Logger().Error(errors.Wrap(err, fmt.Sprintf("error getting amm pool: %d", pool.AmmPoolId)).Error())
				continue
			}
			if k.IsPoolEnabled(ctx, pool.AmmPoolId) {
				positions, _, _ := k.GetPositionsForPool(ctx, pool.AmmPoolId, nil)
				for _, position := range positions {
					k.LiquidatePositionIfUnhealthy(ctx, position, pool, ammPool)
				}
			}
			k.SetPool(ctx, pool)
		}
	}
}

func (k Keeper) LiquidatePositionIfUnhealthy(ctx sdk.Context, position *types.Position, pool types.Pool, ammPool ammtypes.Pool) {
	defer func() {
		if r := recover(); r != nil {
			if msg, ok := r.(string); ok {
				ctx.Logger().Error(msg)
			}
		}
	}()
	h, err := k.GetPositionHealth(ctx, *position, ammPool)
	if err != nil {
		ctx.Logger().Error(errors.Wrap(err, fmt.Sprintf("error updating position health: %s", position.String())).Error())
		return
	}
	position.PositionHealth = h
	k.SetPosition(ctx, position)

	lpTokenPrice, err := k.GetLpTokenPrice(ctx, &ammPool)
	if err != nil {
		return
	}

	params := k.GetParams(ctx)
	if position.PositionHealth.GT(params.SafetyFactor) && lpTokenPrice.GT(position.StopLossPrice) {
		return
	}

	repayAmount, err := k.ForceCloseLong(ctx, *position, pool)
	if err == nil {
		ctx.EventManager().EmitEvent(sdk.NewEvent(types.EventClose,
			sdk.NewAttribute("id", strconv.FormatInt(int64(position.Id), 10)),
			sdk.NewAttribute("address", position.Address),
			sdk.NewAttribute("collateral", position.Collateral.String()),
			sdk.NewAttribute("repay_amount", repayAmount.String()),
			sdk.NewAttribute("leverage", position.Leverage.String()),
			sdk.NewAttribute("liabilities", position.Liabilities.String()),
			sdk.NewAttribute("health", position.PositionHealth.String()),
		))
	} else {
		ctx.Logger().Error(errors.Wrap(err, "error executing force close").Error())
	}
}
