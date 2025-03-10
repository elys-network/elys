package types

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// StableStakeHooks event hooks for stablestake processing
type StableStakeHooks interface {
	AfterBond(ctx sdk.Context, sender sdk.AccAddress, shareAmount math.Int, poolId uint64) error
	AfterUnbond(ctx sdk.Context, sender sdk.AccAddress, shareAmount math.Int, poolId uint64) error
}
