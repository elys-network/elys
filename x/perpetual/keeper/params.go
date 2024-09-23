package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

// SetParams set the params
func (k Keeper) SetParams(ctx sdk.Context, params *types.Params) error {
	if err := params.Validate(); err != nil {
		return err
	}

	store := ctx.KVStore(k.storeKey)
	bz, err := k.cdc.Marshal(params)
	if err != nil {
		return err
	}
	store.Set(types.KeyPrefix(types.ParamsKey), bz)

	return nil
}

// GetParams get all parameters as types.Params
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.KeyPrefix(types.ParamsKey))
	if bz == nil {
		return params
	}

	k.cdc.MustUnmarshal(bz, &params)
	return params
}

func (k Keeper) GetLeagcyParams(ctx sdk.Context) (params types.LegacyParams) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.KeyPrefix(types.ParamsKey))
	if bz == nil {
		return params
	}

	k.cdc.MustUnmarshal(bz, &params)
	return params
}

func (k Keeper) GetMaxLeverageParam(ctx sdk.Context) sdk.Dec {
	return k.GetParams(ctx).LeverageMax
}

func (k Keeper) GetBorrowInterestRateMax(ctx sdk.Context) sdk.Dec {
	return k.GetParams(ctx).BorrowInterestRateMax
}

func (k Keeper) GetBorrowInterestRateMin(ctx sdk.Context) sdk.Dec {
	return k.GetParams(ctx).BorrowInterestRateMin
}

func (k Keeper) GetBorrowInterestRateIncrease(ctx sdk.Context) sdk.Dec {
	return k.GetParams(ctx).BorrowInterestRateIncrease
}

func (k Keeper) GetBorrowInterestRateDecrease(ctx sdk.Context) sdk.Dec {
	return k.GetParams(ctx).BorrowInterestRateDecrease
}

func (k Keeper) GetHealthGainFactor(ctx sdk.Context) sdk.Dec {
	return k.GetParams(ctx).HealthGainFactor
}

func (k Keeper) GetPoolOpenThreshold(ctx sdk.Context) sdk.Dec {
	return k.GetParams(ctx).PoolOpenThreshold
}

func (k Keeper) GetForceCloseFundPercentage(ctx sdk.Context) sdk.Dec {
	return k.GetParams(ctx).ForceCloseFundPercentage
}

func (k Keeper) GetForceCloseFundAddress(ctx sdk.Context) sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(k.GetParams(ctx).ForceCloseFundAddress)
	if err != nil {
		panic(err)
	}

	return addr
}

func (k Keeper) GetIncrementalBorrowInterestPaymentFundPercentage(ctx sdk.Context) sdk.Dec {
	return k.GetParams(ctx).IncrementalBorrowInterestPaymentFundPercentage
}

func (k Keeper) GetIncrementalBorrowInterestPaymentFundAddress(ctx sdk.Context) sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(k.GetParams(ctx).IncrementalBorrowInterestPaymentFundAddress)
	if err != nil {
		panic(err)
	}

	return addr
}

func (k Keeper) GetMaxOpenPositions(ctx sdk.Context) uint64 {
	return (uint64)(k.GetParams(ctx).MaxOpenPositions)
}

func (k Keeper) GetIncrementalBorrowInterestPaymentEnabled(ctx sdk.Context) bool {
	return k.GetParams(ctx).IncrementalBorrowInterestPaymentEnabled
}

func (k Keeper) GetSafetyFactor(ctx sdk.Context) sdk.Dec {
	return k.GetParams(ctx).SafetyFactor
}

func (k Keeper) GetEnabledPools(ctx sdk.Context) []uint64 {
	poolIds := make([]uint64, 0)
	pools := k.GetAllPools(ctx)
	for _, p := range pools {
		if p.Enabled {
			poolIds = append(poolIds, p.AmmPoolId)
		}
	}

	return poolIds
}

func (k Keeper) SetEnabledPools(ctx sdk.Context, pools []uint64) {
	for _, poolId := range pools {
		pool, found := k.GetPool(ctx, poolId)
		if !found {
			pool = types.NewPool(poolId)
			k.SetPool(ctx, pool)
		}
		pool.Enabled = true

		k.SetPool(ctx, pool)
	}
}

func (k Keeper) IsPoolEnabled(ctx sdk.Context, poolId uint64) bool {
	pool, found := k.GetPool(ctx, poolId)
	if !found {
		pool = types.NewPool(poolId)
		k.SetPool(ctx, pool)
	}

	return pool.Enabled
}

func (k Keeper) IsPoolClosed(ctx sdk.Context, poolId uint64) bool {
	pool, found := k.GetPool(ctx, poolId)
	if !found {
		pool = types.NewPool(poolId)
		k.SetPool(ctx, pool)
	}

	return pool.Closed
}

func (k Keeper) IsWhitelistingEnabled(ctx sdk.Context) bool {
	return k.GetParams(ctx).WhitelistingEnabled
}

func (k Keeper) GetTakeProfitBorrowInterestRateMin(ctx sdk.Context) sdk.Dec {
	return k.GetParams(ctx).TakeProfitBorrowInterestRateMin
}

func (k Keeper) GetFundingFeeBaseRate(ctx sdk.Context) sdk.Dec {
	return k.GetParams(ctx).FundingFeeBaseRate
}

func (k Keeper) GetFundingFeeMaxRate(ctx sdk.Context) sdk.Dec {
	return k.GetParams(ctx).FundingFeeMaxRate
}

func (k Keeper) GetFundingFeeMinRate(ctx sdk.Context) sdk.Dec {
	return k.GetParams(ctx).FundingFeeMinRate
}

func (k Keeper) GetFundingFeeCollectionAddress(ctx sdk.Context) sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(k.GetParams(ctx).FundingFeeCollectionAddress)
	if err != nil {
		panic(err)
	}

	return addr
}

func (k Keeper) GetSwapFee(ctx sdk.Context) sdk.Dec {
	return k.GetParams(ctx).SwapFee
}

func (k Keeper) GetEpochLength(ctx sdk.Context) int64 {
	return k.GetParams(ctx).EpochLength
}
