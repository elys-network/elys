package types

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// JoinPoolNoSwap calculates the number of shares needed for an all-asset join given tokensIn with swapFee applied.
// It updates the liquidity if the pool is joined successfully. If not, returns error.
func (p *Pool) JoinPoolNoSwap(tokensIn sdk.Coins, swapFee sdk.Dec) (numShares math.Int, err error) {
	numShares, tokensJoined, err := p.CalcJoinPoolNoSwapShares(tokensIn, swapFee)
	if err != nil {
		return math.Int{}, err
	}

	// update pool with the calculated share and liquidity needed to join pool
	p.IncreaseLiquidity(numShares, tokensJoined)
	return numShares, nil
}
