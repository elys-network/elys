package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/x/margin/types"
)

func (k msgServer) Open(goCtx context.Context, msg *types.MsgOpen) (*types.MsgOpenResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := k.CheckUserAuthorization(ctx, msg); err != nil {
		return nil, err
	}

	if err := k.CheckMaxOpenPositions(ctx); err != nil {
		return nil, err
	}

	if err := k.ValidateCollateralAsset(msg.CollateralAsset); err != nil {
		return nil, err
	}

	poolId, err := k.GetFirstValidPool(ctx, msg.BorrowAsset)
	if err != nil {
		return nil, err
	}

	if err := k.CheckPoolHealth(ctx, poolId); err != nil {
		return nil, err
	}

	var mtp *types.MTP
	switch msg.Position {
	case types.Position_LONG:
		mtp, err = k.OpenLong(ctx, poolId, msg)
		if err != nil {
			return nil, err
		}
	default:
		return nil, sdkerrors.Wrap(types.ErrInvalidPosition, msg.Position.String())
	}

	ctx.EventManager().EmitEvent(k.GenerateOpenEvent(mtp))
	return &types.MsgOpenResponse{}, nil
}
