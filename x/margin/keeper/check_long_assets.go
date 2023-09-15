package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/x/margin/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

func (k Keeper) CheckLongingAssets(ctx sdk.Context, collateralAsset string, borrowAsset string) error {
	if borrowAsset == ptypes.BaseCurrency {
		return sdkerrors.Wrap(types.ErrInvalidBorrowingAsset, "invalid borrowing asset")
	}

	if collateralAsset == borrowAsset && collateralAsset == ptypes.BaseCurrency {
		return sdkerrors.Wrap(types.ErrInvalidBorrowingAsset, "invalid borrowing asset")
	}

	if collateralAsset != borrowAsset && collateralAsset != ptypes.BaseCurrency {
		return sdkerrors.Wrap(types.ErrInvalidBorrowingAsset, "invalid borrowing asset")
	}

	return nil
}
