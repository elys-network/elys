package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AfterEpochEnd executes the indicated hook after epochs ends
func (k Keeper) RunHooksAfterEpochEnd(ctx sdk.Context, identifier string, epochNumber int64) error {
	if k.hooks == nil {
		return nil
	}
	return k.hooks.AfterEpochEnd(ctx, identifier, epochNumber)
}

// BeforeEpochStart executes the indicated hook before the epochs
func (k Keeper) RunHooksBeforeEpochStart(ctx sdk.Context, identifier string, epochNumber int64) error {
	if k.hooks == nil {
		return nil
	}

	return k.hooks.BeforeEpochStart(ctx, identifier, epochNumber)
}
