package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/burner/types"
)

// GetParams get all parameters as types.Params
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	return types.NewParams(
		k.EpochIdentifier(ctx),
	)
}

// SetParams set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}

// EpochIdentifier returns the EpochIdentifier param
func (k Keeper) EpochIdentifier(ctx sdk.Context) (res string) {
	k.paramstore.Get(ctx, types.KeyEpochIdentifier, &res)
	return
}
