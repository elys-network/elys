package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
)

type MarginHooks interface {
	// AfterMarginPositionOpen is called after OpenLong or OpenShort position.
	AfterMarginPositionOpen(ctx sdk.Context, ammPool ammtypes.Pool, marginPool Pool)

	// AfterMarginPositionModified is called after a position gets modified.
	AfterMarginPositionModified(ctx sdk.Context, ammPool ammtypes.Pool, marginPool Pool)

	// AfterMarginPositionClosed is called after a position gets closed.
	AfterMarginPositionClosed(ctx sdk.Context, ammPool ammtypes.Pool, marginPool Pool)

	// AfterPoolCreated is called after CreatePool
	AfterAmmPoolCreated(ctx sdk.Context, ammPool ammtypes.Pool)

	// AfterJoinPool is called after JoinPool, JoinSwapExternAmountIn, and JoinSwapShareAmountOut
	AfterAmmJoinPool(ctx sdk.Context, ammPool ammtypes.Pool, marginPool Pool)

	// AfterExitPool is called after ExitPool, ExitSwapShareAmountIn, and ExitSwapExternAmountOut
	AfterAmmExitPool(ctx sdk.Context, ammPool ammtypes.Pool, marginPool Pool)

	// AfterSwap is called after SwapExactAmountIn and SwapExactAmountOut
	AfterAmmSwap(ctx sdk.Context, ammPool ammtypes.Pool, marginPool Pool)
}

var _ MarginHooks = MultiMarginHooks{}

// combine multiple margin hooks, all hook functions are run in array sequence.
type MultiMarginHooks []MarginHooks

// Creates hooks for the Amm Module.
func NewMultiMarginHooks(hooks ...MarginHooks) MultiMarginHooks {
	return hooks
}

func (h MultiMarginHooks) AfterMarginPositionOpen(ctx sdk.Context, ammPool ammtypes.Pool, marginPool Pool) {
	for i := range h {
		h[i].AfterMarginPositionOpen(ctx, ammPool, marginPool)
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

func (h MultiMarginHooks) AfterAmmPoolCreated(ctx sdk.Context, ammPool ammtypes.Pool) {
	for i := range h {
		h[i].AfterAmmPoolCreated(ctx, ammPool)
	}
}

func (h MultiMarginHooks) AfterAmmJoinPool(ctx sdk.Context, ammPool ammtypes.Pool, marginPool Pool) {
	for i := range h {
		h[i].AfterAmmJoinPool(ctx, ammPool, marginPool)
	}
}

func (h MultiMarginHooks) AfterAmmExitPool(ctx sdk.Context, ammPool ammtypes.Pool, marginPool Pool) {
	for i := range h {
		h[i].AfterAmmExitPool(ctx, ammPool, marginPool)
	}
}

func (h MultiMarginHooks) AfterAmmSwap(ctx sdk.Context, ammPool ammtypes.Pool, marginPool Pool) {
	for i := range h {
		h[i].AfterAmmSwap(ctx, ammPool, marginPool)
	}
}
