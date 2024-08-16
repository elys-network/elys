package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/x/commitment/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
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
		return nil, errorsmod.Wrapf(types.ErrInvalidDenom, "denom: %s", msg.Denom)
	}

	creator := sdk.MustAccAddressFromBech32(msg.Creator)
	commitments, err := k.DeductClaimed(ctx, creator, msg.Denom, msg.Amount)
	if err != nil {
		return nil, err
	}

	// Ensure vestingInfo.VestNowFactor is not zero to avoid division by zero
	if vestingInfo.VestNowFactor.IsZero() {
		return nil, types.ErrInvalidAmount
	}

	vestAmount := msg.Amount.Quo(vestingInfo.VestNowFactor)
	withdrawCoins := sdk.NewCoins(sdk.NewCoin(vestingInfo.VestingDenom, vestAmount))

	addr, err := sdk.AccAddressFromBech32(commitments.Creator)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, "unable to convert address from bech32")
	}

	// mint coins if vesting token is ELYS
	if vestingInfo.VestingDenom == ptypes.Elys {
		err := k.bankKeeper.MintCoins(ctx, types.ModuleName, withdrawCoins)
		if err != nil {
			ctx.Logger().Debug("unable to mint vested tokens for ELYS token")
			return nil, err
		}
	}

	// Send the minted coins to the user's account
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, withdrawCoins)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInsufficientFunds, "unable to send withdrawn tokens")
	}

	// Update the commitments
	k.SetCommitments(ctx, commitments)

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
