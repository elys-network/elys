package keeper

import (
	"context"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/elys-network/elys/v5/x/commitment/types"
	ptypes "github.com/elys-network/elys/v5/x/parameter/types"
)

var maxEdenAmountToClaim = sdkmath.NewInt(1000000000000)

func (k msgServer) ClaimRewardProgram(goCtx context.Context, msg *types.MsgClaimRewardProgram) (*types.MsgClaimRewardProgramResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	sender := sdk.MustAccAddressFromBech32(msg.ClaimAddress)
	rewardProgram := k.GetRewardProgram(ctx, sender)

	if rewardProgram.Amount.IsZero() {
		return nil, types.ErrRewardProgramNotFound
	}

	edenAmount := rewardProgram.Amount

	params := k.GetParams(ctx)
	if !params.EnableClaim {
		return nil, types.ErrClaimNotEnabled
	}

	if rewardProgram.Claimed {
		return nil, types.ErrRewardProgramAlreadyClaimed
	}

	currentHeight := uint64(ctx.BlockHeight())
	if currentHeight < params.StartRewardProgramClaimHeight {
		return nil, types.ErrRewardProgramNotStarted
	}

	if currentHeight > params.EndRewardProgramClaimHeight {
		return nil, types.ErrRewardProgramEnded
	}

	total := k.GetTotalRewardProgramClaimed(ctx)
	total.TotalEdenClaimed = total.TotalEdenClaimed.Add(edenAmount)
	k.SetTotalRewardProgramClaimed(ctx, total)

	// Add eden to commitment
	err := k.MintCoins(ctx, types.ModuleName, sdk.NewCoins(sdk.NewCoin(ptypes.Eden, edenAmount)))
	if err != nil {
		return nil, err
	}
	err = k.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sender, sdk.NewCoins(sdk.NewCoin(ptypes.Eden, edenAmount)))
	if err != nil {
		return nil, err
	}

	rewardProgram.Claimed = true
	k.SetRewardProgram(ctx, rewardProgram)

	// This will never be triggered
	if total.TotalEdenClaimed.GT(maxEdenAmountToClaim) {
		return nil, types.ErrMaxEdenAmountReached
	}

	// Emit event for reward program claim
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeClaimRewardProgram,
			sdk.NewAttribute(types.AttributeKeyClaimAddress, msg.ClaimAddress),
			sdk.NewAttribute(types.AttributeKeyEdenAmount, edenAmount.String()),
			sdk.NewAttribute(types.AttributeKeyTotalEdenClaimed, total.TotalEdenClaimed.String()),
		),
	)

	return &types.MsgClaimRewardProgramResponse{
		EdenAmount: edenAmount,
	}, nil
}
