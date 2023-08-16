package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/margin/types"
)

func (k Keeper) GetAmmPool(ctx sdk.Context, poolId uint64, borrowAsset string) (ammtypes.Pool, error) {
	ammPool, found := k.amm.GetPool(ctx, poolId)
	if !found {
		return ammPool, sdkerrors.Wrap(types.ErrPoolDoesNotExist, borrowAsset)
	}
	return ammPool, nil
}
