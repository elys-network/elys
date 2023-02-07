package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/x/assetprofile/types"
)

func (k msgServer) CreateEntry(goCtx context.Context, msg *types.MsgCreateEntry) (*types.MsgCreateEntryResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value already exists
	_, isFound := k.GetEntry(
		ctx,
		msg.BaseDenom,
	)
	if isFound {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "index already set")
	}

	var entry = types.Entry{
		Authority:                msg.Authority,
		BaseDenom:                msg.BaseDenom,
		Decimals:                 msg.Decimals,
		Denom:                    msg.Denom,
		Path:                     msg.Path,
		IbcChannelId:             msg.IbcChannelId,
		IbcCounterpartyChannelId: msg.IbcCounterpartyChannelId,
		DisplayName:              msg.DisplayName,
		DisplaySymbol:            msg.DisplaySymbol,
		Network:                  msg.Network,
		Address:                  msg.Address,
		ExternalSymbol:           msg.ExternalSymbol,
		TransferLimit:            msg.TransferLimit,
		Permissions:              msg.Permissions,
		UnitDenom:                msg.UnitDenom,
		IbcCounterpartyDenom:     msg.IbcCounterpartyDenom,
		IbcCounterpartyChainId:   msg.IbcCounterpartyChainId,
	}

	k.SetEntry(
		ctx,
		entry,
	)
	return &types.MsgCreateEntryResponse{}, nil
}

func (k msgServer) UpdateEntry(goCtx context.Context, msg *types.MsgUpdateEntry) (*types.MsgUpdateEntryResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value exists
	valFound, isFound := k.GetEntry(
		ctx,
		msg.BaseDenom,
	)
	if !isFound {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, "index not set")
	}

	// Checks if the the msg authority is the same as the current owner
	if msg.Authority != valFound.Authority {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	var entry = types.Entry{
		Authority:                msg.Authority,
		BaseDenom:                msg.BaseDenom,
		Decimals:                 msg.Decimals,
		Denom:                    msg.Denom,
		Path:                     msg.Path,
		IbcChannelId:             msg.IbcChannelId,
		IbcCounterpartyChannelId: msg.IbcCounterpartyChannelId,
		DisplayName:              msg.DisplayName,
		DisplaySymbol:            msg.DisplaySymbol,
		Network:                  msg.Network,
		Address:                  msg.Address,
		ExternalSymbol:           msg.ExternalSymbol,
		TransferLimit:            msg.TransferLimit,
		Permissions:              msg.Permissions,
		UnitDenom:                msg.UnitDenom,
		IbcCounterpartyDenom:     msg.IbcCounterpartyDenom,
		IbcCounterpartyChainId:   msg.IbcCounterpartyChainId,
	}

	k.SetEntry(ctx, entry)

	return &types.MsgUpdateEntryResponse{}, nil
}

func (k msgServer) DeleteEntry(goCtx context.Context, msg *types.MsgDeleteEntry) (*types.MsgDeleteEntryResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value exists
	valFound, isFound := k.GetEntry(
		ctx,
		msg.BaseDenom,
	)
	if !isFound {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, "index not set")
	}

	// Checks if the the msg authority is the same as the current owner
	if msg.Authority != valFound.Authority {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	k.RemoveEntry(
		ctx,
		msg.BaseDenom,
	)

	return &types.MsgDeleteEntryResponse{}, nil
}
