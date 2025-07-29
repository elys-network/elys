package keeper

import (
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/v7/x/amm/types"
)

func (k Keeper) CheckAmmPoolBalance(ctx sdk.Context, ammPool ammtypes.Pool) error {
	leveragePool, found := k.GetPool(ctx, ammPool.PoolId)
	if !found {
		// It is possible that this pool haven't been enabled
		return nil
	}
	stablestakeAmmPool := k.stableKeeper.GetAmmPool(ctx, ammPool.PoolId)
	params := k.GetParams(ctx)

	for _, asset := range ammPool.PoolAssets {
		for _, liabilties := range stablestakeAmmPool.TotalLiabilities {
			reducedLiabilities := params.LiabilitiesFactor.Mul(math.LegacyNewDecFromInt(liabilties.Amount))
			if asset.Token.Denom == liabilties.Denom && asset.Token.Amount.LT(reducedLiabilities.TruncateInt()) {
				return fmt.Errorf("insufficient amount of %s after the operation for leveragelp", asset.Token.Denom)
			}
		}
	}

	ratio := leveragePool.LeveragedLpAmount.ToLegacyDec().Quo(ammPool.TotalShares.Amount.ToLegacyDec())

	maxRatio := leveragePool.MaxLeveragelpRatio.Add(params.ExitBuffer)
	if ratio.GT(maxRatio) {
		return fmt.Errorf("operation not allowed: pool leverage position becomes %s (> %s)", ratio.String(), maxRatio.String())
	}

	return nil
}

// AfterPoolCreated is called after CreatePool
func (k Keeper) AfterPoolCreated(ctx sdk.Context, sender sdk.AccAddress, ammPool ammtypes.Pool) error {
	return nil
}

// AfterJoinPool is called after JoinPool, JoinSwapExternAmountIn, and JoinSwapShareAmountOut
func (k Keeper) AfterJoinPool(ctx sdk.Context, sender sdk.AccAddress, ammPool ammtypes.Pool, enterCoins sdk.Coins, shareOutAmount math.Int) error {
	return nil
}

// AfterExitPool is called after ExitPool, ExitSwapShareAmountIn, and ExitSwapExternAmountOut
func (k Keeper) AfterExitPool(ctx sdk.Context, sender sdk.AccAddress, ammPool ammtypes.Pool, shareInAmount math.Int, exitCoins sdk.Coins) error {
	return k.CheckAmmPoolBalance(ctx, ammPool)
}

// AfterSwap is called after SwapExactAmountIn and SwapExactAmountOut
func (k Keeper) AfterSwap(ctx sdk.Context, sender sdk.AccAddress, ammPool ammtypes.Pool, input sdk.Coins, output sdk.Coins) error {
	return k.CheckAmmPoolBalance(ctx, ammPool)
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
func (h AmmHooks) AfterJoinPool(ctx sdk.Context, sender sdk.AccAddress, pool ammtypes.Pool, enterCoins sdk.Coins, shareOutAmount math.Int) error {
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
