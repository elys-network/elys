package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	perpetualtypes "github.com/elys-network/elys/x/perpetual/types"
)

func (k Keeper) AfterPerpetualPositionOpen(ctx sdk.Context, ammPool ammtypes.Pool, perpetualPool perpetualtypes.Pool, sender string) {
	k.UpdateAccountedPool(ctx, ammPool, perpetualPool)
}

func (k Keeper) AfterPerpetualPositionModified(ctx sdk.Context, ammPool ammtypes.Pool, perpetualPool perpetualtypes.Pool, sender string) {
	k.UpdateAccountedPool(ctx, ammPool, perpetualPool)
}

func (k Keeper) AfterPerpetualPositionClosed(ctx sdk.Context, ammPool ammtypes.Pool, perpetualPool perpetualtypes.Pool, sender string) {
	k.UpdateAccountedPool(ctx, ammPool, perpetualPool)
}

// AfterPoolCreated is called after CreatePool
func (k Keeper) AfterAmmPoolCreated(ctx sdk.Context, ammPool ammtypes.Pool, sender string) {
	k.InitiateAccountedPool(ctx, ammPool)
}

// AfterAmmJoinPool is called after JoinPool, JoinSwapExternAmountIn, and JoinSwapShareAmountOut
func (k Keeper) AfterAmmJoinPool(ctx sdk.Context, ammPool ammtypes.Pool, perpetualPool perpetualtypes.Pool, sender string) {
	k.UpdateAccountedPool(ctx, ammPool, perpetualPool)
}

// AfterExitPool is called after ExitPool, ExitSwapShareAmountIn, and ExitSwapExternAmountOut
func (k Keeper) AfterAmmExitPool(ctx sdk.Context, ammPool ammtypes.Pool, perpetualPool perpetualtypes.Pool, sender string) {
	k.UpdateAccountedPool(ctx, ammPool, perpetualPool)
}

// AfterSwap is called after SwapExactAmountIn and SwapExactAmountOut
func (k Keeper) AfterAmmSwap(ctx sdk.Context, ammPool ammtypes.Pool, perpetualPool perpetualtypes.Pool, sender string) {
	k.UpdateAccountedPool(ctx, ammPool, perpetualPool)
}

// Hooks wrapper struct for tvl keeper
type PerpetualHooks struct {
	k Keeper
}

var _ perpetualtypes.PerpetualHooks = PerpetualHooks{}

// Return the wrapper struct
func (k Keeper) PerpetualHooks() PerpetualHooks {
	return PerpetualHooks{k}
}

func (h PerpetualHooks) AfterPerpetualPositionOpen(ctx sdk.Context, ammPool ammtypes.Pool, perpetualPool perpetualtypes.Pool, sender string) {
	h.k.AfterPerpetualPositionOpen(ctx, ammPool, perpetualPool, sender)
}

func (h PerpetualHooks) AfterPerpetualPositionModified(ctx sdk.Context, ammPool ammtypes.Pool, perpetualPool perpetualtypes.Pool, sender string) {
	h.k.AfterPerpetualPositionModified(ctx, ammPool, perpetualPool, sender)
}

func (h PerpetualHooks) AfterPerpetualPositionClosed(ctx sdk.Context, ammPool ammtypes.Pool, perpetualPool perpetualtypes.Pool, sender string) {
	h.k.AfterPerpetualPositionClosed(ctx, ammPool, perpetualPool, sender)
}

func (h PerpetualHooks) AfterAmmPoolCreated(ctx sdk.Context, ammPool ammtypes.Pool, sender string) {
	h.k.AfterAmmPoolCreated(ctx, ammPool, sender)
}

func (h PerpetualHooks) AfterAmmJoinPool(ctx sdk.Context, ammPool ammtypes.Pool, perpetualPool perpetualtypes.Pool, sender string) {
	h.k.AfterAmmJoinPool(ctx, ammPool, perpetualPool, sender)
}

func (h PerpetualHooks) AfterAmmExitPool(ctx sdk.Context, ammPool ammtypes.Pool, perpetualPool perpetualtypes.Pool, sender string) {
	h.k.AfterAmmExitPool(ctx, ammPool, perpetualPool, sender)
}

func (h PerpetualHooks) AfterAmmSwap(ctx sdk.Context, ammPool ammtypes.Pool, perpetualPool perpetualtypes.Pool, sender string) {
	h.k.AfterAmmSwap(ctx, ammPool, perpetualPool, sender)
}
