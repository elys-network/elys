package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
)

// AfterPoolCreated is called after CreatePool
func (k Keeper) AfterPoolCreated(ctx sdk.Context, sender sdk.AccAddress, poolId uint64) {
	k.InitiateAccountedPool(ctx, poolId)
}

// AfterJoinPool is called after JoinPool, JoinSwapExternAmountIn, and JoinSwapShareAmountOut
func (k Keeper) AfterJoinPool(ctx sdk.Context, sender sdk.AccAddress, poolId uint64, enterCoins sdk.Coins, shareOutAmount sdk.Int) {
	k.UpdateAccountedPool(ctx, poolId)
}

// AfterExitPool is called after ExitPool, ExitSwapShareAmountIn, and ExitSwapExternAmountOut
func (k Keeper) AfterExitPool(ctx sdk.Context, sender sdk.AccAddress, poolId uint64, shareInAmount sdk.Int, exitCoins sdk.Coins) {
	k.UpdateAccountedPool(ctx, poolId)
}

// AfterSwap is called after SwapExactAmountIn and SwapExactAmountOut
func (k Keeper) AfterSwap(ctx sdk.Context, sender sdk.AccAddress, poolId uint64, input sdk.Coins, output sdk.Coins) {
	k.UpdateAccountedPool(ctx, poolId)
}

// Hooks wrapper struct for tvl keeper
type AmmHooks struct {
	k Keeper
}

var _ ammtypes.AmmHooks = AmmHooks{}

// Return the wrapper struct
func (k Keeper) AmmHooks() AmmHooks {
	return AmmHooks{k}
}

// AfterPoolCreated is called after CreatePool
func (h AmmHooks) AfterPoolCreated(ctx sdk.Context, sender sdk.AccAddress, poolId uint64) {
	h.k.AfterPoolCreated(ctx, sender, poolId)
}

// AfterJoinPool is called after JoinPool, JoinSwapExternAmountIn, and JoinSwapShareAmountOut
func (h AmmHooks) AfterJoinPool(ctx sdk.Context, sender sdk.AccAddress, poolId uint64, enterCoins sdk.Coins, shareOutAmount sdk.Int) {
	h.k.AfterJoinPool(ctx, sender, poolId, enterCoins, shareOutAmount)
}

// AfterExitPool is called after ExitPool, ExitSwapShareAmountIn, and ExitSwapExternAmountOut
func (h AmmHooks) AfterExitPool(ctx sdk.Context, sender sdk.AccAddress, poolId uint64, shareInAmount sdk.Int, exitCoins sdk.Coins) {
	h.k.AfterExitPool(ctx, sender, poolId, shareInAmount, exitCoins)
}

// AfterSwap is called after SwapExactAmountIn and SwapExactAmountOut
func (h AmmHooks) AfterSwap(ctx sdk.Context, sender sdk.AccAddress, poolId uint64, input sdk.Coins, output sdk.Coins) {
	h.k.AfterSwap(ctx, sender, poolId, input, output)
}
