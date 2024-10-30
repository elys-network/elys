package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

// Swap estimation using amm CalcOutAmtGivenIn function
func (k Keeper) EstimateSwapGivenIn(ctx sdk.Context, tokenInAmount sdk.Coin, tokenOutDenom string, ammPool ammtypes.Pool) (math.Int, math.LegacyDec, error) {
	swapFee := k.GetPerpetualSwapFee(ctx)
	// Estimate swap
	snapshot := k.amm.GetPoolSnapshotOrSet(ctx, ammPool)
	tokensIn := sdk.Coins{tokenInAmount}
	tokenOut, slippage, _, _, err := k.amm.SwapOutAmtGivenIn(ctx, ammPool.PoolId, k.oracleKeeper, &snapshot, tokensIn, tokenOutDenom, swapFee)
	if err != nil {
		return sdkmath.ZeroInt(), math.LegacyZeroDec(), err
	}

	if tokenOut.IsZero() {
		return sdkmath.ZeroInt(), math.LegacyZeroDec(), types.ErrAmountTooLow
	}
	return tokenOut.Amount, slippage, nil
}

// Swap estimation using amm CalcInAmtGivenOut function
func (k Keeper) EstimateSwapGivenOut(ctx sdk.Context, tokenOutAmount sdk.Coin, tokenInDenom string, ammPool ammtypes.Pool) (math.Int, math.LegacyDec, error) {
	perpetualSwapFee := k.GetPerpetualSwapFee(ctx)
	tokensOut := sdk.Coins{tokenOutAmount}
	// Estimate swap
	snapshot := k.amm.GetPoolSnapshotOrSet(ctx, ammPool)
	tokenIn, slippage, _, _, err := k.amm.SwapInAmtGivenOut(ctx, ammPool.PoolId, k.oracleKeeper, &snapshot, tokensOut, tokenInDenom, perpetualSwapFee)
	if err != nil {
		return sdk.ZeroInt(), math.LegacyZeroDec(), err
	}

	if tokenIn.IsZero() {
		return sdk.ZeroInt(), math.LegacyZeroDec(), types.ErrAmountTooLow
	}
	return tokenIn.Amount, slippage, nil
}
