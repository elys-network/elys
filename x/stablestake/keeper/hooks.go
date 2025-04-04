package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/stablestake/types"
)

var _ types.StableStakeHooks = MultiStableStakeHooks{}

// combine multiple stablestake hooks, all hook functions are run in array sequence
type MultiStableStakeHooks []types.StableStakeHooks

func NewMultiStableStakeHooks(hooks ...types.StableStakeHooks) MultiStableStakeHooks {
	return hooks
}

// Committed is called when staker committed his token
func (mh MultiStableStakeHooks) AfterBond(ctx sdk.Context, sender sdk.AccAddress, shareAmount math.Int, poolId uint64) error {
	for i := range mh {
		err := mh[i].AfterBond(ctx, sender, shareAmount, poolId)
		if err != nil {
			return err
		}
	}
	return nil
}

// Committed is called when staker committed his token
func (mh MultiStableStakeHooks) AfterUnbond(ctx sdk.Context, sender sdk.AccAddress, shareAmount math.Int, poolId uint64) error {
	for i := range mh {
		err := mh[i].AfterUnbond(ctx, sender, shareAmount, poolId)
		if err != nil {
			return err
		}
	}
	return nil
}
