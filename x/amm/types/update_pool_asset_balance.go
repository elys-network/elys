package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (p *Pool) UpdatePoolAssetBalance(coin sdk.Coin) error {
	// Check that PoolAsset exists.
	assetIndex, existingAsset, err := p.GetPoolAssetAndIndex(coin.Denom)
	if err != nil {
		return err
	}

	if coin.Amount.LTE(sdk.ZeroInt()) {
		return fmt.Errorf("can't set the pool's balance of a token to be zero or negative")
	}

	// Update the supply of the asset
	existingAsset.Token = coin
	p.PoolAssets[assetIndex] = &existingAsset
	return nil
}
