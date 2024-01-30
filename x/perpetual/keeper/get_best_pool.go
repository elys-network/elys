package keeper

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (k Keeper) GetBestPool(ctx sdk.Context, collateralAsset string, tradingAsset string) (uint64, error) {
	denoms := []string{collateralAsset, tradingAsset}
	pool, found := k.amm.GetBestPoolWithDenoms(ctx, denoms)
	if !found {
		return 0, errorsmod.Wrap(types.ErrPoolDoesNotExist, fmt.Sprintf("%s", denoms))
	}
	return pool.PoolId, nil
}
