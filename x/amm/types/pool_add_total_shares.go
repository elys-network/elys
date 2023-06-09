package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (p *Pool) AddTotalShares(amt sdk.Int) {
	p.TotalShares.Amount = p.TotalShares.Amount.Add(amt)
}
