package keeper

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/leveragelp/types"
	stablestaketypes "github.com/elys-network/elys/x/stablestake/types"
)

func (k Keeper) AfterBond(ctx sdk.Context, sender sdk.AccAddress, shareAmount math.Int) error {
	return nil
}

func (k Keeper) AfterUnbond(ctx sdk.Context, sender sdk.AccAddress, shareAmount math.Int) error {
	// loop over all leverage pools and throw error if any leverage pool is unhealthy
	for _, pool := range k.GetAllPools(ctx) {
		// print that after unbound is called and checking pool health of pool id
		fmt.Printf("after unbound is called and checking pool health of pool id %d\n", pool.AmmPoolId)
		if err := k.CheckPoolHealth(ctx, pool.AmmPoolId); err != nil {
			return errorsmod.Wrapf(types.ErrUnbondingPoolHealth, "pool health too low to unbond for pool %d", pool.AmmPoolId)
		}
	}

	return nil
}

type StableStakeHooks struct {
	k Keeper
}

var _ stablestaketypes.StableStakeHooks = StableStakeHooks{}

// Return the wrapper struct
func (k Keeper) StableStakeHooks() StableStakeHooks {
	return StableStakeHooks{k}
}

// AfterPoolCreated is called after CreatePool
func (h StableStakeHooks) AfterBond(ctx sdk.Context, sender sdk.AccAddress, shareAmount math.Int) error {
	return h.k.AfterBond(ctx, sender, shareAmount)
}

// AfterJoinPool is called after JoinPool, JoinSwapExternAmountIn, and JoinSwapShareAmountOut
func (h StableStakeHooks) AfterUnbond(ctx sdk.Context, sender sdk.AccAddress, shareAmount math.Int) error {
	return h.k.AfterUnbond(ctx, sender, shareAmount)
}
