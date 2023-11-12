package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/x/commitment/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

func (k msgServer) ClaimRewards(goCtx context.Context, msg *types.MsgClaimRewards) (*types.MsgClaimRewardsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get the Commitments for the creator
	commitments, found := k.GetCommitments(ctx, msg.Creator)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrCommitmentsNotFound, "creator: %s", msg.Creator)
	}

	unclaimed := commitments.RewardsUnclaimed
	commitments.RewardsUnclaimed = sdk.Coins{}

	edenAmount := unclaimed.AmountOf(ptypes.Eden)
	commitments.AddClaimed(sdk.NewCoin(ptypes.Eden, edenAmount))
	k.SetCommitments(ctx, commitments)

	withdrawCoins := unclaimed.Sub(sdk.NewCoin(ptypes.Eden, edenAmount))

	// Mint the withdrawn tokens to the module account
	err := k.bankKeeper.MintCoins(ctx, types.ModuleName, withdrawCoins)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInsufficientFunds, "unable to mint withdrawn tokens")
	}

	addr, err := sdk.AccAddressFromBech32(commitments.Creator)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "unable to convert address from bech32")
	}

	// Send the minted coins to the user's account
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, withdrawCoins)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInsufficientFunds, "unable to send withdrawn tokens")
	}

	// Emit Hook commitment changed
	k.AfterCommitmentChange(ctx, msg.Creator, withdrawCoins)

	return &types.MsgClaimRewardsResponse{}, nil
}
