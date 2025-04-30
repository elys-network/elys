package keeper

import (
	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/perpetual/types"
	"github.com/osmosis-labs/osmosis/osmomath"
)

// SetParams set the params
func (k Keeper) SetParams(ctx sdk.Context, params *types.Params) error {
	if err := params.Validate(); err != nil {
		return err
	}

	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	bz, err := k.cdc.Marshal(params)
	if err != nil {
		return err
	}
	store.Set(types.KeyPrefix(types.ParamsKey), bz)

	return nil
}

// GetParams get all parameters as types.Params
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	bz := store.Get(types.KeyPrefix(types.ParamsKey))
	if bz == nil {
		return params
	}

	k.cdc.MustUnmarshal(bz, &params)
	return params
}

func (k Keeper) GetLegacyParams(ctx sdk.Context) (params types.LegacyParams) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	bz := store.Get(types.KeyPrefix(types.ParamsKey))
	if bz == nil {
		return params
	}

	k.cdc.MustUnmarshal(bz, &params)
	return params
}

func (k Keeper) GetMaxLeverageParam(ctx sdk.Context) sdkmath.LegacyDec {
	return k.GetParams(ctx).LeverageMax
}

func (k Keeper) GetBigDecMaxLeverageParam(ctx sdk.Context) osmomath.BigDec {
	return osmomath.BigDecFromDec(k.GetParams(ctx).LeverageMax)
}

func (k Keeper) GetBorrowInterestRateMax(ctx sdk.Context) sdkmath.LegacyDec {
	return k.GetParams(ctx).BorrowInterestRateMax
}

func (k Keeper) GetBigDecBorrowInterestRateMax(ctx sdk.Context) osmomath.BigDec {
	return osmomath.BigDecFromDec(k.GetParams(ctx).BorrowInterestRateMax)
}

func (k Keeper) GetBorrowInterestRateMin(ctx sdk.Context) sdkmath.LegacyDec {
	return k.GetParams(ctx).BorrowInterestRateMin
}

func (k Keeper) GetBigDecBorrowInterestRateMin(ctx sdk.Context) osmomath.BigDec {
	return osmomath.BigDecFromDec(k.GetParams(ctx).BorrowInterestRateMin)
}

func (k Keeper) GetBorrowInterestRateIncrease(ctx sdk.Context) sdkmath.LegacyDec {
	return k.GetParams(ctx).BorrowInterestRateIncrease
}

func (k Keeper) GetBigDecBorrowInterestRateIncrease(ctx sdk.Context) osmomath.BigDec {
	return osmomath.BigDecFromDec(k.GetParams(ctx).BorrowInterestRateIncrease)
}

func (k Keeper) GetBorrowInterestRateDecrease(ctx sdk.Context) sdkmath.LegacyDec {
	return k.GetParams(ctx).BorrowInterestRateDecrease
}
func (k Keeper) GetBigDecBorrowInterestRateDecrease(ctx sdk.Context) osmomath.BigDec {
	return osmomath.BigDecFromDec(k.GetParams(ctx).BorrowInterestRateDecrease)
}

func (k Keeper) GetHealthGainFactor(ctx sdk.Context) sdkmath.LegacyDec {
	return k.GetParams(ctx).HealthGainFactor
}
func (k Keeper) GetBigDecHealthGainFactor(ctx sdk.Context) osmomath.BigDec {
	return osmomath.BigDecFromDec(k.GetParams(ctx).HealthGainFactor)
}

func (k Keeper) GetPoolOpenThreshold(ctx sdk.Context) sdkmath.LegacyDec {
	return k.GetParams(ctx).PoolOpenThreshold
}

func (k Keeper) GetBigDecPoolOpenThreshold(ctx sdk.Context) osmomath.BigDec {
	return osmomath.BigDecFromDec(k.GetParams(ctx).PoolOpenThreshold)
}

func (k Keeper) GetBorrowInterestPaymentFundPercentage(ctx sdk.Context) sdkmath.LegacyDec {
	return k.GetParams(ctx).BorrowInterestPaymentFundPercentage
}

func (k Keeper) GetBorrowInterestPaymentFundAddress(ctx sdk.Context) sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(k.GetParams(ctx).BorrowInterestPaymentFundAddress)
	if err != nil {
		panic(err)
	}

	return addr
}

func (k Keeper) GetMaxOpenPositions(ctx sdk.Context) int64 {
	return k.GetParams(ctx).MaxOpenPositions
}

func (k Keeper) GetBorrowInterestPaymentEnabled(ctx sdk.Context) bool {
	return k.GetParams(ctx).BorrowInterestPaymentEnabled
}

func (k Keeper) GetSafetyFactor(ctx sdk.Context) sdkmath.LegacyDec {
	return k.GetParams(ctx).SafetyFactor
}

func (k Keeper) IsWhitelistingEnabled(ctx sdk.Context) bool {
	return k.GetParams(ctx).WhitelistingEnabled
}

func (k Keeper) GetPerpetualSwapFee(ctx sdk.Context) sdkmath.LegacyDec {
	return k.GetParams(ctx).PerpetualSwapFee
}

func (k Keeper) GetMaxLimitOrder(ctx sdk.Context) int64 {
	return k.GetParams(ctx).MaxLimitOrder
}
