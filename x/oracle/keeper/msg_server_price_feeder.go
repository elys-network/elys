package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/oracle/types"
)

func (k msgServer) SetPriceFeeder(goCtx context.Context, msg *types.MsgSetPriceFeeder) (*types.MsgSetPriceFeederResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value exists
	// valFound, isFound := k.GetPriceFeeder(ctx, msg.Feeder)
	// if !isFound {
	// 	return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, "index not set")
	// }

	// Checks if the the msg creator is the same as the current owner
	// if msg.Creator != valFound.Creator {
	// 	return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	// }

	var priceFeeder = types.PriceFeeder{
		Feeder:   msg.Feeder,
		IsActive: msg.IsActive,
	}

	k.Keeper.SetPriceFeeder(ctx, priceFeeder)

	return &types.MsgSetPriceFeederResponse{}, nil
}

func (k msgServer) DeletePriceFeeder(goCtx context.Context, msg *types.MsgDeletePriceFeeder) (*types.MsgDeletePriceFeederResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value exists
	// valFound, isFound := k.GetPriceFeeder(ctx,msg.Feeder,)
	// if !isFound {
	// 	return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, "index not set")
	// }

	// // Checks if the the msg creator is the same as the current owner
	// if msg.Creator != valFound.Creator {
	// 	return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	// }

	k.RemovePriceFeeder(ctx, msg.Feeder)

	return &types.MsgDeletePriceFeederResponse{}, nil
}
