package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/leveragelp/types"
)

// Get Leveragelp Pool Balance
func (k Keeper) GetLeveragelpPoolBalances(leveragelpPool types.Pool, denom string) (sdk.Int, sdk.Int, sdk.Int) {
	for _, asset := range leveragelpPool.PoolAssets {
		if asset.AssetDenom == denom {
			return asset.AssetBalance, asset.Liabilities, asset.Custody
		}
	}

	return sdk.ZeroInt(), sdk.ZeroInt(), sdk.ZeroInt()
}
