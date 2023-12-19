package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (pool Pool) Validate(
	poolId uint64,
) error {
	if pool.GetPoolId() != poolId {
		return errorsmod.Wrapf(ErrInvalidPool,
			"Pool was attempted to be created with incorrect pool ID.")
	}
	address, err := sdk.AccAddressFromBech32(pool.GetAddress())
	if err != nil {
		return errorsmod.Wrapf(ErrInvalidPool,
			"Pool was attempted to be created with invalid pool address.")
	}
	if !address.Equals(NewPoolAddress(poolId)) {
		return errorsmod.Wrapf(ErrInvalidPool,
			"Pool was attempted to be created with incorrect pool address.")
	}
	return nil
}
