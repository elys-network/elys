package types

import (
	"cosmossdk.io/math"
	ammtypes "github.com/elys-network/elys/v6/x/amm/types"
)

// Check if amm pool has sufficcient balance
func HasSufficientPoolBalance(ammPool ammtypes.Pool, assetDenom string, requiredAmount math.Int) bool {
	balance, err := ammPool.GetAmmPoolBalance(assetDenom)
	if err != nil {
		return false
	}

	// Balance check
	if balance.GTE(requiredAmount) {
		return true
	}

	return false
}
