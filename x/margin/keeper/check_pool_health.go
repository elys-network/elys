package keeper

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/margin/types"
)

func (k Keeper) CheckPoolHealth(ctx sdk.Context, poolId uint64) error {
	pool, found := k.PoolChecker.GetPool(ctx, poolId)
	if !found {
		return errorsmod.Wrap(types.ErrInvalidBorrowingAsset, "invalid collateral asset")
	}

	if !k.PoolChecker.IsPoolEnabled(ctx, poolId) || k.PoolChecker.IsPoolClosed(ctx, poolId) {
		return errorsmod.Wrap(types.ErrMTPDisabled, "pool is disabled or closed")
	}

	if !pool.Health.IsNil() && pool.Health.LTE(k.PoolChecker.GetPoolOpenThreshold(ctx)) {
		return errorsmod.Wrap(types.ErrInvalidPosition, "pool health too low to open new positions")
	}
	return nil
}
