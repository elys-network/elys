package keeper

import (
	"context"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/leveragelp/types"
)

func (k msgServer) Close(goCtx context.Context, msg *types.MsgClose) (*types.MsgCloseResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	return k.Keeper.Close(ctx, msg)
}

func (k Keeper) Close(ctx sdk.Context, msg *types.MsgClose) (*types.MsgCloseResponse, error) {
	closedPosition, repayAmount, err := k.CloseLong(ctx, msg)
	if err != nil {
		return nil, err
	}

	if k.hooks != nil {
		err := k.hooks.AfterLeverageLpPositionClose(ctx, sdk.MustAccAddressFromBech32(msg.Creator))
		if err != nil {
			return nil, err
		}
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(types.EventClose,
		sdk.NewAttribute("id", strconv.FormatInt(int64(closedPosition.Id), 10)),
		sdk.NewAttribute("address", closedPosition.Address),
		sdk.NewAttribute("collateral", closedPosition.Collateral.String()),
		sdk.NewAttribute("repay_amount", repayAmount.String()),
		sdk.NewAttribute("liabilities", closedPosition.Liabilities.String()),
		sdk.NewAttribute("health", closedPosition.PositionHealth.String()),
	))

	return &types.MsgCloseResponse{}, nil
}
