package keeper

import (
	"cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/commitment/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

// GetParams get all parameters as types.Params
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))

	b := store.Get(types.ParamsKey)
	if b == nil {
		return
	}

	k.cdc.MustUnmarshal(b, &params)
	return
}

// GetLegacyParams get all parameters as types.Params
func (k Keeper) GetLegacyParams(ctx sdk.Context) (params types.LegacyParams) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))

	b := store.Get(types.ParamsKey)
	if b == nil {
		return
	}

	k.cdc.MustUnmarshal(b, &params)
	return
}

// SetParams set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	b := k.cdc.MustMarshal(&params)
	store.Set(types.ParamsKey, b)
}

// GetVestingDenom returns the vesting denom for the given base denom
func (k Keeper) GetVestingInfo(ctx sdk.Context, baseDenom string) (*types.VestingInfo, int) {
	params := k.GetParams(ctx)

	for i, vestingInfo := range params.VestingInfos {
		if vestingInfo.BaseDenom == baseDenom {
			return &vestingInfo, i
		}
	}

	return nil, 0
}

func (k Keeper) V9_ParamsMigration(ctx sdk.Context) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	iterator := storetypes.KVStorePrefixIterator(store, types.CommitmentsKeyPrefix)

	defer iterator.Close()

	totalEden := math.ZeroInt()
	totalEdenB := math.ZeroInt()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Commitments
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		for _, token := range val.CommittedTokens {
			if token.Denom == ptypes.Eden {
				totalEden.Add(token.Amount)
			}
			if token.Denom == ptypes.EdenB {
				totalEdenB.Add(token.Amount)
			}
		}

		totalEden.Add(val.Claimed.AmountOf(ptypes.Eden))
		totalEdenB.Add(val.Claimed.AmountOf(ptypes.EdenB))
	}

	params := k.GetParams(ctx)
	params.TotalEdenSupply = totalEden
	params.TotalEdenbSupply = totalEdenB
	k.SetParams(ctx, params)
}
