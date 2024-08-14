package types

import (
	errorsmod "cosmossdk.io/errors"
)

func CheckLongAssets(collateralAsset string, borrowAsset string, baseCurrency string) error {
	if borrowAsset == baseCurrency {
		return errorsmod.Wrap(ErrInvalidBorrowingAsset, "invalid operation: the borrowed asset cannot be the base currency")
	}

	if collateralAsset == borrowAsset && collateralAsset == baseCurrency {
		return errorsmod.Wrap(ErrInvalidBorrowingAsset, "invalid operation: collateral and borrowed assets cannot both be the base currency")
	}

	if collateralAsset != borrowAsset && collateralAsset != baseCurrency {
		return errorsmod.Wrap(ErrInvalidBorrowingAsset, "invalid collateral: collateral must either match the borrowed asset or be the base currency")
	}

	return nil
}

func CheckShortAssets(collateralAsset string, borrowAsset string, baseCurrency string) error {
	// You shouldn't be shorting the base currency (like USDC).
	if borrowAsset == baseCurrency {
		return errorsmod.Wrap(ErrInvalidBorrowingAsset, "borrowing not allowed: cannot take a short position against the base currency")
	}

	// If both the collateralAsset and borrowAsset are the same, it doesn't make sense.
	if collateralAsset == borrowAsset {
		return errorsmod.Wrap(ErrInvalidCollateralAsset, "invalid operation: collateral asset cannot be identical to the borrowed asset for a short position")
	}

	// The collateral for a short must be the base currency.
	if collateralAsset != baseCurrency {
		return errorsmod.Wrap(ErrInvalidCollateralAsset, "invalid collateral: the collateral asset for a short position must be the base currency")
	}

	return nil
}
