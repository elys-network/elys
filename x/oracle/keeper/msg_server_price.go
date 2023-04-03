package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/x/oracle/types"
)

func (k msgServer) CreatePrice(goCtx context.Context, msg *types.MsgCreatePrice) (*types.MsgCreatePriceResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value already exists
	_, isFound := k.GetPrice(
		ctx,
		msg.Asset,
	)
	if isFound {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "index already set")
	}

	var price = types.Price{
		Provider: msg.Provider,
		Asset:    msg.Asset,
		Price:    msg.Price,
		Source:   msg.Source,
	}

	k.SetPrice(ctx, price)
	return &types.MsgCreatePriceResponse{}, nil
}

func (k msgServer) UpdatePrice(goCtx context.Context, msg *types.MsgUpdatePrice) (*types.MsgUpdatePriceResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value exists
	valFound, isFound := k.GetPrice(
		ctx,
		msg.Asset,
	)
	if !isFound {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, "index not set")
	}

	// Checks if the the msg creator is the same as the current owner
	if msg.Provider != valFound.Provider {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	var price = types.Price{
		Provider: msg.Provider,
		Asset:    msg.Asset,
	}

	k.SetPrice(ctx, price)

	return &types.MsgUpdatePriceResponse{}, nil
}

func (k msgServer) DeletePrice(goCtx context.Context, msg *types.MsgDeletePrice) (*types.MsgDeletePriceResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value exists
	valFound, isFound := k.GetPrice(
		ctx,
		msg.Asset,
	)
	if !isFound {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, "index not set")
	}

	// Checks if the the msg creator is the same as the current owner
	if msg.Creator != valFound.Provider {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	k.RemovePrice(
		ctx,
		msg.Asset,
	)

	return &types.MsgDeletePriceResponse{}, nil
}
