package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/commitment/types"
)

// GetParams returns the current parameters of the Commitment module
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	var params types.Params
	k.paramstore.GetParamSet(ctx, &params)
	return params
}

// SetParams set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}

// GetVestingDenom returns the vesting denom for the given base denom
func (k Keeper) GetVestingInfo(ctx sdk.Context, baseDenom string) (*types.VestingInfo, int) {
	params := k.GetParams(ctx)

	for i, vestingInfo := range params.VestingInfos {
		if vestingInfo.BaseDenom == baseDenom {
			return vestingInfo, i
		}
	}

	return nil, 0
}
