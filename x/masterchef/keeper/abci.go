package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// EndBlocker of amm module
func (k Keeper) EndBlocker(ctx sdk.Context) {
	// TODO: reward distributions

	// TODO: calculate APR

	// TODO: build poolInfo object

	// TODO: remove expired external incentives

	// TODO: manage pool external_reward_denoms array
}
