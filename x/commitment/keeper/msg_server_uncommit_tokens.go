package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	assetprofiletypes "github.com/elys-network/elys/v5/x/assetprofile/types"
	"github.com/elys-network/elys/v5/x/commitment/types"
	ptypes "github.com/elys-network/elys/v5/x/parameter/types"
)

func (k Keeper) UncommitTokens(ctx sdk.Context, addr sdk.AccAddress, denom string, amount math.Int, isLiquidation bool) error {
	assetProfile, found := k.assetProfileKeeper.GetEntry(ctx, denom)
	if !found {
		return errorsmod.Wrapf(assetprofiletypes.ErrAssetProfileNotFound, "denom: %s, address: %s", denom, addr.String())
	}

	if !assetProfile.WithdrawEnabled {
		return errorsmod.Wrapf(types.ErrCommitDisabled, "denom: %s, address: %s", denom, addr.String())
	}

	// Get the Commitments for the creator
	commitments := k.GetCommitments(ctx, addr)

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
	err := commitments.DeductFromCommitted(denom, amount, uint64(ctx.BlockTime().Unix()), isLiquidation)
	if err != nil {
		return err
	}
	k.SetCommitments(ctx, commitments)

	liquidCoins := sdk.NewCoins(sdk.NewCoin(denom, amount))
	edenAmount := liquidCoins.AmountOf(ptypes.Eden)
	edenBAmount := liquidCoins.AmountOf(ptypes.EdenB)
	commitments.AddClaimed(sdk.NewCoin(ptypes.Eden, edenAmount))
	commitments.AddClaimed(sdk.NewCoin(ptypes.EdenB, edenBAmount))
	k.SetCommitments(ctx, commitments)

	// Emit Hook commitment changed
	err = k.CommitmentChanged(ctx, addr, sdk.Coins{sdk.NewCoin(denom, amount)})
	if err != nil {
		return err
	}

	withdrawCoins := liquidCoins.
		Sub(sdk.NewCoin(ptypes.Eden, edenAmount)).
		Sub(sdk.NewCoin(ptypes.EdenB, edenBAmount))

	if !withdrawCoins.Empty() {
		err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, withdrawCoins)
		if err != nil {
			return err
		}
	}

	// Update total commitment
	params := k.GetParams(ctx)
	params.TotalCommitted = params.TotalCommitted.Sub(liquidCoins...)
	k.SetParams(ctx, params)

	// Emit Hook if Eden is uncommitted
	if denom == ptypes.Eden {
		err = k.EdenUncommitted(ctx, addr, sdk.NewCoin(denom, amount))
		if err != nil {
			return err
		}
	}

	// Emit event
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

	err = k.Keeper.UncommitTokens(ctx, addr, msg.Denom, msg.Amount, false)
	if err != nil {
		return nil, err
	}

	return &types.MsgUncommitTokensResponse{}, nil
}
