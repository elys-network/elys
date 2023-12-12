package keeper

import (
	"context"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/elys-network/elys/x/tokenomics/types"
)

func (k msgServer) CreateAirdrop(goCtx context.Context, msg *types.MsgCreateAirdrop) (*types.MsgCreateAirdropResponse, error) {
	if k.authority != msg.Authority {
		return nil, errors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value already exists
	_, isFound := k.GetAirdrop(
		ctx,
		msg.Intent,
	)
	if isFound {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "index already set")
	}

	airdrop := types.Airdrop{
		Authority: msg.Authority,
		Intent:    msg.Intent,
		Amount:    msg.Amount,
	}

	k.SetAirdrop(
		ctx,
		airdrop,
	)
	return &types.MsgCreateAirdropResponse{}, nil
}

func (k msgServer) UpdateAirdrop(goCtx context.Context, msg *types.MsgUpdateAirdrop) (*types.MsgUpdateAirdropResponse, error) {
	if k.authority != msg.Authority {
		return nil, errors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value exists
	valFound, isFound := k.GetAirdrop(
		ctx,
		msg.Intent,
	)
	if !isFound {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, "index not set")
	}

	// Checks if the the msg authority is the same as the current owner
	if msg.Authority != valFound.Authority {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	airdrop := types.Airdrop{
		Authority: msg.Authority,
		Intent:    msg.Intent,
		Amount:    msg.Amount,
	}

	k.SetAirdrop(ctx, airdrop)

	return &types.MsgUpdateAirdropResponse{}, nil
}

func (k msgServer) DeleteAirdrop(goCtx context.Context, msg *types.MsgDeleteAirdrop) (*types.MsgDeleteAirdropResponse, error) {
	if k.authority != msg.Authority {
		return nil, errors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value exists
	valFound, isFound := k.GetAirdrop(
		ctx,
		msg.Intent,
	)
	if !isFound {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, "index not set")
	}

	// Checks if the the msg authority is the same as the current owner
	if msg.Authority != valFound.Authority {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	k.RemoveAirdrop(
		ctx,
		msg.Intent,
	)

	return &types.MsgDeleteAirdropResponse{}, nil
}
