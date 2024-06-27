package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/leveragelp/types"
)

func (k msgServer) AddCollateral(goCtx context.Context, msg *types.MsgAddCollateral) (*types.MsgAddCollateralResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	err := k.ProcessAddCollateral(ctx, msg.Creator, msg.Id, msg.Collateral)
	if err != nil {
		return nil, err
	}

	return &types.MsgAddCollateralResponse{}, nil
}
