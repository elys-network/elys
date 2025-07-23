package keeper

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	"fmt"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/elys-network/elys/v6/x/leveragelp/types"
	"strconv"
)

func (k Keeper) SetADLCounter(ctx sdk.Context, adlCounter types.ADLCounter) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.GetADLCounterKey(adlCounter.PoolId)
	bz := k.cdc.MustMarshal(&adlCounter)
	store.Set(key, bz)
}

func (k Keeper) GetADLCounter(ctx sdk.Context, poolId uint64) (adlCounter types.ADLCounter) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.GetADLCounterKey(poolId)
	bz := store.Get(key)
	if bz == nil {
		return types.ADLCounter{
			PoolId:  poolId,
			Counter: 0,
		}
	}
	k.cdc.MustUnmarshal(bz, &adlCounter)
	return adlCounter
}

func (k Keeper) GetAllADLCounter(ctx sdk.Context) []types.ADLCounter {

	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	iterator := storetypes.KVStorePrefixIterator(store, types.ADLCounterKeyPrefix)
	defer iterator.Close()

	var list []types.ADLCounter
	for ; iterator.Valid(); iterator.Next() {
		var adlCounter types.ADLCounter
		k.cdc.MustUnmarshal(iterator.Value(), &adlCounter)
		list = append(list, adlCounter)
	}

	return list
}

func (k Keeper) ClosePositionsOnADL(ctx sdk.Context, leveragePool types.Pool) error {
	// closing ratio = (current ratio - max leverage) / current ratio
	// we use max leverage instead of adl trigger ratio because then this whole thing will be asymptotic and
	// process will never end. adl trigger stops that by being higher
	ammPool, err := k.GetAmmPool(ctx, leveragePool.AmmPoolId)
	if err != nil {
		return err
	}

	// Check for division by zero
	if ammPool.TotalShares.Amount.IsZero() {
		return fmt.Errorf("amm pool %d has zero total shares", leveragePool.AmmPoolId)
	}
	currentLeverageRatio := leveragePool.LeveragedLpAmount.ToLegacyDec().Quo(ammPool.TotalShares.Amount.ToLegacyDec())

	if currentLeverageRatio.LTE(leveragePool.AdlTriggerRatio) {
		return nil
	}
	closingRatio := currentLeverageRatio.Sub(leveragePool.MaxLeveragelpRatio).Quo(currentLeverageRatio)
	if closingRatio.IsZero() || closingRatio.IsNegative() {
		err := fmt.Errorf("closing ratio is <= 0 for pool while triggering adl for pool id %d", leveragePool.AmmPoolId)
		ctx.Logger().Error(err.Error())
		return err
	}

	if closingRatio.GT(math.LegacyOneDec()) {
		closingRatio = math.LegacyOneDec()
	}

	params := k.GetParams(ctx)

	pageReq := &query.PageRequest{
		Limit:      uint64(params.NumberPerBlock),
		CountTotal: true,
	}
	adlCounter := k.GetADLCounter(ctx, leveragePool.AmmPoolId)
	if len(adlCounter.NextKey) != 0 {
		pageReq.Key = adlCounter.NextKey
	} else {
		pageReq.Offset = 0
	}
	totalOpen := k.GetPositionCounter(ctx, leveragePool.AmmPoolId).TotalOpen
	if adlCounter.Counter+uint64(params.NumberPerBlock) >= totalOpen {
		adlCounter.Counter = 0
	} else {
		adlCounter.Counter += uint64(params.NumberPerBlock)
	}

	positions, pageResponse, err := k.GetPositionsForPool(ctx, leveragePool.AmmPoolId, pageReq)
	if err != nil {
		ctx.Logger().Error(errorsmod.Wrap(err, "error fetching paginated positions").Error())
		return err
	}
	adlCounter.NextKey = pageResponse.NextKey
	k.SetADLCounter(ctx, adlCounter)

	for _, position := range positions {
		finalClosingRatio, totalLpAmountToClose, coinsForAmm, repayAmount, userReturnTokens, exitFeeOnClosingPosition, stopLossReached, _, exitSlippageFee, swapFee, takerFee, err := k.CheckHealthStopLossThenRepayAndClose(ctx, &position, &leveragePool, closingRatio, false)
		if err != nil {
			ctx.Logger().Error(errorsmod.Wrap(err, "error executing close for stopLossPrice").Error())
			return err
		}
		ctx.EventManager().EmitEvent(sdk.NewEvent(types.EventADLTriggeredClosePosition,
			sdk.NewAttribute("id", strconv.FormatUint(position.Id, 10)),
			sdk.NewAttribute("address", position.Address),
			sdk.NewAttribute("poolId", strconv.FormatUint(position.AmmPoolId, 10)),
			sdk.NewAttribute("closing_ratio", finalClosingRatio.String()),
			sdk.NewAttribute("lp_amount_closed", totalLpAmountToClose.String()),
			sdk.NewAttribute("coins_to_amm", coinsForAmm.String()),
			sdk.NewAttribute("repay_amount", repayAmount.String()),
			sdk.NewAttribute("user_return_tokens", userReturnTokens.String()),
			sdk.NewAttribute("exit_fee", exitFeeOnClosingPosition.String()),
			sdk.NewAttribute("health", position.PositionHealth.String()),
			sdk.NewAttribute("stop_loss_reached", strconv.FormatBool(stopLossReached)),
			sdk.NewAttribute("exit_slippage_fee", exitSlippageFee.String()),
			sdk.NewAttribute("exit_swap_fee", swapFee.String()),
			sdk.NewAttribute("exit_taker_fee", takerFee.String()),
		))
	}
	return nil
}
