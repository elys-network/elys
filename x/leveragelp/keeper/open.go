package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/leveragelp/types"
)

func (k Keeper) Open(ctx sdk.Context, msg *types.MsgOpen) (*types.MsgOpenResponse, error) {

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

	if err := k.OpenChecker.CheckPoolHealth(ctx, msg.AmmPoolId); err != nil {
		return nil, err
	}

	mtp, err := k.OpenChecker.OpenLong(ctx, msg.AmmPoolId, msg)
	if err != nil {
		return nil, err
	}

	k.OpenChecker.EmitOpenEvent(ctx, mtp)

	if k.hooks != nil {
		k.hooks.AfterLeveragelpPositionOpen(ctx, ammPool, pool)
	}

	return &types.MsgOpenResponse{}, nil
}
