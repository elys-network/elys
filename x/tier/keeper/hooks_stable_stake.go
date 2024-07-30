package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stablestaketypes "github.com/elys-network/elys/x/stablestake/types"
)

func (k Keeper) AfterBond(ctx sdk.Context, sender string, shareAmount math.Int) error {
	k.RetrieveAllPortfolio(ctx, sender)
	return nil
}

func (k Keeper) AfterUnbond(ctx sdk.Context, sender string, shareAmount math.Int) error {
	k.RetrieveAllPortfolio(ctx, sender)
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
	h.k.AfterBond(ctx, sender, shareAmount)
	return nil
}

func (h StableStakeHooks) AfterUnbond(ctx sdk.Context, sender string, shareAmount math.Int) error {
	h.k.AfterUnbond(ctx, sender, shareAmount)
	return nil
}
