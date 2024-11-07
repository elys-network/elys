package keeper

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func getWeightBreakingFee(weightBalanceBonus math.LegacyDec) math.LegacyDec {
	if weightBalanceBonus.IsNegative() {
		return weightBalanceBonus.Neg()
	} else {
		return math.LegacyZeroDec()
	}
}

// Swap estimation using amm CalcOutAmtGivenIn function
func (k Keeper) EstimateSwapGivenIn(ctx sdk.Context, tokenInAmount sdk.Coin, tokenOutDenom string, ammPool ammtypes.Pool) (math.Int, math.LegacyDec, math.LegacyDec, error) {
	if tokenInAmount.IsZero() {
		return math.Int{}, math.LegacyDec{}, math.LegacyDec{}, fmt.Errorf("tokenInAmount is zero for EstimateSwapGivenIn")
	}
	params := k.GetParams(ctx)
	// Estimate swap
	snapshot := k.amm.GetAccountedPoolSnapshotOrSet(ctx, ammPool)
	tokensIn := sdk.Coins{tokenInAmount}
	tokenOut, slippage, _, weightBalanceBonus, err := k.amm.SwapOutAmtGivenIn(ctx, ammPool.PoolId, k.oracleKeeper, &snapshot, tokensIn, tokenOutDenom, params.PerpetualSwapFee, params.WeightBreakingFeeFactor)
	if err != nil {
		return sdk.ZeroInt(), math.LegacyZeroDec(), math.LegacyZeroDec(), errorsmod.Wrapf(err, "unable to swap (EstimateSwapGivenIn) for in %s and out denom %s", tokenInAmount.String(), tokenOutDenom)
	}

	if tokenOut.IsZero() {
		return sdk.ZeroInt(), math.LegacyZeroDec(), math.LegacyZeroDec(), errorsmod.Wrapf(types.ErrAmountTooLow, "tokenOut is zero for swap (EstimateSwapGivenIn) for in %s and out denom %s", tokenInAmount.String(), tokenOutDenom)
	}
	return tokenOut.Amount, slippage, getWeightBreakingFee(weightBalanceBonus), nil
}

// Swap estimation using amm CalcInAmtGivenOut function
func (k Keeper) EstimateSwapGivenOut(ctx sdk.Context, tokenOutAmount sdk.Coin, tokenInDenom string, ammPool ammtypes.Pool) (math.Int, math.LegacyDec, math.LegacyDec, error) {
	if tokenOutAmount.IsZero() {
		return math.Int{}, math.LegacyDec{}, math.LegacyDec{}, fmt.Errorf("tokenOutAmount is zero for EstimateSwapGivenOut")
	}
	params := k.GetParams(ctx)
	tokensOut := sdk.Coins{tokenOutAmount}
	// Estimate swap
	snapshot := k.amm.GetAccountedPoolSnapshotOrSet(ctx, ammPool)
	tokenIn, slippage, _, weightBalanceBonus, err := k.amm.SwapInAmtGivenOut(ctx, ammPool.PoolId, k.oracleKeeper, &snapshot, tokensOut, tokenInDenom, params.PerpetualSwapFee, params.WeightBreakingFeeFactor)
	if err != nil {
		return sdk.ZeroInt(), math.LegacyZeroDec(), math.LegacyZeroDec(), errorsmod.Wrapf(err, "unable to swap (EstimateSwapGivenOut) for out %s and in denom %s", tokenOutAmount.String(), tokenInDenom)
	}

	if tokenIn.IsZero() {
		return sdk.ZeroInt(), math.LegacyZeroDec(), math.LegacyZeroDec(), errorsmod.Wrapf(types.ErrAmountTooLow, "tokenIn is zero for swap (EstimateSwapGivenOut) for out %s and in denom %s", tokenOutAmount.String(), tokenInDenom)
	}
	return tokenIn.Amount, slippage, getWeightBreakingFee(weightBalanceBonus), nil
}
