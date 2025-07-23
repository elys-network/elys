package keeper

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	"fmt"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/elys-network/elys/v6/x/perpetual/types"
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
	iterator := storetypes.KVStorePrefixIterator(store, types.ADLCounterPrefix)
	defer iterator.Close()

	var list []types.ADLCounter
	for ; iterator.Valid(); iterator.Next() {
		var adlCounter types.ADLCounter
		k.cdc.MustUnmarshal(iterator.Value(), &adlCounter)
		list = append(list, adlCounter)
	}

	return list
}

func (k Keeper) ClosePositionsOnADL(ctx sdk.Context, perpetualPool types.Pool) error {
	// closing ratio = (current ratio - max leverage) / current ratio
	// we use max leverage instead of adl trigger ratio because then this whole thing will be asymptotic and
	// process will never end. adl trigger stops that by being higher
	currentMaxLiabilitiesRatio := math.LegacyMaxDec(perpetualPool.BaseAssetLiabilitiesRatio, perpetualPool.QuoteAssetLiabilitiesRatio)
	params := k.GetParams(ctx)
	if currentMaxLiabilitiesRatio.LTE(params.PoolMaxLiabilitiesThreshold) {
		return nil
	}
	closingRatio := currentMaxLiabilitiesRatio.Sub(params.PoolMaxLiabilitiesThreshold).Quo(currentMaxLiabilitiesRatio)
	if closingRatio.IsZero() || closingRatio.IsNegative() {
		err := fmt.Errorf("closing ratio is <= 0 for pool while triggering adl for pool id %d in perpetual", perpetualPool.AmmPoolId)
		ctx.Logger().Error(err.Error())
		return err
	}

	if closingRatio.GT(math.LegacyOneDec()) {
		closingRatio = math.LegacyOneDec()
	}

	pageReq := &query.PageRequest{
		Limit:      params.NumberPerBlock,
		CountTotal: true,
	}
	adlCounter := k.GetADLCounter(ctx, perpetualPool.AmmPoolId)
	if len(adlCounter.NextKey) != 0 {
		pageReq.Key = adlCounter.NextKey
	} else {
		pageReq.Offset = 0
	}
	totalOpen := k.GetPerpetualCounter(ctx, perpetualPool.AmmPoolId).TotalOpen
	if adlCounter.Counter+params.NumberPerBlock >= totalOpen {
		adlCounter.Counter = 0
	} else {
		adlCounter.Counter += params.NumberPerBlock
	}

	mtps, pageResponse, err := k.GetPositionsForPool(ctx, perpetualPool.AmmPoolId, pageReq)
	if err != nil {
		ctx.Logger().Error(errorsmod.Wrap(err, "error fetching paginated positions").Error())
		return err
	}

	if adlCounter.Counter == 0 {
		adlCounter.NextKey = nil
	} else {
		adlCounter.NextKey = pageResponse.NextKey
	}
	k.SetADLCounter(ctx, adlCounter)

	for _, mtp := range mtps {
		msg := types.NewMsgClose(mtp.Address, mtp.Id, math.ZeroInt(), perpetualPool.AmmPoolId, closingRatio)
		closedMtp, repayAmount, finalClosingRatio, returnAmt, fundingFeeAmt, fundingAmtDistributed, interestAmt, insuranceAmt, allInterestsPaid, forceClosed, totalPerpetualFeesCoins, closingPrice, initialCollateral, initialCustody, initialLiabilities, err := k.ClosePosition(ctx, msg)
		if err != nil {
			return err
		}

		err = k.EmitClose(ctx, "adl", closedMtp, repayAmount, finalClosingRatio, returnAmt, fundingFeeAmt, fundingAmtDistributed, interestAmt, insuranceAmt, allInterestsPaid, forceClosed, totalPerpetualFeesCoins, closingPrice, initialCollateral, initialCustody, initialLiabilities)
		if err != nil {
			return err
		}
	}
	return nil
}
