package keeper

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammkeeper "github.com/elys-network/elys/v6/x/amm/keeper"
	ammtypes "github.com/elys-network/elys/v6/x/amm/types"
	"github.com/elys-network/elys/v6/x/perpetual/types"
	"github.com/osmosis-labs/osmosis/osmomath"
)

func (k Keeper) GetAmmPool(ctx sdk.Context, poolId uint64) (ammtypes.Pool, error) {
	ammPool, found := k.amm.GetPool(ctx, poolId)
	if !found {
		return ammPool, errorsmod.Wrapf(types.ErrPoolDoesNotExist, "pool id %d", poolId)
	}
	return ammPool, nil
}

func getWeightBreakingFee(weightBalanceBonus osmomath.BigDec) osmomath.BigDec {
	// when weightBalanceBonus is 0, then breaking fee is also 0
	// when it's > 0, then breaking fee is still 0
	// when it's < 0, breaking fee is it's negative
	if weightBalanceBonus.IsNegative() {
		return weightBalanceBonus.Neg()
	} else {
		return osmomath.ZeroBigDec()
	}
}

// Swap estimation using amm CalcOutAmtGivenIn function
func (k Keeper) EstimateSwapGivenIn(ctx sdk.Context, tokenInAmount sdk.Coin, tokenOutDenom string, ammPool ammtypes.Pool, owner string, emitFeesEvent bool) (math.Int, osmomath.BigDec, osmomath.BigDec, math.LegacyDec, math.LegacyDec, error) {
	if tokenInAmount.IsZero() {
		return math.Int{}, osmomath.BigDec{}, osmomath.BigDec{}, math.LegacyDec{}, math.LegacyDec{}, fmt.Errorf("tokenInAmount is zero for EstimateSwapGivenIn")
	}
	params := k.GetParams(ctx)

	addr, err := sdk.AccAddressFromBech32(owner)
	if err != nil {
		addr = sdk.AccAddress{}
	}
	_, tier := k.tierKeeper.GetMembershipTier(ctx, addr)
	perpetualFees := ammtypes.ApplyDiscount(params.GetBigDecPerpetualSwapFee(), tier.GetBigDecDiscount())
	takersFee := k.parameterKeeper.GetParams(ctx).GetBigDecTakerFees()
	// Estimate swap
	snapshot := k.amm.GetPoolWithAccountedBalance(ctx, ammPool.PoolId)
	tokensIn := sdk.Coins{tokenInAmount}
	tokenOut, slippage, slippageAmount, weightBalanceBonus, _, _, err := k.amm.SwapOutAmtGivenIn(ctx, ammPool.PoolId, k.oracleKeeper, snapshot, tokensIn, tokenOutDenom, perpetualFees, params.GetBigDecWeightBreakingFeeFactor(), takersFee)
	if err != nil {
		return math.ZeroInt(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), math.LegacyDec{}, math.LegacyDec{}, errorsmod.Wrapf(err, "unable to swap (EstimateSwapGivenIn) for in %s and out denom %s", tokenInAmount.String(), tokenOutDenom)
	}

	if tokenOut.IsZero() {
		return math.ZeroInt(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), math.LegacyDec{}, math.LegacyDec{}, errorsmod.Wrapf(types.ErrAmountTooLow, "tokenOut is zero for swap (EstimateSwapGivenIn) for in %s and out denom %s", tokenInAmount.String(), tokenOutDenom)
	}

	weightBreakingFee := getWeightBreakingFee(weightBalanceBonus)

	if emitFeesEvent {
		perpFeesValueInUSD := osmomath.ZeroBigDec()
		if perpetualFees.IsPositive() {
			perpetualFeesCoins := ammkeeper.PortionCoins(tokensIn, perpetualFees)
			perpFeesValueInUSD = k.amm.CalculateCoinsUSDValue(ctx, perpetualFeesCoins)
		}

		takerFeesAmountInUSD := osmomath.ZeroBigDec()
		if takersFee.IsPositive() {
			takerFeesInCoins := ammkeeper.PortionCoins(tokensIn, takersFee)
			takerFeesAmountInUSD = k.amm.CalculateCoinsUSDValue(ctx, takerFeesInCoins)
		}

		slippageAmountInUSD := k.amm.CalculateUSDValue(ctx, tokenOut.Denom, slippageAmount.Dec().TruncateInt())
		weightBreakingFeesAmountInUSD := osmomath.ZeroBigDec()
		if !weightBreakingFee.IsZero() {
			weightBreakingFeeAmount := osmomath.BigDecFromSDKInt(tokenInAmount.Amount).Mul(weightBreakingFee).Dec().RoundInt()
			weightBreakingFeesAmountInUSD = k.amm.CalculateUSDValue(ctx, tokenInAmount.Denom, weightBreakingFeeAmount)
		}

		if !(perpFeesValueInUSD.IsZero() && slippageAmountInUSD.IsZero() && weightBreakingFeesAmountInUSD.IsZero() && takerFeesAmountInUSD.IsZero()) {
			types.EmitPerpetualFeesEvent(ctx, perpFeesValueInUSD.String(), slippageAmountInUSD.String(), weightBreakingFeesAmountInUSD.String(), takerFeesAmountInUSD.String())
		}
	}

	return tokenOut.Amount, slippage, weightBreakingFee, perpetualFees.Dec(), takersFee.Dec(), nil
}

// Swap estimation using amm CalcInAmtGivenOut function
func (k Keeper) EstimateSwapGivenOut(ctx sdk.Context, tokenOutAmount sdk.Coin, tokenInDenom string, ammPool ammtypes.Pool, owner string, emitFeesEvent bool) (math.Int, osmomath.BigDec, osmomath.BigDec, math.LegacyDec, math.LegacyDec, error) {
	if tokenOutAmount.IsZero() {
		return math.Int{}, osmomath.BigDec{}, osmomath.BigDec{}, math.LegacyDec{}, math.LegacyDec{}, fmt.Errorf("tokenOutAmount is zero for EstimateSwapGivenOut")
	}
	params := k.GetParams(ctx)
	tokensOut := sdk.Coins{tokenOutAmount}

	addr, err := sdk.AccAddressFromBech32(owner)
	if err != nil {
		addr = sdk.AccAddress{}
	}
	_, tier := k.tierKeeper.GetMembershipTier(ctx, addr)
	perpetualFees := ammtypes.ApplyDiscount(params.GetBigDecPerpetualSwapFee(), tier.GetBigDecDiscount())
	takersFee := k.parameterKeeper.GetParams(ctx).GetBigDecTakerFees()

	// Estimate swap
	snapshot := k.amm.GetPoolWithAccountedBalance(ctx, ammPool.PoolId)
	tokenIn, slippage, slippageAmount, weightBalanceBonus, oracleIn, _, err := k.amm.SwapInAmtGivenOut(ctx, ammPool.PoolId, k.oracleKeeper, snapshot, tokensOut, tokenInDenom, perpetualFees, params.GetBigDecWeightBreakingFeeFactor(), takersFee)
	if err != nil {
		return math.ZeroInt(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), math.LegacyDec{}, math.LegacyDec{}, errorsmod.Wrapf(err, "unable to swap (EstimateSwapGivenOut) for out %s and in denom %s", tokenOutAmount.String(), tokenInDenom)
	}

	if tokenIn.IsZero() {
		return math.ZeroInt(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), math.LegacyDec{}, math.LegacyDec{}, errorsmod.Wrapf(types.ErrAmountTooLow, "tokenIn is zero for swap (EstimateSwapGivenOut) for out %s and in denom %s", tokenOutAmount.String(), tokenInDenom)
	}

	weightBreakingFee := getWeightBreakingFee(weightBalanceBonus)

	if emitFeesEvent {
		takeFeesFrom := sdk.NewCoins(sdk.NewCoin(tokenInDenom, oracleIn.Dec().TruncateInt()))
		if !ammPool.PoolParams.UseOracle {
			takeFeesFrom = sdk.Coins{tokenIn}
		}
		perpFeesValueInUSD := osmomath.ZeroBigDec()
		if perpetualFees.IsPositive() {
			perpetualFeesCoins := ammkeeper.PortionCoins(takeFeesFrom, perpetualFees)
			perpFeesValueInUSD = k.amm.CalculateCoinsUSDValue(ctx, perpetualFeesCoins)
		}

		takerFeesAmountInUSD := osmomath.ZeroBigDec()
		if takersFee.IsPositive() {
			takerFeesInCoins := ammkeeper.PortionCoins(takeFeesFrom, takersFee)
			takerFeesAmountInUSD = k.amm.CalculateCoinsUSDValue(ctx, takerFeesInCoins)
		}

		slippageAmountInUSD := k.amm.CalculateUSDValue(ctx, tokenInDenom, slippageAmount.Dec().TruncateInt())
		weightBreakingFeesAmountInUSD := osmomath.ZeroBigDec()
		if !weightBreakingFee.IsZero() {
			weightBreakingFeeAmount := oracleIn.Mul(weightBreakingFee).Dec().RoundInt()
			weightBreakingFeesAmountInUSD = k.amm.CalculateUSDValue(ctx, tokenInDenom, weightBreakingFeeAmount)
		}

		if !(perpFeesValueInUSD.IsZero() && slippageAmountInUSD.IsZero() && weightBreakingFeesAmountInUSD.IsZero() && takerFeesAmountInUSD.IsZero()) {
			types.EmitPerpetualFeesEvent(ctx, perpFeesValueInUSD.String(), slippageAmountInUSD.String(), weightBreakingFeesAmountInUSD.String(), takerFeesAmountInUSD.String())
		}
	}
	return tokenIn.Amount, slippage, getWeightBreakingFee(weightBalanceBonus), perpetualFees.Dec(), takersFee.Dec(), nil
}
