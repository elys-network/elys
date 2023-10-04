package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/leveragelp/types"
)

// SetParams set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) error {
	if err := params.Validate(); err != nil {
		return err
	}

	store := ctx.KVStore(k.storeKey)
	bz, err := k.cdc.Marshal(&params)
	if err != nil {
		return err
	}
	store.Set(types.KeyPrefix(types.ParamsKey), bz)

	return nil
}

func (k Keeper) GetMaxLeverageParam(ctx sdk.Context) sdk.Dec {
	return k.GetParams(ctx).LeverageMax
}

func (k Keeper) GetInterestRateMax(ctx sdk.Context) sdk.Dec {
	return k.GetParams(ctx).InterestRateMax
}

func (k Keeper) GetInterestRateMin(ctx sdk.Context) sdk.Dec {
	return k.GetParams(ctx).InterestRateMin
}

func (k Keeper) GetInterestRateIncrease(ctx sdk.Context) sdk.Dec {
	return k.GetParams(ctx).InterestRateIncrease
}

func (k Keeper) GetInterestRateDecrease(ctx sdk.Context) sdk.Dec {
	return k.GetParams(ctx).InterestRateDecrease
}

func (k Keeper) GetHealthGainFactor(ctx sdk.Context) sdk.Dec {
	return k.GetParams(ctx).HealthGainFactor
}

func (k Keeper) GetPoolOpenThreshold(ctx sdk.Context) sdk.Dec {
	return k.GetParams(ctx).PoolOpenThreshold
}

func (k Keeper) GetRemovalQueueThreshold(ctx sdk.Context) sdk.Dec {
	return k.GetParams(ctx).RemovalQueueThreshold
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

func (k Keeper) GetIncrementalInterestPaymentFundPercentage(ctx sdk.Context) sdk.Dec {
	return k.GetParams(ctx).IncrementalInterestPaymentFundPercentage
}

func (k Keeper) GetIncrementalInterestPaymentFundAddress(ctx sdk.Context) sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(k.GetParams(ctx).IncrementalInterestPaymentFundAddress)
	if err != nil {
		panic(err)
	}

	return addr
}

func (k Keeper) GetMaxOpenPositions(ctx sdk.Context) uint64 {
	return (uint64)(k.GetParams(ctx).MaxOpenPositions)
}

func (k Keeper) GetIncrementalInterestPaymentEnabled(ctx sdk.Context) bool {
	return k.GetParams(ctx).IncrementalInterestPaymentEnabled
}
func (k Keeper) GetSafetyFactor(ctx sdk.Context) sdk.Dec {
	return k.GetParams(ctx).SafetyFactor
}

func (k Keeper) GetSqModifier(ctx sdk.Context) sdk.Dec {
	return k.GetParams(ctx).SqModifier
}

func (k Keeper) IsWhitelistingEnabled(ctx sdk.Context) bool {
	return k.GetParams(ctx).WhitelistingEnabled
}
