package keeper

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (k Keeper) CheckPoolHealth(ctx sdk.Context, poolId uint64) error {
	pool, found := k.GetPool(ctx, poolId)
	if !found {
		return types.ErrPoolDoesNotExist
	}

	if !k.IsPoolEnabled(ctx, poolId) || k.IsPoolClosed(ctx, poolId) {
		return errorsmod.Wrap(types.ErrMTPDisabled, "pool is disabled or closed")
	}

	if !pool.Health.IsNil() && pool.Health.LTE(k.GetPoolOpenThreshold(ctx)) {
		return errorsmod.Wrap(types.ErrInvalidPosition, "pool health too low to open new positions")
	}
	return nil
}
