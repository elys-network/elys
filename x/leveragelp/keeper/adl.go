package keeper

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/elys-network/elys/v6/x/leveragelp/types"
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

func (k Keeper) AutoClosePositions(ctx sdk.Context, leveragePool types.Pool) error {

	closingRatio := math.LegacyOneDec()

	pageReq := &query.PageRequest{
		Limit:      100,
		CountTotal: false,
	}

	positions, _, err := k.GetPositionsForPool(ctx, leveragePool.AmmPoolId, pageReq)
	if err != nil {
		ctx.Logger().Error(errorsmod.Wrap(err, "error fetching paginated positions").Error())
		return err
	}

	for _, position := range positions {
		cacheCtx, writeCache := ctx.CacheContext()
		finalClosingRatio, totalLpAmountToClose, coinsForAmm, repayAmount, userReturnTokens, exitFeeOnClosingPosition, stopLossReached, _, exitSlippageFee, swapFee, takerFee, slippageValue, swapFeeValue, takerFeeValue, weightBreakingFeeValue, err := k.CheckHealthStopLossThenRepayAndClose(cacheCtx, &position, &leveragePool, closingRatio, false)
		if err != nil {
			ctx.Logger().Error(errorsmod.Wrap(err, "error executing auto close").Error())
			continue
		} else {
			writeCache()
			k.EmitCloseEvent(ctx, "auto_close", position, finalClosingRatio, totalLpAmountToClose, coinsForAmm, repayAmount, userReturnTokens, exitFeeOnClosingPosition, stopLossReached, exitSlippageFee, swapFee, takerFee, slippageValue, swapFeeValue, takerFeeValue, weightBreakingFeeValue)
		}
	}
	return nil
}
