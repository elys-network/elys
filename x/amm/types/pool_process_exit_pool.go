package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// exitPool exits the pool given exitingCoins and exitingShares.
// updates the pool's liquidity and totalShares.
func (p *Pool) processExitPool(ctx sdk.Context, exitingCoins sdk.Coins, exitingShares sdk.Int) error {
	balances := p.GetTotalPoolLiquidity().Sub(exitingCoins...)
	if err := p.UpdatePoolAssetBalances(balances); err != nil {
		return err
	}

	totalShares := p.GetTotalShares().Amount
	p.TotalShares = sdk.NewCoin(p.TotalShares.Denom, totalShares.Sub(exitingShares))

	return nil
}
