package types

import sdk "github.com/cosmos/cosmos-sdk/types"

type MarginHooks interface {
	// AfterMarginPositionOpended is called after OpenLong or OpenShort position.
	AfterMarginPositionOpended(ctx sdk.Context, poolId uint64)

	// AfterMarginPositionModified is called after a position gets modified.
	AfterMarginPositionModified(ctx sdk.Context, poolId uint64)

	// AfterMarginPositionClosed is called after a position gets closed.
	AfterMarginPositionClosed(ctx sdk.Context, poolId uint64)
}

var _ MarginHooks = MultiMarginHooks{}

// combine multiple margin hooks, all hook functions are run in array sequence.
type MultiMarginHooks []MarginHooks

// Creates hooks for the Amm Module.
func NewMultiMarginHooks(hooks ...MarginHooks) MultiMarginHooks {
	return hooks
}

func (h MultiMarginHooks) AfterMarginPositionOpended(ctx sdk.Context, poolId uint64) {
	for i := range h {
		h[i].AfterMarginPositionOpended(ctx, poolId)
	}
}

func (h MultiMarginHooks) AfterMarginPositionModified(ctx sdk.Context, poolId uint64) {
	for i := range h {
		h[i].AfterMarginPositionModified(ctx, poolId)
	}
}

func (h MultiMarginHooks) AfterMarginPositionClosed(ctx sdk.Context, poolId uint64) {
	for i := range h {
		h[i].AfterMarginPositionClosed(ctx, poolId)
	}
}
