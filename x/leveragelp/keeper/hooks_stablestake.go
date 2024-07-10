package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stablestaketypes "github.com/elys-network/elys/x/stablestake/types"
)

func (k Keeper) AfterUpdateInterestStacked(ctx sdk.Context, address string, old sdk.Int, new sdk.Int) error {
	k.SetSortedLiquidation(ctx, address, old, new)
	return nil
}

// Hooks wrapper struct for incentive keeper
type StableStakeHooks struct {
	k Keeper
}

var _ stablestaketypes.StableStakeHooks = StableStakeHooks{}

// Return the wrapper struct
func (k Keeper) StableStakeHooks() StableStakeHooks {
	return StableStakeHooks{k}
}

func (h StableStakeHooks) AfterBond(ctx sdk.Context, sender string, shareAmount math.Int) error {
	return nil
}

func (h StableStakeHooks) AfterUnbond(ctx sdk.Context, sender string, shareAmount math.Int) error {
	return nil
}

func (h StableStakeHooks) AfterUpdateInterestStacked(ctx sdk.Context, address string, old sdk.Int, new sdk.Int) error {
	h.k.AfterUpdateInterestStacked(ctx, address, old, new)
	return nil
}
