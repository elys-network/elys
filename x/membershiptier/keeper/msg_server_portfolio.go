package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/x/membershiptier/types"
)

func (k msgServer) CreatePortfolio(goCtx context.Context, msg *types.MsgCreatePortfolio) (*types.MsgCreatePortfolioResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value already exists
	_, isFound := k.GetPortfolio(
		ctx,
		msg.Index,
	)
	if isFound {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "index already set")
	}

	var portfolio = types.Portfolio{
		Creator:   msg.Creator,
		Index:     msg.Index,
		Assetkey:  msg.Assetkey,
		Coinvalue: msg.Coinvalue,
	}

	k.SetPortfolio(
		ctx,
		portfolio,
	)
	return &types.MsgCreatePortfolioResponse{}, nil
}

func (k msgServer) UpdatePortfolio(goCtx context.Context, msg *types.MsgUpdatePortfolio) (*types.MsgUpdatePortfolioResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value exists
	valFound, isFound := k.GetPortfolio(
		ctx,
		msg.Index,
	)
	if !isFound {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, "index not set")
	}

	// Checks if the the msg creator is the same as the current owner
	if msg.Creator != valFound.Creator {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	var portfolio = types.Portfolio{
		Creator:   msg.Creator,
		Index:     msg.Index,
		Assetkey:  msg.Assetkey,
		Coinvalue: msg.Coinvalue,
	}

	k.SetPortfolio(ctx, portfolio)

	return &types.MsgUpdatePortfolioResponse{}, nil
}

func (k msgServer) DeletePortfolio(goCtx context.Context, msg *types.MsgDeletePortfolio) (*types.MsgDeletePortfolioResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value exists
	valFound, isFound := k.GetPortfolio(
		ctx,
		msg.Index,
	)
	if !isFound {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, "index not set")
	}

	// Checks if the the msg creator is the same as the current owner
	if msg.Creator != valFound.Creator {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	k.RemovePortfolio(
		ctx,
		msg.Index,
	)

	return &types.MsgDeletePortfolioResponse{}, nil
}
