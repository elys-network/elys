package keeper

import (
	"context"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/leveragelp/types"
)

func (k msgServer) Open(goCtx context.Context, msg *types.MsgOpen) (*types.MsgOpenResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	return k.Keeper.Open(ctx, msg)
}

func (k Keeper) Open(ctx sdk.Context, msg *types.MsgOpen) (*types.MsgOpenResponse, error) {
	if err := k.CheckUserAuthorization(ctx, msg); err != nil {
		return nil, err
	}

	// Check if it is the same direction position for the same trader.
	if mtp := k.CheckSamePosition(ctx, msg); mtp != nil {
		return k.OpenConsolidate(ctx, mtp, msg)
	}

	if err := k.CheckMaxOpenPositions(ctx); err != nil {
		return nil, err
	}

	if err := k.CheckPoolHealth(ctx, msg.AmmPoolId); err != nil {
		return nil, err
	}

	mtp, err := k.OpenLong(ctx, msg.AmmPoolId, msg)
	if err != nil {
		return nil, err
	}

	event := sdk.NewEvent(types.EventOpen,
		sdk.NewAttribute("id", strconv.FormatInt(int64(mtp.Id), 10)),
		sdk.NewAttribute("address", mtp.Address),
		sdk.NewAttribute("collateral", mtp.Collateral.String()),
		sdk.NewAttribute("leverage", mtp.Leverage.String()),
		sdk.NewAttribute("liabilities", mtp.Liabilities.String()),
		sdk.NewAttribute("health", mtp.MtpHealth.String()),
	)
	ctx.EventManager().EmitEvent(event)

	return &types.MsgOpenResponse{}, nil
}
