package keeper

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/x/margin/types"
)

func (k Keeper) ValidateCollateralAsset(collateralAsset string, baseCurrency string) error {
	if collateralAsset != baseCurrency {
		return sdkerrors.Wrap(types.ErrInvalidCollateralAsset, "invalid collateral asset")
	}
	return nil
}
