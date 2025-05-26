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

	if k.GetRewardProgramClaimed(ctx, sender).Claimed {
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

	k.SetRewardProgramClaimed(ctx, sender)

	// This will never be triggered
	if total.TotalEdenClaimed.GT(maxEdenAmountToClaim) {
		return nil, types.ErrMaxEdenAmountReached
	}

	return &types.MsgClaimRewardProgramResponse{
		EdenAmount: edenAmount,
	}, nil
}
