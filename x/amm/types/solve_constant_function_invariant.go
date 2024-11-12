package types

import (
	sdkmath "cosmossdk.io/math"
)

// solveConstantFunctionInvariant solves the constant function of an AMM
// that determines the relationship between the differences of two sides
// of assets inside the pool.
// For fixed balanceXBefore, balanceXAfter, weightX, balanceY, weightY,
// we could deduce the balanceYDelta, calculated by:
// balanceYDelta = balanceY * (1 - (balanceXBefore/balanceXAfter)^(weightX/weightY))
// balanceYDelta is positive when the balance liquidity decreases.
// balanceYDelta is negative when the balance liquidity increases.
//
// panics if tokenWeightUnknown is 0.
func solveConstantFunctionInvariant(
	tokenBalanceFixedBefore,
	tokenBalanceFixedAfter,
	tokenWeightFixed,
	tokenBalanceUnknownBefore,
	tokenWeightUnknown sdkmath.LegacyDec,
) (sdkmath.LegacyDec, error) {
	// Ensure tokenWeightUnknown is not zero to avoid division by zero
	if tokenWeightUnknown.IsZero() {
		return sdkmath.LegacyZeroDec(), ErrAmountTooLow
	}

	// weightRatio = (weightX/weightY)
	weightRatio := tokenWeightFixed.Quo(tokenWeightUnknown)

	// Ensure tokenBalanceFixedAfter is not zero to avoid division by zero
	if tokenBalanceFixedAfter.IsZero() || tokenBalanceFixedAfter.IsNegative() {
		return sdkmath.LegacyZeroDec(), ErrAmountTooLow
	}

	// y = balanceXBefore/balanceXAfter
	y := tokenBalanceFixedBefore.Quo(tokenBalanceFixedAfter)

	// amountY = balanceY * (1 - (y ^ weightRatio))
	yToWeightRatio := Pow(y, weightRatio)
	paranthetical := sdkmath.LegacyOneDec().Sub(yToWeightRatio)
	amountY := tokenBalanceUnknownBefore.Mul(paranthetical)
	return amountY, nil
}

// E.g. tokenA: ELYS, tokenB: USDC
func CalculateTokenARate(tokenBalanceA, tokenWeightA, tokenBalanceB, tokenWeightB sdkmath.LegacyDec) sdkmath.LegacyDec {
	if tokenBalanceA.IsZero() || tokenWeightB.IsZero() {
		return sdkmath.LegacyZeroDec()
	}
	return tokenBalanceB.
		Mul(tokenWeightA).
		Quo(tokenWeightB).
		Quo(tokenBalanceA)
}

// feeRatio returns the fee ratio that is defined as follows:
// 1 - ((1 - normalizedTokenWeightOut) * spreadFactor)
func feeRatio(normalizedWeight, spreadFactor sdkmath.LegacyDec) sdkmath.LegacyDec {
	return sdkmath.LegacyOneDec().Sub((sdkmath.LegacyOneDec().Sub(normalizedWeight)).Mul(spreadFactor))
}

// balancer notation: pAo - pool shares amount out, given single asset in
// the second argument requires the tokenWeightIn / total token weight.
func calcPoolSharesOutGivenSingleAssetIn(
	tokenBalanceIn,
	normalizedTokenWeightIn,
	poolShares,
	tokenAmountIn,
	spreadFactor sdkmath.LegacyDec,
) (sdkmath.LegacyDec, error) {
	// deduct spread factor on the in asset.
	// We don't charge spread factor on the token amount that we imagine as unswapped (the normalized weight).
	// So effective_swapfee = spread factor * (1 - normalized_token_weight)
	tokenAmountInAfterFee := tokenAmountIn.Mul(feeRatio(normalizedTokenWeightIn, spreadFactor))
	// To figure out the number of shares we add, first notice that in balancer we can treat
	// the number of shares as linearly related to the `k` value function. This is due to the normalization.
	// e.g.
	// if x^.5 y^.5 = k, then we `n` x the liquidity to `(nx)^.5 (ny)^.5 = nk = k'`
	// We generalize this linear relation to do the liquidity add for the not-all-asset case.
	// Suppose we increase the supply of x by x', so we want to solve for `k'/k`.
	// This is `(x + x')^{weight} * old_terms / (x^{weight} * old_terms) = (x + x')^{weight} / (x^{weight})`
	// The number of new shares we need to make is then `old_shares * ((k'/k) - 1)`
	// Whats very cool, is that this turns out to be the exact same `solveConstantFunctionInvariant` code
	// with the answer's sign reversed.
	poolAmountOut, err := solveConstantFunctionInvariant(
		tokenBalanceIn.Add(tokenAmountInAfterFee),
		tokenBalanceIn,
		normalizedTokenWeightIn,
		poolShares,
		sdkmath.LegacyOneDec())
	if err != nil {
		return sdkmath.LegacyDec{}, err
	}
	poolAmountOut = poolAmountOut.Neg()
	return poolAmountOut, nil
}
