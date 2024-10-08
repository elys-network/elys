package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/tradeshield/types"
)

func (k msgServer) CreatePendingSpotOrder(goCtx context.Context, msg *types.MsgCreatePendingSpotOrder) (*types.MsgCreatePendingSpotOrderResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	var pendingSpotOrder = types.SpotOrder{
		OrderType:    msg.OrderType,
		OrderId:      uint64(0),
		OrderPrice:   msg.OrderPrice,
		OrderAmount:  *msg.OrderAmount,
		OwnerAddress: msg.OwnerAddress,
	}

	id := k.AppendPendingSpotOrder(
		ctx,
		pendingSpotOrder,
	)

	return &types.MsgCreatePendingSpotOrderResponse{
		OrderId: id,
	}, nil
}

func (k msgServer) UpdatePendingSpotOrder(goCtx context.Context, msg *types.MsgUpdatePendingSpotOrder) (*types.MsgUpdatePendingSpotOrderResponse, error) {
	// _ := sdk.UnwrapSDKContext(goCtx)

	return &types.MsgUpdatePendingSpotOrderResponse{}, nil
}
