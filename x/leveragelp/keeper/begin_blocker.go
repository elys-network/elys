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
	params := k.GetParams(ctx)

	if epochPosition == 0 && params.FallbackEnabled { // if epoch has passed
		positions := k.GetAllPositions(ctx)
		for _, position := range positions {
			pool, found := k.GetPool(ctx, position.AmmPoolId)
			if !found {
				continue
			}
			ammPool, err := k.GetAmmPool(ctx, pool.AmmPoolId)
			if err != nil {
				ctx.Logger().Error(errors.Wrap(err, fmt.Sprintf("error getting amm pool: %d", pool.AmmPoolId)).Error())
				continue
			}
			isHealthy, _ := k.LiquidatePositionIfUnhealthy(ctx, &position, pool, ammPool)
			if !isHealthy {
				continue
			}
			k.ClosePositionIfUnderStopLossPrice(ctx, &position, pool, ammPool)
		}
	}
}

func (k Keeper) LiquidatePositionIfUnhealthy(ctx sdk.Context, position *types.Position, pool types.Pool, ammPool ammtypes.Pool) (isHealthy, earlyReturn bool) {
	defer func() {
		if r := recover(); r != nil {
			if msg, ok := r.(string); ok {
				ctx.Logger().Error(msg)
			}
		}
	}()
	h, err := k.GetPositionHealth(ctx, *position)
	if err != nil {
		ctx.Logger().Error(errors.Wrap(err, fmt.Sprintf("error updating position health: %s", position.String())).Error())
		return false, true
	}
	position.PositionHealth = h
	k.SetPosition(ctx, position)

	params := k.GetParams(ctx)
	isHealthy = position.PositionHealth.GT(params.SafetyFactor)
	if isHealthy {
		return isHealthy, false
	}

	repayAmount, err := k.ForceCloseLong(ctx, *position, pool, position.LeveragedLpAmount)
	if err != nil {
		ctx.Logger().Error(errors.Wrap(err, "error executing liquidation").Error())
		return isHealthy, true
	}
	ctx.EventManager().EmitEvent(sdk.NewEvent(types.EventClose,
		sdk.NewAttribute("id", strconv.FormatInt(int64(position.Id), 10)),
		sdk.NewAttribute("address", position.Address),
		sdk.NewAttribute("collateral", position.Collateral.String()),
		sdk.NewAttribute("repay_amount", repayAmount.String()),
		sdk.NewAttribute("leverage", position.Leverage.String()),
		sdk.NewAttribute("liabilities", position.Liabilities.String()),
		sdk.NewAttribute("health", position.PositionHealth.String()),
	))
	return isHealthy, false
}

func (k Keeper) ClosePositionIfUnderStopLossPrice(ctx sdk.Context, position *types.Position, pool types.Pool, ammPool ammtypes.Pool) (underStopLossPrice, earlyReturn bool) {
	defer func() {
		if r := recover(); r != nil {
			if msg, ok := r.(string); ok {
				ctx.Logger().Error(msg)
			}
		}
	}()
	h, err := k.GetPositionHealth(ctx, *position)
	if err != nil {
		ctx.Logger().Error(errors.Wrap(err, fmt.Sprintf("error updating position health: %s", position.String())).Error())
		return false, true
	}
	position.PositionHealth = h
	k.SetPosition(ctx, position)

	lpTokenPrice, err := ammPool.LpTokenPrice(ctx, k.oracleKeeper)
	if err != nil {
		return false, true
	}

	underStopLossPrice = !position.StopLossPrice.IsNil() && lpTokenPrice.LTE(position.StopLossPrice)
	if !underStopLossPrice {
		return underStopLossPrice, false
	}

	repayAmount, err := k.ForceCloseLong(ctx, *position, pool, position.LeveragedLpAmount)
	if err != nil {
		ctx.Logger().Error(errors.Wrap(err, "error executing close for stopLossPrice").Error())
		return underStopLossPrice, true
	}
	ctx.EventManager().EmitEvent(sdk.NewEvent(types.EventClose,
		sdk.NewAttribute("id", strconv.FormatInt(int64(position.Id), 10)),
		sdk.NewAttribute("address", position.Address),
		sdk.NewAttribute("collateral", position.Collateral.String()),
		sdk.NewAttribute("repay_amount", repayAmount.String()),
		sdk.NewAttribute("leverage", position.Leverage.String()),
		sdk.NewAttribute("liabilities", position.Liabilities.String()),
		sdk.NewAttribute("health", position.PositionHealth.String()),
	))
	return underStopLossPrice, false
}
