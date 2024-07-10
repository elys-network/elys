package types

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// StableStakeHooks event hooks for stablestake processing
type StableStakeHooks interface {
	AfterBond(ctx sdk.Context, sender string, shareAmount math.Int) error
	AfterUnbond(ctx sdk.Context, sender string, shareAmount math.Int) error
	AfterUpdateInterestStacked(ctx sdk.Context, address string, old math.Int, new math.Int) error
}
