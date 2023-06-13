package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (p *Pool) IncreaseLiquidity(sharesOut sdk.Int, coinsIn sdk.Coins) {
	err := p.addToPoolAssetBalances(coinsIn)
	if err != nil {
		panic(err)
	}
	p.AddTotalShares(sharesOut)
}
