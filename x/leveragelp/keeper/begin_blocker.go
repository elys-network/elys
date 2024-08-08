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
				// Liquidate positions liquidation health threshold
				// Design
				// - `Health = PositionValue / liability`, PositionValue is based on LpToken price change
				// - Debt growth speed is relying on liability.
				// - Things are sorted by `LeveragedLpAmount / liability` per pool to liquidate efficiently
				k.IteratePoolPosIdsLiquidationSorted(ctx, pool.AmmPoolId, func(addressId types.AddressId) bool {
					position, err := k.GetPosition(ctx, addressId.GetPositionCreatorAddress(), addressId.Id)
					if err != nil {
						return false
					}
					isHealthy, earlyReturn := k.LiquidatePositionIfUnhealthy(ctx, &position, pool, ammPool)
					if !earlyReturn && isHealthy {
						return true
					}
					return false
				})

				// Close stopLossPrice reached positions
				k.IteratePoolPosIdsStopLossSorted(ctx, pool.AmmPoolId, func(addressId types.AddressId) bool {
					position, err := k.GetPosition(ctx, addressId.GetPositionCreatorAddress(), addressId.Id)
					if err != nil {
						return false
					}
					underStopLossPrice, earlyReturn := k.ClosePositionIfUnderStopLossPrice(ctx, &position, pool, ammPool)
					if !earlyReturn && underStopLossPrice {
						return true
					}
					return false
				})
			}
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
	debt := k.stableKeeper.GetDebt(ctx, position.GetPositionAddress())
	position.PositionHealth = h
	k.SetPosition(ctx, position, debt.Borrowed.Add(debt.InterestStacked).Sub(debt.InterestPaid))

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
	debt := k.stableKeeper.GetDebt(ctx, position.GetPositionAddress())
	position.PositionHealth = h
	k.SetPosition(ctx, position, debt.Borrowed.Add(debt.InterestStacked).Sub(debt.InterestPaid))

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
