package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	aptypes "github.com/elys-network/elys/x/assetprofile/types"
	"github.com/elys-network/elys/x/commitment/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

func (k Keeper) UncommitTokens(ctx sdk.Context, addr sdk.AccAddress, denom string, amount math.Int) error {
	assetProfile, found := k.assetProfileKeeper.GetEntry(ctx, denom)
	if !found {
		return errorsmod.Wrapf(aptypes.ErrAssetProfileNotFound, "denom: %s", denom)
	}

	if !assetProfile.WithdrawEnabled {
		return errorsmod.Wrapf(types.ErrCommitDisabled, "denom: %s", denom)
	}

	// Get the Commitments for the creator
	commitments := k.GetCommitments(ctx, addr.String())

	if denom == ptypes.Eden {
		err := k.hooks.BeforeEdenCommitChange(ctx, addr)
		if err != nil {
			return err
		}
	}

	if denom == ptypes.EdenB {
		err := k.hooks.BeforeEdenBCommitChange(ctx, addr)
		if err != nil {
			return err
		}
	}

	// Deduct from committed tokens
	err := commitments.DeductFromCommitted(denom, amount, uint64(ctx.BlockTime().Unix()))
	if err != nil {
		return err
	}
	k.SetCommitments(ctx, commitments)

	liquidCoins := sdk.NewCoins(sdk.NewCoin(denom, amount))

	err = k.HandleWithdrawFromCommitment(ctx, &commitments, liquidCoins, true, addr)
	if err != nil {
		return err
	}

	// Emit Hook if Eden is uncommitted
	if denom == ptypes.Eden {
		err = k.EdenUncommitted(ctx, addr.String(), sdk.NewCoin(denom, amount))
		if err != nil {
			return err
		}
	}

	// Emit Hook commitment changed
	err = k.CommitmentChanged(ctx, addr.String(), sdk.Coins{sdk.NewCoin(denom, amount)})
	if err != nil {
		return err
	}

	// Update total commitment
	params := k.GetParams(ctx)
	params.TotalCommitted = params.TotalCommitted.Add(liquidCoins...)
	k.SetParams(ctx, params)

	// Emit blockchain event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCommitmentChanged,
			sdk.NewAttribute(types.AttributeCreator, addr.String()),
			sdk.NewAttribute(types.AttributeAmount, amount.String()),
			sdk.NewAttribute(types.AttributeDenom, denom),
		),
	)
	return nil
}

// UncommitTokens uncommits the tokens from committed store and make it liquid immediately.
func (k msgServer) UncommitTokens(goCtx context.Context, msg *types.MsgUncommitTokens) (*types.MsgUncommitTokensResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	addr, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, "unable to convert address from bech32")
	}

	if msg.Denom != ptypes.Eden && msg.Denom != ptypes.EdenB {
		return nil, types.ErrUnsupportedUncommitToken
	}

	err = k.Keeper.UncommitTokens(ctx, addr, msg.Denom, msg.Amount)
	if err != nil {
		return nil, err
	}

	return &types.MsgUncommitTokensResponse{}, nil
}
