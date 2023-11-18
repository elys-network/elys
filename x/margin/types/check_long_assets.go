package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func CheckLongAssets(collateralAsset string, borrowAsset string, baseCurrency string) error {
	if borrowAsset == baseCurrency {
		return sdkerrors.Wrap(ErrInvalidBorrowingAsset, "invalid borrowing asset")
	}

	if collateralAsset == borrowAsset && collateralAsset == baseCurrency {
		return sdkerrors.Wrap(ErrInvalidBorrowingAsset, "invalid borrowing asset")
	}

	if collateralAsset != borrowAsset && collateralAsset != baseCurrency {
		return sdkerrors.Wrap(ErrInvalidBorrowingAsset, "invalid borrowing asset")
	}

	return nil
}
