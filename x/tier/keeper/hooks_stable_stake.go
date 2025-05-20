package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stablestaketypes "github.com/elys-network/elys/v4/x/stablestake/types"
)

// Hooks wrapper struct for incentive keeper
type StableStakeHooks struct {
	k Keeper
}

var _ stablestaketypes.StableStakeHooks = StableStakeHooks{}

// Return the wrapper struct
func (k Keeper) StableStakeHooks() StableStakeHooks {
	return StableStakeHooks{k}
}

func (h StableStakeHooks) AfterBond(ctx sdk.Context, sender sdk.AccAddress, shareAmount math.Int, poolId uint64) error {
	h.k.RetrieveAllPortfolio(ctx, sender)
	return nil
}

func (h StableStakeHooks) AfterUnbond(ctx sdk.Context, sender sdk.AccAddress, shareAmount math.Int, poolId uint64) error {
	h.k.RetrieveAllPortfolio(ctx, sender)
	return nil
}
