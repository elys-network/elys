package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// poolAssetsCoins returns all the coins corresponding to a slice of pool assets.
func poolAssetsCoins(assets []PoolAsset) sdk.Coins {
	coins := sdk.Coins{}
	for _, asset := range assets {
		coins = coins.Add(asset.Token)
	}
	return coins
}
