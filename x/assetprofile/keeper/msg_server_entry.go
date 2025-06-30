package keeper

import (
	"context"
	"strings"

	"cosmossdk.io/errors"
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	"github.com/elys-network/elys/v6/x/assetprofile/types"
)

func (k msgServer) UpdateEntry(goCtx context.Context, msg *types.MsgUpdateEntry) (*types.MsgUpdateEntryResponse, error) {
	if k.authority != msg.Authority {
		return nil, errors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value exists
	entry, isFound := k.GetEntry(ctx, msg.BaseDenom)
	if !isFound {
		return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, "entry not set")
	}

	// Checks if the msg authority is the same as the current owner
	if msg.Authority != entry.Authority {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	// check the validity of ibc denom & channel
	hash, err := ibctransfertypes.ParseHexHash(strings.TrimPrefix(msg.Denom, "ibc/"))
	if err == nil && k.transferKeeper != nil {
		denomTrace, ok := k.transferKeeper.GetDenomTrace(ctx, hash)
		if !ok {
			return nil, types.ErrNotValidIbcDenom
		}
		if !strings.Contains(denomTrace.Path, msg.IbcChannelId) {
			return nil, types.ErrChannelIdAndDenomHashMismatch
		}
	}

	entry = types.Entry{
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
		CommitEnabled:            msg.CommitEnabled,
		WithdrawEnabled:          msg.WithdrawEnabled,
	}

	k.SetEntry(ctx, entry)

	return &types.MsgUpdateEntryResponse{}, nil
}

func (k msgServer) DeleteEntry(goCtx context.Context, msg *types.MsgDeleteEntry) (*types.MsgDeleteEntryResponse, error) {
	if k.authority != msg.Authority {
		return nil, errors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value exists
	entry, isFound := k.GetEntry(ctx, msg.BaseDenom)
	if !isFound {
		return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, "entry not set")
	}

	// Checks if the msg authority is the same as the current owner
	if msg.Authority != entry.Authority {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	k.RemoveEntry(ctx, msg.BaseDenom)

	return &types.MsgDeleteEntryResponse{}, nil
}
