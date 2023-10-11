package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/x/margin/types"
)

func (k Keeper) CheckPoolHealth(ctx sdk.Context, poolId uint64) error {
	pool, found := k.PoolChecker.GetPool(ctx, poolId)
	if !found {
		return sdkerrors.Wrap(types.ErrInvalidBorrowingAsset, "invalid collateral asset")
	}

	if !k.PoolChecker.IsPoolEnabled(ctx, poolId) || k.PoolChecker.IsPoolClosed(ctx, poolId) {
		return sdkerrors.Wrap(types.ErrMTPDisabled, "pool is disabled or closed")
	}

	if !pool.Health.IsNil() && pool.Health.LTE(k.PoolChecker.GetPoolOpenThreshold(ctx)) {
		return sdkerrors.Wrap(types.ErrInvalidPosition, "pool health too low to open new positions")
	}
	return nil
}
