package types

import (
	"cosmossdk.io/math"
	ammtypes "github.com/elys-network/elys/x/amm/types"
)

// Check if amm pool has sufficcient balance
func HasSufficientPoolBalance(ammPool ammtypes.Pool, assetDenom string, requiredAmount math.Int) bool {
	balance, err := GetAmmPoolBalance(ammPool, assetDenom)
	if err != nil {
		return false
	}

	// Balance check
	if balance.GTE(requiredAmount) {
		return true
	}

	return false
}
