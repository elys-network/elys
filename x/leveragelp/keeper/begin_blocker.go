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

			isHealthy, closeAttempted, _, err := k.CheckAndLiquidateUnhealthyPosition(ctx, position, pool)
			if err == nil {
				continue
			}
			if isHealthy && !closeAttempted {
				ammPool, poolErr := k.GetAmmPool(ctx, pool.AmmPoolId)
				if poolErr != nil {
					ctx.Logger().Error(fmt.Sprintf("error getting for amm pool %d: %s", position.AmmPoolId, poolErr.Error()))
					continue
				}
				_, _, _ = k.CheckAndCloseAtStopLoss(ctx, position, pool, ammPool)
			}
		}
	}
}

func (k Keeper) CheckAndLiquidateUnhealthyPosition(ctx sdk.Context, position *types.Position, pool types.Pool) (isHealthy, closeAttempted bool, health math.LegacyDec, err error) {
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

	if isHealthy {
		return true, false, h, errors.New("position is healthy to close")
	}

	finalClosingRatio, totalLpAmountToClose, coinsForAmm, repayAmount, userReturnTokens, exitFeeOnClosingPosition, stopLossReached, _, exitSlippageFee, swapFee, err := k.CheckHealthStopLossThenRepayAndClose(ctx, position, &pool, math.LegacyOneDec(), true)
	if err != nil {
		ctx.Logger().Error(errorsmod.Wrap(err, "error executing liquidation for unhealthy").Error())
		return isHealthy, true, h, err
	}
	ctx.EventManager().EmitEvent(sdk.NewEvent(types.EventAutomatedClosePosition,
		sdk.NewAttribute("id", strconv.FormatUint(position.Id, 10)),
		sdk.NewAttribute("address", position.Address),
		sdk.NewAttribute("closing_ratio", finalClosingRatio.String()),
		sdk.NewAttribute("lp_amount_closed", totalLpAmountToClose.String()),
		sdk.NewAttribute("coins_to_amm", coinsForAmm.String()),
		sdk.NewAttribute("repay_amount", repayAmount.String()),
		sdk.NewAttribute("user_return_tokens", userReturnTokens.String()),
		sdk.NewAttribute("exit_fee", exitFeeOnClosingPosition.String()),
		sdk.NewAttribute("reason", "unhealthy"),
		sdk.NewAttribute("stop_loss_reached", strconv.FormatBool(stopLossReached)),
		sdk.NewAttribute("exit_slippage_fee", exitSlippageFee.String()),
		sdk.NewAttribute("exit_swap_fee", swapFee.String()),
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

	lpTokenPrice, err := ammPool.LpTokenPriceForShare(ctx, k.oracleKeeper, k.accountedPoolKeeper)
	if err != nil {
		return false, false, err
	}

	underStopLossPrice = !position.StopLossPrice.IsNil() && lpTokenPrice.LTE(position.StopLossPrice)
	if !underStopLossPrice {
		return underStopLossPrice, false, errors.New("position stop loss price is not <= lp token price")
	}

	finalClosingRatio, totalLpAmountToClose, coinsForAmm, repayAmount, userReturnTokens, exitFeeOnClosingPosition, stopLossReached, _, exitSlippageFee, swapFee, err := k.CheckHealthStopLossThenRepayAndClose(ctx, position, &pool, math.LegacyOneDec(), false)
	if err != nil {
		ctx.Logger().Error(errorsmod.Wrap(err, "error executing close for stopLossPrice").Error())
		return underStopLossPrice, true, err
	}
	ctx.EventManager().EmitEvent(sdk.NewEvent(types.EventAutomatedClosePosition,
		sdk.NewAttribute("id", strconv.FormatUint(position.Id, 10)),
		sdk.NewAttribute("address", position.Address),
		sdk.NewAttribute("closing_ratio", finalClosingRatio.String()),
		sdk.NewAttribute("lp_amount_closed", totalLpAmountToClose.String()),
		sdk.NewAttribute("coins_to_amm", coinsForAmm.String()),
		sdk.NewAttribute("repay_amount", repayAmount.String()),
		sdk.NewAttribute("user_return_tokens", userReturnTokens.String()),
		sdk.NewAttribute("exit_fee", exitFeeOnClosingPosition.String()),
		sdk.NewAttribute("reason", "stop_loss"),
		sdk.NewAttribute("stop_loss_reached", strconv.FormatBool(stopLossReached)),
		sdk.NewAttribute("exit_slippage_fee", exitSlippageFee.String()),
		sdk.NewAttribute("exit_swap_fee", swapFee.String()),
	))
	return underStopLossPrice, true, nil
}
