package types

import sdk "github.com/cosmos/cosmos-sdk/types"

type AmmHooks interface {
	// AfterPoolCreated is called after CreatePool
	AfterPoolCreated(ctx sdk.Context, sender sdk.AccAddress, pool Pool)

	// AfterJoinPool is called after JoinPool, JoinSwapExternAmountIn, and JoinSwapShareAmountOut
	AfterJoinPool(ctx sdk.Context, sender sdk.AccAddress, pool Pool, enterCoins sdk.Coins, shareOutAmount sdk.Int)

	// AfterExitPool is called after ExitPool, ExitSwapShareAmountIn, and ExitSwapExternAmountOut
	AfterExitPool(ctx sdk.Context, sender sdk.AccAddress, pool Pool, shareInAmount sdk.Int, exitCoins sdk.Coins) error

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

func (h MultiAmmHooks) AfterPoolCreated(ctx sdk.Context, sender sdk.AccAddress, pool Pool) {
	for i := range h {
		h[i].AfterPoolCreated(ctx, sender, pool)
	}
}

func (h MultiAmmHooks) AfterJoinPool(ctx sdk.Context, sender sdk.AccAddress, pool Pool, enterCoins sdk.Coins, shareOutAmount sdk.Int) {
	for i := range h {
		h[i].AfterJoinPool(ctx, sender, pool, enterCoins, shareOutAmount)
	}
}

func (h MultiAmmHooks) AfterExitPool(ctx sdk.Context, sender sdk.AccAddress, pool Pool, shareInAmount sdk.Int, exitCoins sdk.Coins) error {
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
