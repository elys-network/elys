package keeper

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (k Keeper) GetFirstValidPool(ctx sdk.Context, collateralAsset string, tradingAsset string) (uint64, error) {
	denoms := []string{collateralAsset, tradingAsset}
	poolId, found := k.amm.GetPoolIdWithAllDenoms(ctx, denoms)
	if !found {
		return 0, errorsmod.Wrap(types.ErrPoolDoesNotExist, fmt.Sprintf("%s", denoms))
	}
	return poolId, nil
}
