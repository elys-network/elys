package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
)

type MarginHooks interface {
	// AfterMarginPositionOpended is called after OpenLong or OpenShort position.
	AfterMarginPositionOpended(ctx sdk.Context, ammPool ammtypes.Pool, marginPool Pool)

	// AfterMarginPositionModified is called after a position gets modified.
	AfterMarginPositionModified(ctx sdk.Context, ammPool ammtypes.Pool, marginPool Pool)

	// AfterMarginPositionClosed is called after a position gets closed.
	AfterMarginPositionClosed(ctx sdk.Context, ammPool ammtypes.Pool, marginPool Pool)
}

var _ MarginHooks = MultiMarginHooks{}

// combine multiple margin hooks, all hook functions are run in array sequence.
type MultiMarginHooks []MarginHooks

// Creates hooks for the Amm Module.
func NewMultiMarginHooks(hooks ...MarginHooks) MultiMarginHooks {
	return hooks
}

func (h MultiMarginHooks) AfterMarginPositionOpended(ctx sdk.Context, ammPool ammtypes.Pool, marginPool Pool) {
	for i := range h {
		h[i].AfterMarginPositionOpended(ctx, ammPool, marginPool)
	}
}

func (h MultiMarginHooks) AfterMarginPositionModified(ctx sdk.Context, ammPool ammtypes.Pool, marginPool Pool) {
	for i := range h {
		h[i].AfterMarginPositionModified(ctx, ammPool, marginPool)
	}
}

func (h MultiMarginHooks) AfterMarginPositionClosed(ctx sdk.Context, ammPool ammtypes.Pool, marginPool Pool) {
	for i := range h {
		h[i].AfterMarginPositionClosed(ctx, ammPool, marginPool)
	}
}
