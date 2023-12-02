package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/x/margin/types"
)

func (k Keeper) GetFirstValidPool(ctx sdk.Context, collateralAsset string, tradingAsset string) (uint64, error) {
	denoms := []string{collateralAsset, tradingAsset}
	poolId, found := k.amm.GetPoolIdWithAllDenoms(ctx, denoms)
	if !found {
		return 0, sdkerrors.Wrap(types.ErrPoolDoesNotExist, fmt.Sprintf("%s", denoms))
	}
	return poolId, nil
}
