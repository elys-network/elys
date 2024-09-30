package types

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	ammtypes "github.com/elys-network/elys/x/amm/types"
)

// Get balance of a denom
func GetAmmPoolBalance(ammPool ammtypes.Pool, assetDenom string) (math.Int, error) {
	for _, asset := range ammPool.PoolAssets {
		if asset.Token.Denom == assetDenom {
			return asset.Token.Amount, nil
		}
	}

	return math.ZeroInt(), errorsmod.Wrap(ErrBalanceNotAvailable, "Balance not available")
}
