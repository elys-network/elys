package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/x/margin/types"
)

func (k Keeper) Open(ctx sdk.Context, msg *types.MsgOpen) (*types.MsgOpenResponse, error) {
	// Determine the type of position (long or short) and validate assets accordingly.
	switch msg.Position {
	case types.Position_LONG:
		if err := k.OpenChecker.CheckLongAssets(ctx, msg.CollateralAsset, msg.BorrowAsset); err != nil {
			return nil, err
		}
	case types.Position_SHORT:
		if err := k.OpenChecker.CheckShortAssets(ctx, msg.CollateralAsset, msg.BorrowAsset); err != nil {
			return nil, err
		}
	default:
		return nil, sdkerrors.Wrap(types.ErrInvalidPosition, msg.Position.String())
	}

	if err := k.OpenChecker.CheckUserAuthorization(ctx, msg); err != nil {
		return nil, err
	}

	// Check if it is the same direction position for the same trader.
	if mtp := k.OpenChecker.CheckSamePosition(ctx, msg); mtp != nil {
		return k.OpenConsolidate(ctx, mtp, msg)
	}

	if err := k.OpenChecker.CheckMaxOpenPositions(ctx); err != nil {
		return nil, err
	}

	// Get token asset other than base currency
	tradingAsset := k.OpenChecker.GetTradingAsset(msg.CollateralAsset, msg.BorrowAsset)

	// Get pool id, amm pool, and margin pool
	poolId, ammPool, pool, err := k.OpenChecker.PreparePools(ctx, tradingAsset)
	if err != nil {
		return nil, err
	}

	if err := k.OpenChecker.CheckPoolHealth(ctx, poolId); err != nil {
		return nil, err
	}

	var mtp *types.MTP
	switch msg.Position {
	case types.Position_LONG:
		mtp, err = k.OpenChecker.OpenLong(ctx, poolId, msg)
		if err != nil {
			return nil, err
		}
	case types.Position_SHORT:
		mtp, err = k.OpenChecker.OpenShort(ctx, poolId, msg)
		if err != nil {
			return nil, err
		}
	default:
		return nil, sdkerrors.Wrap(types.ErrInvalidPosition, msg.Position.String())
	}

	k.OpenChecker.EmitOpenEvent(ctx, mtp)

	if k.hooks != nil {
		k.hooks.AfterMarginPositionOpended(ctx, ammPool, pool)
	}

	return &types.MsgOpenResponse{}, nil
}
