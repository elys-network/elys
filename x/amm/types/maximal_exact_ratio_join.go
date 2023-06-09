package types

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MaximalExactRatioJoin calculates the maximal amount of tokens that can be joined whilst maintaining pool asset's ratio
// returning the number of shares that'd be and how many coins would be left over.
//
//	e.g) suppose we have a pool of 10 foo tokens and 10 bar tokens, with the total amount of 100 shares.
//		 if `tokensIn` provided is 1 foo token and 2 bar tokens, `MaximalExactRatioJoin`
//		 would be returning (10 shares, 1 bar token, nil)
//
// This can be used when `tokensIn` are not guaranteed the same ratio as assets in the pool.
// Calculation for this is done in the following steps.
//  1. iterate through all the tokens provided as an argument, calculate how much ratio it accounts for the asset in the pool
//  2. get the minimal share ratio that would work as the benchmark for all tokens.
//  3. calculate the number of shares that could be joined (total share * min share ratio), return the remaining coins
func MaximalExactRatioJoin(p *Pool, tokensIn sdk.Coins) (numShares sdk.Int, remCoins sdk.Coins, err error) {
	coinShareRatios := make([]sdk.Dec, len(tokensIn))
	minShareRatio := sdk.MaxSortableDec
	maxShareRatio := sdk.ZeroDec()

	poolLiquidity := p.GetTotalPoolLiquidity()
	totalShares := p.GetTotalShares()

	for i, coin := range tokensIn {
		// Note: QuoInt implements floor division, unlike Quo
		// This is because it calls the native golang routine big.Int.Quo
		// https://pkg.go.dev/math/big#Int.Quo
		shareRatio := sdk.NewDecFromBigInt(coin.Amount.BigInt()).QuoInt(poolLiquidity.AmountOfNoDenomValidation(coin.Denom))
		if shareRatio.LT(minShareRatio) {
			minShareRatio = shareRatio
		}
		if shareRatio.GT(maxShareRatio) {
			maxShareRatio = shareRatio
		}
		coinShareRatios[i] = shareRatio
	}

	if minShareRatio.Equal(sdk.MaxSortableDec) {
		return numShares, remCoins, errors.New("unexpected error in MaximalExactRatioJoin")
	}

	remCoins = sdk.Coins{}
	// critically we round down here (TruncateInt), to ensure that the returned LP shares
	// are always less than or equal to % liquidity added.
	numShares = minShareRatio.MulInt(totalShares.Amount).TruncateInt()

	// if we have multiple share values, calculate remainingCoins
	if !minShareRatio.Equal(maxShareRatio) {
		// we have to calculate remCoins
		for i, coin := range tokensIn {
			// if coinShareRatios[i] == minShareRatio, no remainder
			if coinShareRatios[i].Equal(minShareRatio) {
				continue
			}

			usedAmount := minShareRatio.MulInt(poolLiquidity.AmountOfNoDenomValidation(coin.Denom)).Ceil().TruncateInt()
			newAmt := coin.Amount.Sub(usedAmount)
			// if newAmt is non-zero, add to RemCoins. (It could be zero due to rounding)
			if !newAmt.IsZero() {
				remCoins = remCoins.Add(sdk.Coin{Denom: coin.Denom, Amount: newAmt})
			}
		}
	}

	return numShares, remCoins, nil
}
