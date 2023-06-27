package types

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// CalcJoinPoolNoSwapShares calculates the number of shares created to execute an all-asset pool join with the provided amount of `tokensIn`.
// The input tokens must contain the same tokens as in the pool.
//
// Returns the number of shares created, the amount of coins actually joined into the pool, (in case of not being able to fully join),
// and the remaining tokens in `tokensIn` after joining. If an all-asset join is not possible, returns an error.
//
// Since CalcJoinPoolNoSwapShares is non-mutative, the steps for updating pool shares / liquidity are
// more complex / don't just alter the state.
// We should simplify this logic further in the future using multi-join equations.
func (p *Pool) CalcJoinPoolNoSwapShares(tokensIn sdk.Coins) (numShares sdk.Int, tokensJoined sdk.Coins, err error) {
	// get all 'pool assets' (aka current pool liquidity + balancer weight)
	poolAssetsByDenom, err := GetPoolAssetsByDenom(p.GetAllPoolAssets())
	if err != nil {
		return sdk.ZeroInt(), sdk.NewCoins(), err
	}

	err = EnsureDenomInPool(poolAssetsByDenom, tokensIn)
	if err != nil {
		return sdk.ZeroInt(), sdk.NewCoins(), err
	}

	// ensure that there aren't too many or too few assets in `tokensIn`
	if tokensIn.Len() != len(p.PoolAssets) {
		return sdk.ZeroInt(), sdk.NewCoins(), errors.New("no-swap joins require LP'ing with all assets in pool")
	}

	// execute a no-swap join with as many tokens as possible given a perfect ratio:
	// * numShares is how many shares are perfectly matched.
	// * remainingTokensIn is how many coins we have left to join that have not already been used.
	numShares, remainingTokensIn, err := MaximalExactRatioJoin(p, tokensIn)
	if err != nil {
		return sdk.ZeroInt(), sdk.NewCoins(), err
	}

	// ensure that no more tokens have been joined than is possible with the given `tokensIn`
	tokensJoined = tokensIn.Sub(remainingTokensIn...)
	if tokensJoined.IsAnyGT(tokensIn) {
		return sdk.ZeroInt(), sdk.NewCoins(), errors.New("an error has occurred, more coins joined than token In")
	}

	return numShares, tokensJoined, nil
}
