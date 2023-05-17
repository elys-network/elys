package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
)

// GetParams get all parameters as types.Params
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	return types.NewParams(
		k.PoolCreationFee(ctx),
	)
}

// SetParams set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}

// PoolCreationFee returns the PoolCreationFee param
func (k Keeper) PoolCreationFee(ctx sdk.Context) (res uint64) {
	k.paramstore.Get(ctx, types.KeyPoolCreationFee, &res)
	return
}
