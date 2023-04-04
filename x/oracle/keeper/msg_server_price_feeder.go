package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/x/oracle/types"
)

func (k msgServer) CreatePriceFeeder(goCtx context.Context, msg *types.MsgCreatePriceFeeder) (*types.MsgCreatePriceFeederResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value already exists
	_, isFound := k.GetPriceFeeder(ctx, msg.Feeder)
	if isFound {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "index already set")
	}

	var priceFeeder = types.PriceFeeder{
		Feeder:   msg.Feeder,
		IsActive: msg.IsActive,
	}

	k.SetPriceFeeder(
		ctx,
		priceFeeder,
	)
	return &types.MsgCreatePriceFeederResponse{}, nil
}

func (k msgServer) UpdatePriceFeeder(goCtx context.Context, msg *types.MsgUpdatePriceFeeder) (*types.MsgUpdatePriceFeederResponse, error) {
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

	k.SetPriceFeeder(ctx, priceFeeder)

	return &types.MsgUpdatePriceFeederResponse{}, nil
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
