package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func ValidateCollateralAsset(collateralAsset string, baseCurrency string) error {
	if collateralAsset != baseCurrency {
		return sdkerrors.Wrap(ErrInvalidCollateralAsset, "invalid collateral asset")
	}
	return nil
}
