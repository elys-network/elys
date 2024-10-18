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

func (k Keeper) IsWhitelistingEnabled(ctx sdk.Context) bool {
	return k.GetParams(ctx).WhitelistingEnabled
}

func (k Keeper) GetTakeProfitBorrowInterestRateMin(ctx sdk.Context) sdk.Dec {
	return k.GetParams(ctx).TakeProfitBorrowInterestRateMin
}

func (k Keeper) GetPerpetualSwapFee(ctx sdk.Context) sdk.Dec {
	return k.GetParams(ctx).PerpetualSwapFee
}

func (k Keeper) GetMaxLimitOrder(ctx sdk.Context) int64 {
	return k.GetParams(ctx).MaxLimitOrder
}
