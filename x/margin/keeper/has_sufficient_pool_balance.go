package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
)

func (k Keeper) HasSufficientPoolBalance(ctx sdk.Context, ammPool ammtypes.Pool, borrowAsset string, requiredAmount sdk.Int) bool {
	for _, asset := range ammPool.PoolAssets {
		if borrowAsset == asset.Token.Denom && asset.Token.Amount.GTE(requiredAmount) {
			return true
		}
	}
	return false
}
