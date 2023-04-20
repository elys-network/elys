package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	aptypes "github.com/elys-network/elys/x/assetprofile/types"
	"github.com/elys-network/elys/x/commitment/types"
)

func (k msgServer) DepositTokens(goCtx context.Context, msg *types.MsgDepositTokens) (*types.MsgDepositTokensResponse, error) {
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
			Creator:           msg.Creator,
			CommittedTokens:   []*types.CommittedTokens{},
			UncommittedTokens: []*types.UncommittedTokens{},
		}
	}
	// Get the uncommitted tokens for the creator
	uncommittedToken, _ := commitments.GetUncommittedTokensForDenom(msg.Denom)
	if !found {
		uncommittedTokens := commitments.GetUncommittedTokens()
		uncommittedToken = &types.UncommittedTokens{
			Denom:  msg.Denom,
			Amount: sdk.ZeroInt(),
		}
		uncommittedTokens = append(uncommittedTokens, uncommittedToken)
		commitments.UncommittedTokens = uncommittedTokens
	}
	// Update the uncommitted tokens amount
	uncommittedToken.Amount = uncommittedToken.Amount.Add(msg.Amount)

	// Update the commitments
	k.SetCommitments(ctx, commitments)

	// Emit Hook commitment changed
	k.HookCommitmentChanged(ctx, msg.Creator, sdk.NewCoin(msg.Denom, msg.Amount))

	// Emit blockchain event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCommitmentChanged,
			sdk.NewAttribute(types.AttributeCreator, msg.Creator),
			sdk.NewAttribute(types.AttributeAmount, msg.Amount.String()),
			sdk.NewAttribute(types.AttributeDenom, msg.Denom),
		),
	)

	return &types.MsgDepositTokensResponse{}, nil
}
