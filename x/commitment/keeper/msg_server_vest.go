package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/elys-network/elys/x/commitment/types"
)

func (k msgServer) Vest(goCtx context.Context, msg *types.MsgVest) (*types.MsgVestResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	vestingInfo := k.GetVestingInfo(ctx, msg.Denom)

	if vestingInfo == nil {
		return nil, sdkerrors.Wrapf(types.ErrInvalidDenom, "denom: %s", msg.Denom)
	}

	commitments, err := k.DeductCommitments(ctx, msg.Creator, msg.Denom, msg.Amount)
	if err != nil {
		return nil, err
	}

	// Create vesting tokens entry and add to commitments
	vestingTokens := commitments.GetVestingTokens()
	vestingTokens = append(vestingTokens, &types.VestingTokens{
		Denom:           vestingInfo.VestingDenom,
		TotalAmount:     msg.Amount,
		UnvestedAmount:  msg.Amount,
		EpochIdentifier: vestingInfo.EpochIdentifier,
		NumEpochs:       vestingInfo.NumEpochs,
		CurrentEpoch:    0,
	})
	commitments.VestingTokens = vestingTokens

	// Update the commitments
	k.SetCommitments(ctx, commitments)

	return &types.MsgVestResponse{}, nil
}
