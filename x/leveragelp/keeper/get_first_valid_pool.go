package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/x/leveragelp/types"
)

func (k Keeper) GetFirstValidPool(ctx sdk.Context, borrowAsset string) (uint64, error) {
	poolIds := k.amm.GetAllPoolIdsWithDenom(ctx, borrowAsset)
	if len(poolIds) < 1 {
		return 0, sdkerrors.Wrap(types.ErrInvalidBorrowingAsset, "invalid collateral asset")
	}
	return poolIds[0], nil
}
