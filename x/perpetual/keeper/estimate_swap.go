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
	// when weightBalanceBonus is 0, then breaking fee is also 0
	// when it's > 0, then breaking fee is still 0
	// when it's < 0, breaking fee is it's negative
	if weightBalanceBonus.IsNegative() {
		return weightBalanceBonus.Neg()
	} else {
		return math.LegacyZeroDec()
	}
}

// Swap estimation using amm CalcOutAmtGivenIn function
func (k Keeper) EstimateSwapGivenIn(ctx sdk.Context, tokenInAmount sdk.Coin, tokenOutDenom string, ammPool ammtypes.Pool, owner string) (math.Int, math.LegacyDec, math.LegacyDec, error) {
	if tokenInAmount.IsZero() {
		return math.Int{}, math.LegacyDec{}, math.LegacyDec{}, fmt.Errorf("tokenInAmount is zero for EstimateSwapGivenIn")
	}
	params := k.GetParams(ctx)

	addr, err := sdk.AccAddressFromBech32(owner)
	if err != nil {
		addr = sdk.AccAddress{}
	}
	_, tier := k.tierKeeper.GetMembershipTier(ctx, addr)
	perpetualFees := ammtypes.ApplyDiscount(params.PerpetualSwapFee, tier.Discount)
	// Estimate swap
	snapshot := k.amm.GetAccountedPoolSnapshotOrSet(ctx, ammPool)
	tokensIn := sdk.Coins{tokenInAmount}
	tokenOut, slippage, _, weightBalanceBonus, _, _, err := k.amm.SwapOutAmtGivenIn(ctx, ammPool.PoolId, k.oracleKeeper, &snapshot, tokensIn, tokenOutDenom, perpetualFees, params.WeightBreakingFeeFactor)
	if err != nil {
		return math.ZeroInt(), math.LegacyZeroDec(), math.LegacyZeroDec(), errorsmod.Wrapf(err, "unable to swap (EstimateSwapGivenIn) for in %s and out denom %s", tokenInAmount.String(), tokenOutDenom)
	}

	if tokenOut.IsZero() {
		return math.ZeroInt(), math.LegacyZeroDec(), math.LegacyZeroDec(), errorsmod.Wrapf(types.ErrAmountTooLow, "tokenOut is zero for swap (EstimateSwapGivenIn) for in %s and out denom %s", tokenInAmount.String(), tokenOutDenom)
	}
	return tokenOut.Amount, slippage, getWeightBreakingFee(weightBalanceBonus), nil
}

// Swap estimation using amm CalcInAmtGivenOut function
func (k Keeper) EstimateSwapGivenOut(ctx sdk.Context, tokenOutAmount sdk.Coin, tokenInDenom string, ammPool ammtypes.Pool, owner string) (math.Int, math.LegacyDec, math.LegacyDec, error) {
	if tokenOutAmount.IsZero() {
		return math.Int{}, math.LegacyDec{}, math.LegacyDec{}, fmt.Errorf("tokenOutAmount is zero for EstimateSwapGivenOut")
	}
	params := k.GetParams(ctx)
	tokensOut := sdk.Coins{tokenOutAmount}

	addr, err := sdk.AccAddressFromBech32(owner)
	if err != nil {
		addr = sdk.AccAddress{}
	}
	_, tier := k.tierKeeper.GetMembershipTier(ctx, addr)
	perpetualFees := ammtypes.ApplyDiscount(params.PerpetualSwapFee, tier.Discount)

	// Estimate swap
	snapshot := k.amm.GetAccountedPoolSnapshotOrSet(ctx, ammPool)
	tokenIn, slippage, _, weightBalanceBonus, _, _, err := k.amm.SwapInAmtGivenOut(ctx, ammPool.PoolId, k.oracleKeeper, &snapshot, tokensOut, tokenInDenom, perpetualFees, params.WeightBreakingFeeFactor)
	if err != nil {
		return math.ZeroInt(), math.LegacyZeroDec(), math.LegacyZeroDec(), errorsmod.Wrapf(err, "unable to swap (EstimateSwapGivenOut) for out %s and in denom %s", tokenOutAmount.String(), tokenInDenom)
	}

	if tokenIn.IsZero() {
		return math.ZeroInt(), math.LegacyZeroDec(), math.LegacyZeroDec(), errorsmod.Wrapf(types.ErrAmountTooLow, "tokenIn is zero for swap (EstimateSwapGivenOut) for out %s and in denom %s", tokenOutAmount.String(), tokenInDenom)
	}
	return tokenIn.Amount, slippage, getWeightBreakingFee(weightBalanceBonus), nil
}
