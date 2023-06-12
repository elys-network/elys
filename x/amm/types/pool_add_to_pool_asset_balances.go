package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (p *Pool) addToPoolAssetBalances(coins sdk.Coins) error {
	for _, coin := range coins {
		i, poolAsset, err := p.GetPoolAssetAndIndex(coin.Denom)
		if err != nil {
			return err
		}
		poolAsset.Token.Amount = poolAsset.Token.Amount.Add(coin.Amount)
		p.PoolAssets[i] = &poolAsset
	}
	return nil
}
