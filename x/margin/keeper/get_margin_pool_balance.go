package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/margin/types"
)

// Get Margin Pool Balance
func (k Keeper) GetMarginPoolBalances(marginPool types.Pool, denom string) (sdk.Int, sdk.Int, sdk.Int) {
	for _, asset := range marginPool.PoolAssets {
		if asset.AssetDenom == denom {
			return asset.AssetBalance, asset.Liabilities, asset.Custody
		}
	}

	return sdk.ZeroInt(), sdk.ZeroInt(), sdk.ZeroInt()
}
