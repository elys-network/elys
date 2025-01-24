package keeper

import (
	"errors"
	"fmt"
	"strconv"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/types/query"
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
		pageReq := &query.PageRequest{
			Limit:      uint64(params.NumberPerBlock),
			CountTotal: true,
		}
		offset, _ := k.GetOffset(ctx)
		pageReq.Offset = offset
		positions, _, err := k.GetPositions(ctx, pageReq)
		if err != nil {
			ctx.Logger().Error(errorsmod.Wrap(err, "error fetching paginated positions").Error())
			return
		}
		if offset+uint64(params.NumberPerBlock) >= k.GetOpenPositionCount(ctx) {
			k.DeleteOffset(ctx)
		} else {
			k.SetOffset(ctx, offset+uint64(params.NumberPerBlock))
		}

		for _, position := range positions {
			pool, found := k.GetPool(ctx, position.AmmPoolId)
			if !found {
				ctx.Logger().Error(fmt.Sprintf("pool not found for id: %d", position.AmmPoolId))
				continue
			}
			ammPool, poolErr := k.GetAmmPool(ctx, pool.AmmPoolId)
			if poolErr != nil {
				ctx.Logger().Error(fmt.Sprintf("error getting for amm pool %d: %s", position.AmmPoolId, poolErr.Error()))
				continue
			}

			isHealthy, closeAttempted, _, err := k.CheckAndLiquidateUnhealthyPosition(ctx, position, pool, ammPool)
			if err == nil {
				continue
			}
			if isHealthy && !closeAttempted {
				_, _, _ = k.CheckAndCloseAtStopLoss(ctx, position, pool, ammPool)
			}
		}
	}
}

func (k Keeper) CheckAndLiquidateUnhealthyPosition(ctx sdk.Context, position *types.Position, pool types.Pool, ammPool ammtypes.Pool) (isHealthy, closeAttempted bool, health math.LegacyDec, err error) {
	defer func() {
		if r := recover(); r != nil {
			if msg, ok := r.(string); ok {
				ctx.Logger().Error(msg)
			}
		}
	}()
	h, err := k.GetPositionHealth(ctx, *position)
	if err != nil {
		ctx.Logger().Error(errorsmod.Wrap(err, fmt.Sprintf("error updating position health: %s", position.String())).Error())
		return false, false, math.LegacyZeroDec(), err
	}
	position.PositionHealth = h
	k.SetPosition(ctx, position)

	params := k.GetParams(ctx)
	isHealthy = position.PositionHealth.GT(params.SafetyFactor)

	debt := k.stableKeeper.UpdateInterestAndGetDebt(ctx, position.GetPositionAddress(), position.AmmPoolId, position.Collateral.Denom)
	liab := debt.GetTotalLiablities()
	if isHealthy || liab.IsZero() {
		return true, false, h, errors.New("position is healthy to close")
	}

	repayAmount, err := k.ForceCloseLong(ctx, *position, pool, position.LeveragedLpAmount, true)
	if err != nil {
		ctx.Logger().Debug(errorsmod.Wrap(err, "error executing liquidation").Error())
		return isHealthy, true, h, err
	}
	ctx.EventManager().EmitEvent(sdk.NewEvent(types.EventCloseUnhealthyPosition,
		sdk.NewAttribute("id", strconv.FormatInt(int64(position.Id), 10)),
		sdk.NewAttribute("address", position.Address),
		sdk.NewAttribute("collateral", position.Collateral.String()),
		sdk.NewAttribute("repay_amount", repayAmount.String()),
		sdk.NewAttribute("liabilities", position.Liabilities.String()),
		sdk.NewAttribute("health", position.PositionHealth.String()),
	))
	return isHealthy, true, h, nil
}

func (k Keeper) CheckAndCloseAtStopLoss(ctx sdk.Context, position *types.Position, pool types.Pool, ammPool ammtypes.Pool) (underStopLossPrice, closeAttempted bool, err error) {
	defer func() {
		if r := recover(); r != nil {
			if msg, ok := r.(string); ok {
				ctx.Logger().Error(msg)
			}
		}
	}()
	h, err := k.GetPositionHealth(ctx, *position)
	if err != nil {
		ctx.Logger().Error(errorsmod.Wrap(err, fmt.Sprintf("error updating position health: %s", position.String())).Error())
		return false, false, err
	}
	position.PositionHealth = h
	k.SetPosition(ctx, position)

	lpTokenPrice, err := ammPool.LpTokenPrice(ctx, k.oracleKeeper, k.accountedPoolKeeper)
	if err != nil {
		return false, false, err
	}

	underStopLossPrice = !position.StopLossPrice.IsNil() && lpTokenPrice.LTE(position.StopLossPrice)
	if !underStopLossPrice {
		return underStopLossPrice, false, errors.New("position stop loss price is not <= lp token price")
	}

	repayAmount, err := k.ForceCloseLong(ctx, *position, pool, position.LeveragedLpAmount, false)
	if err != nil {
		ctx.Logger().Error(errorsmod.Wrap(err, "error executing close for stopLossPrice").Error())
		return underStopLossPrice, true, err
	}
	ctx.EventManager().EmitEvent(sdk.NewEvent(types.EventClosePositionStopLoss,
		sdk.NewAttribute("id", strconv.FormatInt(int64(position.Id), 10)),
		sdk.NewAttribute("address", position.Address),
		sdk.NewAttribute("collateral", position.Collateral.String()),
		sdk.NewAttribute("repay_amount", repayAmount.String()),
		sdk.NewAttribute("liabilities", position.Liabilities.String()),
		sdk.NewAttribute("health", position.PositionHealth.String()),
	))
	return underStopLossPrice, true, nil
}
