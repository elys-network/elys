package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/x/commitment/types"
)

// VestNow is not enabled at this stage
var VestNowEnabled = false

// VestNow provides functionality to get the token immediately but lower amount than original
// e.g. user can burn 1000 ueden and get 800 uelys when the ratio is 80%
func (k msgServer) VestNow(goCtx context.Context, msg *types.MsgVestNow) (*types.MsgVestNowResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if !VestNowEnabled {
		return nil, types.ErrVestNowIsNotEnabled
	}

	vestingInfo, _ := k.GetVestingInfo(ctx, msg.Denom)
	if vestingInfo == nil {
		return nil, sdkerrors.Wrapf(types.ErrInvalidDenom, "denom: %s", msg.Denom)
	}

	commitments, err := k.DeductCommitments(ctx, msg.Creator, msg.Denom, msg.Amount)
	if err != nil {
		return nil, err
	}

	vestAmount := msg.Amount.Quo(vestingInfo.VestNowFactor)
	withdrawCoins := sdk.NewCoins(sdk.NewCoin(vestingInfo.VestingDenom, vestAmount))

	// Mint the vested tokens to the module account
	err = k.bankKeeper.MintCoins(ctx, types.ModuleName, withdrawCoins)
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

	// Update the commitments
	k.SetCommitments(ctx, commitments)

	// Emit Hook commitment changed
	k.AfterCommitmentChange(ctx, msg.Creator, sdk.Coins{sdk.NewCoin(msg.Denom, msg.Amount)})

	// Emit blockchain event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCommitmentChanged,
			sdk.NewAttribute(types.AttributeCreator, msg.Creator),
			sdk.NewAttribute(types.AttributeAmount, msg.Amount.String()),
			sdk.NewAttribute(types.AttributeDenom, msg.Denom),
		),
	)

	return &types.MsgVestNowResponse{}, nil
}
