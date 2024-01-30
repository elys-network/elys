package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
)

// Get balance of a denom
func GetAmmPoolBalance(ammPool ammtypes.Pool, assetDenom string) (sdk.Int, error) {
	for _, asset := range ammPool.PoolAssets {
		if asset.Token.Denom == assetDenom {
			return asset.Token.Amount, nil
		}
	}

	return sdk.ZeroInt(), errorsmod.Wrap(ErrBalanceNotAvailable, "Balance not available")
}
