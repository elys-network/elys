package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	aptypes "github.com/elys-network/elys/x/assetprofile/types"
	"github.com/elys-network/elys/x/commitment/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

// UncommitTokens uncommits the tokens from committed store and make it liquid immediately.
func (k msgServer) UncommitTokens(goCtx context.Context, msg *types.MsgUncommitTokens) (*types.MsgUncommitTokensResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	addr, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, "unable to convert address from bech32")
	}

	assetProfile, found := k.assetProfileKeeper.GetEntry(ctx, msg.Denom)
	if !found {
		return nil, errorsmod.Wrapf(aptypes.ErrAssetProfileNotFound, "denom: %s", msg.Denom)
	}

	if !assetProfile.WithdrawEnabled {
		return nil, errorsmod.Wrapf(types.ErrCommitDisabled, "denom: %s", msg.Denom)
	}

	// Get the Commitments for the creator
	commitments := k.GetCommitments(ctx, msg.Creator)

	if msg.Denom == ptypes.Eden {
		err := k.hooks.BeforeEdenCommitChange(ctx, addr)
		if err != nil {
			return nil, err
		}
	}

	if msg.Denom == ptypes.EdenB {
		err := k.hooks.BeforeEdenBCommitChange(ctx, addr)
		if err != nil {
			return nil, err
		}
	}

	// Deduct from committed tokens
	err = commitments.DeductFromCommitted(msg.Denom, msg.Amount, uint64(ctx.BlockTime().Unix()))
	if err != nil {
		return nil, err
	}
	k.SetCommitments(ctx, commitments)

	liquidCoins := sdk.NewCoins(sdk.NewCoin(msg.Denom, msg.Amount))

	err = k.HandleWithdrawFromCommitment(ctx, &commitments, liquidCoins, true, addr)
	if err != nil {
		return nil, err
	}

	// Emit Hook if Eden is uncommitted
	if msg.Denom == ptypes.Eden {
		err = k.EdenUncommitted(ctx, msg.Creator, sdk.NewCoin(msg.Denom, msg.Amount))
		if err != nil {
			return nil, err
		}
	}

	// Emit Hook commitment changed
	err = k.AfterCommitmentChange(ctx, msg.Creator, sdk.Coins{sdk.NewCoin(msg.Denom, msg.Amount)})
	if err != nil {
		return nil, err
	}

	// Update total commitment
	params := k.GetParams(ctx)
	params.TotalCommitted = params.TotalCommitted.Add(liquidCoins...)
	k.SetParams(ctx, params)

	// Emit blockchain event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCommitmentChanged,
			sdk.NewAttribute(types.AttributeCreator, msg.Creator),
			sdk.NewAttribute(types.AttributeAmount, msg.Amount.String()),
			sdk.NewAttribute(types.AttributeDenom, msg.Denom),
		),
	)

	return &types.MsgUncommitTokensResponse{}, nil
}
