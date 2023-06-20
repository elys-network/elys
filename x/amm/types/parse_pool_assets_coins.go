package types

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (p Pool) parsePoolAssetsCoins(tokensA sdk.Coins, tokensB sdk.Coins) (
	Aasset PoolAsset, Basset PoolAsset, err error,
) {
	if len(tokensB) != 1 {
		return Aasset, Basset, errors.New("expected tokensA to be of length one")
	}
	_, Aasset, Basset, err = p.parsePoolAssets(tokensA, tokensB[0].Denom)
	return Aasset, Basset, err
}
