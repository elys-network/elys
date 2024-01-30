package types

import (
	errorsmod "cosmossdk.io/errors"
)

func CheckLongAssets(collateralAsset string, borrowAsset string, baseCurrency string) error {
	if borrowAsset == baseCurrency {
		return errorsmod.Wrap(ErrInvalidBorrowingAsset, "invalid borrowing asset")
	}

	if collateralAsset == borrowAsset && collateralAsset == baseCurrency {
		return errorsmod.Wrap(ErrInvalidBorrowingAsset, "invalid borrowing asset")
	}

	if collateralAsset != borrowAsset && collateralAsset != baseCurrency {
		return errorsmod.Wrap(ErrInvalidBorrowingAsset, "invalid borrowing asset")
	}

	return nil
}
