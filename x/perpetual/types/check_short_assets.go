package types

import (
	errorsmod "cosmossdk.io/errors"
)

func CheckShortAssets(collateralAsset string, borrowAsset string, baseCurrency string) error {
	// You shouldn't be shorting the base currency (like USDC).
	if borrowAsset == baseCurrency {
		return errorsmod.Wrap(ErrInvalidBorrowingAsset, "cannot short the base currency")
	}

	// If both the collateralAsset and borrowAsset are the same, it doesn't make sense.
	if collateralAsset == borrowAsset {
		return errorsmod.Wrap(ErrInvalidCollateralAsset, "collateral asset cannot be the same as the borrowed asset in a short position")
	}

	// The collateral for a short must be the base currency.
	if collateralAsset != baseCurrency {
		return errorsmod.Wrap(ErrInvalidCollateralAsset, "collateral asset for a short position must be the base currency")
	}

	return nil
}
