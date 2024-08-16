package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
)

// AfterPoolCreated is called after CreatePool
func (k Keeper) AfterPoolCreated(ctx sdk.Context, sender sdk.AccAddress, ammPool ammtypes.Pool) error {
	if k.hooks != nil {
		k.hooks.AfterAmmPoolCreated(ctx, ammPool, sender)
	}
	return nil
}

// AfterJoinPool is called after JoinPool, JoinSwapExternAmountIn, and JoinSwapShareAmountOut
func (k Keeper) AfterJoinPool(ctx sdk.Context, sender sdk.AccAddress, ammPool ammtypes.Pool, enterCoins sdk.Coins, shareOutAmount math.Int) error {
	perpetualPool, found := k.GetPool(ctx, ammPool.PoolId)
	if !found {
		return nil
	}

	if k.hooks != nil {
		k.hooks.AfterAmmJoinPool(ctx, ammPool, perpetualPool, sender)
	}
	return nil
}

// AfterExitPool is called after ExitPool, ExitSwapShareAmountIn, and ExitSwapExternAmountOut
func (k Keeper) AfterExitPool(ctx sdk.Context, sender sdk.AccAddress, ammPool ammtypes.Pool, shareInAmount math.Int, exitCoins sdk.Coins) error {
	perpetualPool, found := k.GetPool(ctx, ammPool.PoolId)
	if !found {
		return nil
	}

	if k.hooks != nil {
		k.hooks.AfterAmmExitPool(ctx, ammPool, perpetualPool, sender)
	}
	return nil
}

// AfterSwap is called after SwapExactAmountIn and SwapExactAmountOut
func (k Keeper) AfterSwap(ctx sdk.Context, sender sdk.AccAddress, ammPool ammtypes.Pool, input sdk.Coins, output sdk.Coins) error {
	perpetualPool, found := k.GetPool(ctx, ammPool.PoolId)
	if !found {
		return nil
	}
	if k.hooks != nil {
		k.hooks.AfterAmmSwap(ctx, ammPool, perpetualPool, sender)
	}
	return nil
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
func (h AmmHooks) AfterPoolCreated(ctx sdk.Context, sender sdk.AccAddress, pool ammtypes.Pool) error {
	return h.k.AfterPoolCreated(ctx, sender, pool)
}

// AfterJoinPool is called after JoinPool, JoinSwapExternAmountIn, and JoinSwapShareAmountOut
func (h AmmHooks) AfterJoinPool(ctx sdk.Context, sender sdk.AccAddress, pool ammtypes.Pool, enterCoins sdk.Coins, shareOutAmount math.Int) error  {
	return h.k.AfterJoinPool(ctx, sender, pool, enterCoins, shareOutAmount)
}

// AfterExitPool is called after ExitPool, ExitSwapShareAmountIn, and ExitSwapExternAmountOut
func (h AmmHooks) AfterExitPool(ctx sdk.Context, sender sdk.AccAddress, pool ammtypes.Pool, shareInAmount math.Int, exitCoins sdk.Coins) error {
	return h.k.AfterExitPool(ctx, sender, pool, shareInAmount, exitCoins)
}

// AfterSwap is called after SwapExactAmountIn and SwapExactAmountOut
func (h AmmHooks) AfterSwap(ctx sdk.Context, sender sdk.AccAddress, pool ammtypes.Pool, input sdk.Coins, output sdk.Coins) error {
	return h.k.AfterSwap(ctx, sender, pool, input, output)
}
