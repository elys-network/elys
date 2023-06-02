package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/incentive/types"
)

// GetParams get all parameters as types.Params
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	var params types.Params
	k.paramstore.GetParamSet(ctx, &params)
	return params
}

// SetParams set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}

// GetCommunityTax returns the current distribution community tax.
func (k Keeper) GetCommunityTax(ctx sdk.Context) (percent sdk.Dec) {
	k.paramstore.Get(ctx, types.ParamStoreKeyCommunityTax, &percent)
	return percent
}

// GetWithdrawAddrEnabled returns the current distribution withdraw address
// enabled parameter.
func (k Keeper) GetWithdrawAddrEnabled(ctx sdk.Context) (enabled bool) {
	k.paramstore.Get(ctx, types.ParamStoreKeyWithdrawAddrEnabled, &enabled)
	return enabled
}

// GetDEXRewardPercentForLPs returns the dex revenue percent for Lps
func (k Keeper) GetDEXRewardPercentForLPs(ctx sdk.Context) (percent sdk.Dec) {
	k.paramstore.Get(ctx, types.ParamStoreKeyRewardPercentForLps, &percent)
	if percent.LTE(sdk.ZeroDec()) {
		return sdk.NewDecWithPrec(65, 1)
	}
	return percent
}
