package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/x/commitment/types"
)

func (k msgServer) VestNow(goCtx context.Context, msg *types.MsgVestNow) (*types.MsgVestNowResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	vestingInfo := k.GetVestingInfo(ctx, msg.Denom)

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

	return &types.MsgVestNowResponse{}, nil
}
