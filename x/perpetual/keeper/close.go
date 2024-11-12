package keeper

import (
	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (k Keeper) Close(ctx sdk.Context, msg *types.MsgClose) (*types.MsgCloseResponse, error) {
	creator := sdk.MustAccAddressFromBech32(msg.Creator)
	mtp, err := k.GetMTP(ctx, creator, msg.Id)
	if err != nil {
		return nil, err
	}

	entry, found := k.assetProfileKeeper.GetEntry(ctx, ptypes.BaseCurrency)
	if !found {
		return nil, errorsmod.Wrapf(assetprofiletypes.ErrAssetProfileNotFound, "asset %s not found", ptypes.BaseCurrency)
	}
	baseCurrency := entry.Denom

	var closedMtp *types.MTP
	var repayAmount sdkmath.Int
	switch mtp.Position {
	case types.Position_LONG:
		closedMtp, repayAmount, err = k.ClosePosition(ctx, msg, baseCurrency)
		if err != nil {
			return nil, err
		}
	case types.Position_SHORT:
		closedMtp, repayAmount, err = k.ClosePosition(ctx, msg, baseCurrency)
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

func (k Keeper) EmitCloseEvent(ctx sdk.Context, mtp *types.MTP, repayAmount sdkmath.Int) {
	ctx.EventManager().EmitEvent(types.GenerateCloseEvent(mtp, repayAmount))
}
