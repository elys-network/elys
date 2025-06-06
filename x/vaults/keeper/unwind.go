package keeper

import (
	"strconv"
	"strings"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v6/x/vaults/types"
)

func (k Keeper) UnwindCoin(ctx sdk.Context, coin sdk.Coin, targetDenom string) error {
	// Case 1: coin is a LP share
	if strings.HasPrefix(coin.Denom, "amm/pool/") {
		poolId, err := GetPoolIdFromShareDenom(coin.Denom)
		if err != nil {
			return errorsmod.Wrapf(types.ErrInvalidAction, "action failed with error: %s", err)
		}
		vaultAddress := types.NewVaultAddress(poolId)

		// exit pool
		_, _, _, _, _, err = k.amm.ExitPool(ctx, vaultAddress, poolId, coin.Amount, sdk.Coins{}, coin.Denom, false, true)
		if err != nil {
			return errorsmod.Wrapf(types.ErrInvalidAction, "action failed with error: %s", err)
		}

		//now swap required here, only required for fee conversion

	}
	return nil
}

func GetPoolIdFromShareDenom(shareDenom string) (uint64, error) {
	poolId, err := strconv.Atoi(strings.TrimPrefix(shareDenom, "amm/pool/"))
	if err != nil {
		return 0, err
	}
	return uint64(poolId), nil
}
