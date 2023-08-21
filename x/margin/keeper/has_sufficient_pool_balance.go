package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
)

// Check if amm pool has sufficcient balance
func (k Keeper) HasSufficientPoolBalance(ctx sdk.Context, ammPool ammtypes.Pool, assetDenom string, requiredAmount sdk.Int) bool {
	balance, err := k.GetAmmPoolBalance(ctx, ammPool, assetDenom)
	if err != nil {
		return false
	}

	// Balance check
	if balance.GTE(requiredAmount) {
		return true
	}

	return false
}
