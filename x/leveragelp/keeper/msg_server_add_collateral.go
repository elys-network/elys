package keeper

import (
	"context"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/leveragelp/types"
)

func (k msgServer) AddCollateral(goCtx context.Context, msg *types.MsgAddCollateral) (*types.MsgAddCollateralResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	err := k.ProcessAddCollateral(ctx, msg.Creator, msg.Id, msg.Collateral)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(types.EventAddCollateral,
		sdk.NewAttribute("position_index", strconv.FormatInt(int64(msg.Id), 10)),
		sdk.NewAttribute("creator", msg.Creator),
		sdk.NewAttribute("collateral", msg.Collateral.String()),
	))

	return &types.MsgAddCollateralResponse{}, nil
}
