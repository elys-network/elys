package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (k msgServer) BrokerAddCollateral(goCtx context.Context, msg *types.MsgBrokerAddCollateral) (*types.MsgBrokerAddCollateralResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	_ = ctx

	return &types.MsgBrokerAddCollateralResponse{}, nil
}
