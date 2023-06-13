package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (p *Pool) ExitPool(ctx sdk.Context, exitingShares sdk.Int, exitFee sdk.Dec) (exitingCoins sdk.Coins, err error) {
	exitingCoins, err = p.CalcExitPoolCoinsFromShares(ctx, exitingShares, exitFee)
	if err != nil {
		return sdk.Coins{}, err
	}

	if err := p.processExitPool(ctx, exitingCoins, exitingShares); err != nil {
		return sdk.Coins{}, err
	}

	return exitingCoins, nil
}
