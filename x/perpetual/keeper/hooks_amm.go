package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/v6/x/amm/types"
)

// EpochHooks wrapper struct for tvl keeper
type AmmHooks struct {
	k Keeper
}

var _ ammtypes.AmmHooks = AmmHooks{}

// Return the wrapper struct
func (k Keeper) AmmHooks() AmmHooks {
	return AmmHooks{k}
}

// AfterPoolCreated is called after CreatePool
// We are already creating accounted pool using amm hooks in accounted pool module, so no need to create it here
// ideally we should create accounted pool after perpetual pool is needed but then that would follow a complicated process as perpetual module isn't aware of when amm pool is created directly
// This method also allows if any other module in future requires accounted pool, it doesn't need to do create any new accounted pool.
func (h AmmHooks) AfterPoolCreated(ctx sdk.Context, sender sdk.AccAddress, pool ammtypes.Pool) error {
	return nil
}

// AfterJoinPool is called after JoinPool, JoinSwapExternAmountIn, and JoinSwapShareAmountOut
func (h AmmHooks) AfterJoinPool(ctx sdk.Context, sender sdk.AccAddress, ammPool ammtypes.Pool, enterCoins sdk.Coins, shareOutAmount math.Int) error {
	perpetualPool, found := h.k.GetPool(ctx, ammPool.PoolId)
	if !found {
		// It is possible that this pool haven't been enabled
		return nil
	}

	err := h.k.UpdatePoolHealth(ctx, &perpetualPool)
	if err != nil {
		return err
	}

	err = h.k.CheckLowPoolHealthAndMinimumCustody(ctx, ammPool.PoolId)
	if err != nil {
		return err
	}

	return nil

}

// AfterExitPool is called after ExitPool, ExitSwapShareAmountIn, and ExitSwapExternAmountOut
func (h AmmHooks) AfterExitPool(ctx sdk.Context, sender sdk.AccAddress, ammPool ammtypes.Pool, shareInAmount math.Int, exitCoins sdk.Coins) error {
	perpetualPool, found := h.k.GetPool(ctx, ammPool.PoolId)
	if !found {
		// It is possible that this pool haven't been enabled
		return nil
	}

	err := h.k.UpdatePoolHealth(ctx, &perpetualPool)
	if err != nil {
		return err
	}

	err = h.k.CheckLowPoolHealthAndMinimumCustody(ctx, ammPool.PoolId)
	if err != nil {
		return err
	}

	return nil
}

// AfterSwap is called after SwapExactAmountIn and SwapExactAmountOut
func (h AmmHooks) AfterSwap(ctx sdk.Context, sender sdk.AccAddress, ammPool ammtypes.Pool, input sdk.Coins, output sdk.Coins) error {
	perpetualPool, found := h.k.GetPool(ctx, ammPool.PoolId)
	if !found {
		// It is possible that this pool haven't been enabled
		return nil
	}

	err := h.k.UpdatePoolHealth(ctx, &perpetualPool)
	if err != nil {
		return err
	}

	err = h.k.CheckLowPoolHealthAndMinimumCustody(ctx, ammPool.PoolId)
	if err != nil {
		return err
	}

	return nil
}
