package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	margintypes "github.com/elys-network/elys/x/margin/types"
)

func (k Keeper) AfterMarginPositionOpended(ctx sdk.Context, ammPool ammtypes.Pool, marginPool margintypes.Pool) {
	k.UpdateAccountedPool(ctx, ammPool, marginPool)
}

func (k Keeper) AfterMarginPositionModified(ctx sdk.Context, ammPool ammtypes.Pool, marginPool margintypes.Pool) {
	k.UpdateAccountedPool(ctx, ammPool, marginPool)
}

func (k Keeper) AfterMarginPositionClosed(ctx sdk.Context, ammPool ammtypes.Pool, marginPool margintypes.Pool) {
	k.UpdateAccountedPool(ctx, ammPool, marginPool)
}

// AfterPoolCreated is called after CreatePool
func (k Keeper) AfterAmmPoolCreated(ctx sdk.Context, ammPool ammtypes.Pool) {
	k.InitiateAccountedPool(ctx, ammPool)
}

// AfterJoinPool is called after JoinPool, JoinSwapExternAmountIn, and JoinSwapShareAmountOut
func (k Keeper) AfterAmmJoinPool(ctx sdk.Context, ammPool ammtypes.Pool, marginPool margintypes.Pool) {
	k.UpdateAccountedPool(ctx, ammPool, marginPool)
}

// AfterExitPool is called after ExitPool, ExitSwapShareAmountIn, and ExitSwapExternAmountOut
func (k Keeper) AfterAmmExitPool(ctx sdk.Context, ammPool ammtypes.Pool, marginPool margintypes.Pool) {
	k.UpdateAccountedPool(ctx, ammPool, marginPool)
}

// AfterSwap is called after SwapExactAmountIn and SwapExactAmountOut
func (k Keeper) AfterAmmSwap(ctx sdk.Context, ammPool ammtypes.Pool, marginPool margintypes.Pool) {
	k.UpdateAccountedPool(ctx, ammPool, marginPool)
}

// Hooks wrapper struct for tvl keeper
type MarginHooks struct {
	k Keeper
}

var _ margintypes.MarginHooks = MarginHooks{}

// Return the wrapper struct
func (k Keeper) MarginHooks() MarginHooks {
	return MarginHooks{k}
}

func (h MarginHooks) AfterMarginPositionOpended(ctx sdk.Context, ammPool ammtypes.Pool, marginPool margintypes.Pool) {
	h.k.AfterMarginPositionOpended(ctx, ammPool, marginPool)
}

func (h MarginHooks) AfterMarginPositionModified(ctx sdk.Context, ammPool ammtypes.Pool, marginPool margintypes.Pool) {
	h.k.AfterMarginPositionModified(ctx, ammPool, marginPool)
}

func (h MarginHooks) AfterMarginPositionClosed(ctx sdk.Context, ammPool ammtypes.Pool, marginPool margintypes.Pool) {
	h.k.AfterMarginPositionClosed(ctx, ammPool, marginPool)
}

func (h MarginHooks) AfterAmmPoolCreated(ctx sdk.Context, ammPool ammtypes.Pool) {
	h.k.AfterAmmPoolCreated(ctx, ammPool)
}

func (h MarginHooks) AfterAmmJoinPool(ctx sdk.Context, ammPool ammtypes.Pool, marginPool margintypes.Pool) {
	h.k.AfterAmmJoinPool(ctx, ammPool, marginPool)
}

func (h MarginHooks) AfterAmmExitPool(ctx sdk.Context, ammPool ammtypes.Pool, marginPool margintypes.Pool) {
	h.k.AfterAmmExitPool(ctx, ammPool, marginPool)
}

func (h MarginHooks) AfterAmmSwap(ctx sdk.Context, ammPool ammtypes.Pool, marginPool margintypes.Pool) {
	h.k.AfterAmmSwap(ctx, ammPool, marginPool)
}
