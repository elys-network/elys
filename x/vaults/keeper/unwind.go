package keeper

import (
	"strconv"
	"strings"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/v5/x/amm/types"
	"github.com/elys-network/elys/v5/x/vaults/types"
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
		exitCoins, _, _, _, _, err := k.amm.ExitPool(ctx, vaultAddress, poolId, coin.Amount, sdk.Coins{}, coin.Denom, false, true)
		if err != nil {
			return errorsmod.Wrapf(types.ErrInvalidAction, "action failed with error: %s", err)
		}

		// Swap these exit coins into targetDenom
		for _, exitCoin := range exitCoins {
			_, err := k.amm.SwapByDenom(ctx, &ammtypes.MsgSwapByDenom{
				Sender:    vaultAddress.String(),
				Recipient: vaultAddress.String(),
				DenomIn:   exitCoin.Denom,
				DenomOut:  targetDenom,
				Amount:    exitCoin,
				MinAmount: sdk.NewCoin(targetDenom, sdkmath.NewInt(0)),
				MaxAmount: sdk.NewCoin(targetDenom, sdkmath.NewInt(0)), // TODO: check here, swap order
			})
			if err != nil {
				return errorsmod.Wrapf(types.ErrInvalidAction, "action failed with error: %s", err)
			}
		}

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
