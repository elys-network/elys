package keeper

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/leveragelp/types"
)

func (k Keeper) CheckUserAuthorization(ctx sdk.Context, msg *types.MsgOpen) error {
	if k.IsWhitelistingEnabled(ctx) && !k.CheckIfWhitelisted(ctx, msg.Creator) {
		return errorsmod.Wrap(types.ErrUnauthorised, "unauthorised")
	}
	return nil
}

func (k Keeper) CheckSamePosition(ctx sdk.Context, msg *types.MsgOpen) *types.Position {
	positions, _, err := k.GetPositionsForAddress(ctx, sdk.MustAccAddressFromBech32(msg.Creator), &query.PageRequest{})
	if err != nil {
		return nil
	}
	for _, position := range positions {
		if position.Position.AmmPoolId == msg.AmmPoolId && position.Position.Collateral.Denom == msg.CollateralAsset {
			return position.Position
		}
	}

	return nil
}

func (k Keeper) CheckPoolHealth(ctx sdk.Context, poolId uint64) error {
	pool, found := k.GetPool(ctx, poolId)
	if !found {
		return errorsmod.Wrap(types.ErrInvalidBorrowingAsset, "invalid collateral asset")
	}

	if !k.IsPoolEnabled(ctx, poolId) || k.IsPoolClosed(ctx, poolId) {
		return errorsmod.Wrap(types.ErrPositionDisabled, "pool is disabled or closed")
	}

	if !pool.Health.IsNil() && pool.Health.LTE(k.GetPoolOpenThreshold(ctx)) {
		return errorsmod.Wrap(types.ErrInvalidPosition, "pool health too low to open new positions")
	}
	return nil
}

func (k Keeper) CheckMaxOpenPositions(ctx sdk.Context) error {
	if k.GetOpenPositionCount(ctx) >= k.GetMaxOpenPositions(ctx) {
		return errorsmod.Wrap(types.ErrMaxOpenPositions, "cannot open new positions")
	}
	return nil
}

func (k Keeper) GetAmmPool(ctx sdk.Context, poolId uint64) (ammtypes.Pool, error) {
	ammPool, found := k.amm.GetPool(ctx, poolId)
	if !found {
		return ammPool, errorsmod.Wrap(types.ErrPoolDoesNotExist, fmt.Sprintf("poolId: %d", poolId))
	}
	return ammPool, nil
}
