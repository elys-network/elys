package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/x/margin/types"
)

func (k Keeper) CheckLongAssets(ctx sdk.Context, collateralAsset string, borrowAsset string, baseCurrency string) error {
	if borrowAsset == baseCurrency {
		return sdkerrors.Wrap(types.ErrInvalidBorrowingAsset, "invalid borrowing asset")
	}

	if collateralAsset == borrowAsset && collateralAsset == baseCurrency {
		return sdkerrors.Wrap(types.ErrInvalidBorrowingAsset, "invalid borrowing asset")
	}

	if collateralAsset != borrowAsset && collateralAsset != baseCurrency {
		return sdkerrors.Wrap(types.ErrInvalidBorrowingAsset, "invalid borrowing asset")
	}

	return nil
}
