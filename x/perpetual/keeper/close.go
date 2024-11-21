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
	entry, found := k.assetProfileKeeper.GetEntry(ctx, ptypes.BaseCurrency)
	if !found {
		return nil, errorsmod.Wrapf(assetprofiletypes.ErrAssetProfileNotFound, "asset %s not found", ptypes.BaseCurrency)
	}
	baseCurrency := entry.Denom

	closedMtp, repayAmount, closingRatio, err := k.ClosePosition(ctx, msg, baseCurrency)
	if err != nil {
		return nil, err
	}

	// Emit close event
	k.EmitCloseEvent(ctx, closedMtp, repayAmount, closingRatio)

	return &types.MsgCloseResponse{
		Id:     closedMtp.Id,
		Amount: repayAmount,
	}, nil
}

func (k Keeper) EmitCloseEvent(ctx sdk.Context, mtp *types.MTP, repayAmount sdkmath.Int, closingRatio sdkmath.LegacyDec) {
	ctx.EventManager().EmitEvent(types.GenerateCloseEvent(mtp, repayAmount, closingRatio))
}
