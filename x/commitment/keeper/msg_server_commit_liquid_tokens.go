package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	aptypes "github.com/elys-network/elys/x/assetprofile/types"
	"github.com/elys-network/elys/x/commitment/types"
)

// CommitLiquidTokens commit the tokens from user's balance
func (k msgServer) CommitLiquidTokens(goCtx context.Context, msg *types.MsgCommitLiquidTokens) (*types.MsgCommitLiquidTokensResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	assetProfile, found := k.apKeeper.GetEntry(ctx, msg.Denom)
	if !found {
		return nil, sdkerrors.Wrapf(aptypes.ErrAssetProfileNotFound, "denom: %s", msg.Denom)
	}

	if !assetProfile.CommitEnabled {
		return nil, sdkerrors.Wrapf(types.ErrCommitDisabled, "denom: %s", msg.Denom)
	}

	depositCoins := sdk.NewCoins(sdk.NewCoin(msg.Denom, msg.Amount))

	addr, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "unable to convert address from bech32")
	}

	// send the deposited coins to the module
	err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, addr, types.ModuleName, depositCoins)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInsufficientFunds, fmt.Sprintf("unable to send deposit tokens: %v", depositCoins))
	}
	// burn the deposited coins
	err = k.bankKeeper.BurnCoins(ctx, types.ModuleName, depositCoins)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInsufficientFunds, "unable to burn deposit tokens")
	}

	// Get the Commitments for the creator
	commitments, found := k.GetCommitments(ctx, msg.Creator)
	if !found {
		commitments = types.Commitments{
			Creator:          msg.Creator,
			CommittedTokens:  []*types.CommittedTokens{},
			RewardsUnclaimed: sdk.Coins{},
		}
	}

	// Update the commitments
	commitments.AddCommittedTokens(msg.Denom, msg.Amount, uint64(ctx.BlockTime().Unix())+msg.MinLock)
	k.SetCommitments(ctx, commitments)

	// Emit Hook commitment changed
	k.AfterCommitmentChange(ctx, msg.Creator, sdk.NewCoin(msg.Denom, msg.Amount))

	// Emit blockchain event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCommitmentChanged,
			sdk.NewAttribute(types.AttributeCreator, msg.Creator),
			sdk.NewAttribute(types.AttributeAmount, msg.Amount.String()),
			sdk.NewAttribute(types.AttributeDenom, msg.Denom),
		),
	)

	return &types.MsgCommitLiquidTokensResponse{}, nil
}
