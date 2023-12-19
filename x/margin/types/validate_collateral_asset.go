package types

import (
	errorsmod "cosmossdk.io/errors"
)

func ValidateCollateralAsset(collateralAsset string, baseCurrency string) error {
	if collateralAsset != baseCurrency {
		return errorsmod.Wrap(ErrInvalidCollateralAsset, "invalid collateral asset")
	}
	return nil
}
