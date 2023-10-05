package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
)

type LeveragelpHooks interface {
	// AfterLeveragelpPositionOpended is called after OpenLong or OpenShort position.
	AfterLeveragelpPositionOpended(ctx sdk.Context, ammPool ammtypes.Pool, leveragelpPool Pool)

	// AfterLeveragelpPositionModified is called after a position gets modified.
	AfterLeveragelpPositionModified(ctx sdk.Context, ammPool ammtypes.Pool, leveragelpPool Pool)

	// AfterLeveragelpPositionClosed is called after a position gets closed.
	AfterLeveragelpPositionClosed(ctx sdk.Context, ammPool ammtypes.Pool, leveragelpPool Pool)

	// AfterPoolCreated is called after CreatePool
	AfterAmmPoolCreated(ctx sdk.Context, ammPool ammtypes.Pool)

	// AfterJoinPool is called after JoinPool, JoinSwapExternAmountIn, and JoinSwapShareAmountOut
	AfterAmmJoinPool(ctx sdk.Context, ammPool ammtypes.Pool, leveragelpPool Pool)

	// AfterExitPool is called after ExitPool, ExitSwapShareAmountIn, and ExitSwapExternAmountOut
	AfterAmmExitPool(ctx sdk.Context, ammPool ammtypes.Pool, leveragelpPool Pool)

	// AfterSwap is called after SwapExactAmountIn and SwapExactAmountOut
	AfterAmmSwap(ctx sdk.Context, ammPool ammtypes.Pool, leveragelpPool Pool)
}

var _ LeveragelpHooks = MultiLeveragelpHooks{}

// combine multiple leveragelp hooks, all hook functions are run in array sequence.
type MultiLeveragelpHooks []LeveragelpHooks

// Creates hooks for the Amm Module.
func NewMultiLeveragelpHooks(hooks ...LeveragelpHooks) MultiLeveragelpHooks {
	return hooks
}

func (h MultiLeveragelpHooks) AfterLeveragelpPositionOpended(ctx sdk.Context, ammPool ammtypes.Pool, leveragelpPool Pool) {
	for i := range h {
		h[i].AfterLeveragelpPositionOpended(ctx, ammPool, leveragelpPool)
	}
}

func (h MultiLeveragelpHooks) AfterLeveragelpPositionModified(ctx sdk.Context, ammPool ammtypes.Pool, leveragelpPool Pool) {
	for i := range h {
		h[i].AfterLeveragelpPositionModified(ctx, ammPool, leveragelpPool)
	}
}

func (h MultiLeveragelpHooks) AfterLeveragelpPositionClosed(ctx sdk.Context, ammPool ammtypes.Pool, leveragelpPool Pool) {
	for i := range h {
		h[i].AfterLeveragelpPositionClosed(ctx, ammPool, leveragelpPool)
	}
}

func (h MultiLeveragelpHooks) AfterAmmPoolCreated(ctx sdk.Context, ammPool ammtypes.Pool) {
	for i := range h {
		h[i].AfterAmmPoolCreated(ctx, ammPool)
	}
}

func (h MultiLeveragelpHooks) AfterAmmJoinPool(ctx sdk.Context, ammPool ammtypes.Pool, leveragelpPool Pool) {
	for i := range h {
		h[i].AfterAmmJoinPool(ctx, ammPool, leveragelpPool)
	}
}

func (h MultiLeveragelpHooks) AfterAmmExitPool(ctx sdk.Context, ammPool ammtypes.Pool, leveragelpPool Pool) {
	for i := range h {
		h[i].AfterAmmExitPool(ctx, ammPool, leveragelpPool)
	}
}

func (h MultiLeveragelpHooks) AfterAmmSwap(ctx sdk.Context, ammPool ammtypes.Pool, leveragelpPool Pool) {
	for i := range h {
		h[i].AfterAmmSwap(ctx, ammPool, leveragelpPool)
	}
}
