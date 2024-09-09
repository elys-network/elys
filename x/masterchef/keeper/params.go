package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/masterchef/types"
)

// GetParams get all parameters as types.Params
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	store := ctx.KVStore(k.storeKey)

	b := store.Get(types.ParamsKey)
	if b == nil {
		return
	}

	k.cdc.MustUnmarshal(b, &params)
	return
}

// SetParams set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshal(&params)
	store.Set(types.ParamsKey, b)
}

func (k Keeper) CheckBlockedAddress(params types.Params) bool {
	return k.bankKeeper.BlockedAddr(sdk.MustAccAddressFromBech32(params.ProtocolRevenueAddress))
}

func (k Keeper) GetLegacyParams(ctx sdk.Context) (params types.LegacyParams) {
	store := ctx.KVStore(k.storeKey)

	b := store.Get(types.ParamsKey)
	if b == nil {
		return
	}

	k.cdc.MustUnmarshal(b, &params)
	return
}
