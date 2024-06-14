package keeper

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (k Keeper) Close(ctx sdk.Context, msg *types.MsgClose) (*types.MsgCloseResponse, error) {
	mtp, err := k.GetMTP(ctx, msg.Creator, msg.Id)
	if err != nil {
		return nil, err
	}

	entry, found := k.assetProfileKeeper.GetEntry(ctx, ptypes.BaseCurrency)
	if !found {
		return nil, errorsmod.Wrapf(assetprofiletypes.ErrAssetProfileNotFound, "asset %s not found", ptypes.BaseCurrency)
	}
	baseCurrency := entry.Denom

	var closedMtp *types.MTP
	var repayAmount math.Int
	switch mtp.Position {
	case types.Position_LONG:
		closedMtp, repayAmount, err = k.CloseLong(ctx, msg, baseCurrency)
		if err != nil {
			return nil, err
		}
	case types.Position_SHORT:
		closedMtp, repayAmount, err = k.CloseShort(ctx, msg, baseCurrency)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errorsmod.Wrap(types.ErrInvalidPosition, mtp.Position.String())
	}

	// Emit close event
	k.EmitCloseEvent(ctx, closedMtp, repayAmount)

	return &types.MsgCloseResponse{
		Id:     closedMtp.Id,
		Amount: repayAmount,
	}, nil
}

func (k Keeper) EmitCloseEvent(ctx sdk.Context, mtp *types.MTP, repayAmount math.Int) {
	ctx.EventManager().EmitEvent(types.GenerateCloseEvent(mtp, repayAmount))
}
