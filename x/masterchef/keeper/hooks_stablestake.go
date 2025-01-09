package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stablestaketypes "github.com/elys-network/elys/x/stablestake/types"
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
	h.k.AfterDeposit(ctx, poolId, sender, shareAmount)
	return nil
}

func (h StableStakeHooks) AfterUnbond(ctx sdk.Context, sender sdk.AccAddress, shareAmount math.Int, poolId uint64) error {
	h.k.AfterWithdraw(ctx, poolId, sender, shareAmount)
	return nil
}
