package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/leveragelp/types"
)

func (k Keeper) CheckUserAuthorization(ctx sdk.Context, msg *types.MsgOpen) error {
	if k.IsWhitelistingEnabled(ctx) && !k.CheckIfWhitelisted(ctx, msg.Creator) {
		return sdkerrors.Wrap(types.ErrUnauthorised, "unauthorised")
	}
	return nil
}

func (k Keeper) CheckSamePosition(ctx sdk.Context, msg *types.MsgOpen) *types.Position {
	positions := k.GetAllPositions(ctx)
	for _, position := range positions {
		if position.Address == msg.Creator {
			return &position
		}
	}

	return nil
}

func (k Keeper) CheckPoolHealth(ctx sdk.Context, poolId uint64) error {
	pool, found := k.GetPool(ctx, poolId)
	if !found {
		return sdkerrors.Wrap(types.ErrInvalidBorrowingAsset, "invalid collateral asset")
	}

	if !k.IsPoolEnabled(ctx, poolId) || k.IsPoolClosed(ctx, poolId) {
		return sdkerrors.Wrap(types.ErrPositionDisabled, "pool is disabled or closed")
	}

	if !pool.Health.IsNil() && pool.Health.LTE(k.GetPoolOpenThreshold(ctx)) {
		return sdkerrors.Wrap(types.ErrInvalidPosition, "pool health too low to open new positions")
	}
	return nil
}

func (k Keeper) CheckMaxOpenPositions(ctx sdk.Context) error {
	if k.GetOpenPositionCount(ctx) >= k.GetMaxOpenPositions(ctx) {
		return sdkerrors.Wrap(types.ErrMaxOpenPositions, "cannot open new positions")
	}
	return nil
}

func (k Keeper) GetAmmPool(ctx sdk.Context, poolId uint64) (ammtypes.Pool, error) {
	ammPool, found := k.amm.GetPool(ctx, poolId)
	if !found {
		return ammPool, sdkerrors.Wrap(types.ErrPoolDoesNotExist, fmt.Sprintf("poolId: %d", poolId))
	}
	return ammPool, nil
}
