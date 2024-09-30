package types

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type AmmHooks interface {
	// AfterPoolCreated is called after CreatePool
	AfterPoolCreated(ctx sdk.Context, sender sdk.AccAddress, pool Pool) error

	// AfterJoinPool is called after JoinPool, JoinSwapExternAmountIn, and JoinSwapShareAmountOut
	AfterJoinPool(ctx sdk.Context, sender sdk.AccAddress, pool Pool, enterCoins sdk.Coins, shareOutAmount sdkmath.Int) error

	// AfterExitPool is called after ExitPool, ExitSwapShareAmountIn, and ExitSwapExternAmountOut
	AfterExitPool(ctx sdk.Context, sender sdk.AccAddress, pool Pool, shareInAmount sdkmath.Int, exitCoins sdk.Coins) error

	// AfterSwap is called after SwapExactAmountIn and SwapExactAmountOut
	AfterSwap(ctx sdk.Context, sender sdk.AccAddress, pool Pool, input sdk.Coins, output sdk.Coins) error
}

var _ AmmHooks = MultiAmmHooks{}

// combine multiple amm hooks, all hook functions are run in array sequence.
type MultiAmmHooks []AmmHooks

// Creates hooks for the Amm Module.
func NewMultiAmmHooks(hooks ...AmmHooks) MultiAmmHooks {
	return hooks
}

func (h MultiAmmHooks) AfterPoolCreated(ctx sdk.Context, sender sdk.AccAddress, pool Pool) error {
	for i := range h {
		err := h[i].AfterPoolCreated(ctx, sender, pool)
		if err != nil {
			return err
		}
	}
	return nil
}

func (h MultiAmmHooks) AfterJoinPool(ctx sdk.Context, sender sdk.AccAddress, pool Pool, enterCoins sdk.Coins, shareOutAmount sdkmath.Int) error {
	for i := range h {
		err := h[i].AfterJoinPool(ctx, sender, pool, enterCoins, shareOutAmount)
		if err != nil {
			return err
		}
	}
	return nil
}

func (h MultiAmmHooks) AfterExitPool(ctx sdk.Context, sender sdk.AccAddress, pool Pool, shareInAmount sdkmath.Int, exitCoins sdk.Coins) error {
	for i := range h {
		err := h[i].AfterExitPool(ctx, sender, pool, shareInAmount, exitCoins)
		if err != nil {
			return err
		}
	}
	return nil
}

func (h MultiAmmHooks) AfterSwap(ctx sdk.Context, sender sdk.AccAddress, pool Pool, input sdk.Coins, output sdk.Coins) error {
	for i := range h {
		err := h[i].AfterSwap(ctx, sender, pool, input, output)
		if err != nil {
			return err
		}
	}
	return nil
}
