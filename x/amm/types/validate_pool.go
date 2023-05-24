package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (pool Pool) Validate(
	poolId uint64,
) error {
	if pool.GetPoolId() != poolId {
		return sdkerrors.Wrapf(ErrInvalidPool,
			"Pool was attempted to be created with incorrect pool ID.")
	}
	address, err := sdk.AccAddressFromBech32(pool.GetAddress())
	if err != nil {
		return sdkerrors.Wrapf(ErrInvalidPool,
			"Pool was attempted to be created with invalid pool address.")
	}
	if !address.Equals(NewPoolAddress(poolId)) {
		return sdkerrors.Wrapf(ErrInvalidPool,
			"Pool was attempted to be created with incorrect pool address.")
	}
	return nil
}
