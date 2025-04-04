package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

func PoolAddress() sdk.AccAddress {
	return authtypes.NewModuleAddress(ModuleName)
}

func (p *AmmPool) AddLiabilities(coin sdk.Coin) {
	p.TotalLiabilities = p.TotalLiabilities.Add(coin)
}

func (p *AmmPool) SubLiabilities(coin sdk.Coin) {
	p.TotalLiabilities = p.TotalLiabilities.Sub(coin)
}
