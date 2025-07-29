package keeper

import (
	"context"
	"strings"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	"github.com/elys-network/elys/v7/x/assetprofile/types"
)

func (k msgServer) AddEntry(goCtx context.Context, msg *types.MsgAddEntry) (*types.MsgAddEntryResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the entry already exists
	_, isFound := k.GetEntry(ctx, msg.BaseDenom)
	if isFound {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "entry already set")
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

	entry := types.Entry{
		Authority:                k.authority,
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
	return &types.MsgAddEntryResponse{}, nil
}
