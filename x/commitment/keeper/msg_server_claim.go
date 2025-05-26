package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/elys-network/elys/v5/x/commitment/types"
	ptypes "github.com/elys-network/elys/v5/x/parameter/types"
)

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
		return nil, types.ErrAirdropAlreadyClaimed
	}

	currentHeight := uint64(ctx.BlockHeight())
	if currentHeight < params.StartAirdropClaimHeight {
		return nil, types.ErrAirdropNotStarted
	}

	if currentHeight > params.EndAirdropClaimHeight {
		return nil, types.ErrAirdropEnded
	}

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

	return &types.MsgClaimRewardProgramResponse{
		EdenAmount: edenAmount,
	}, nil
}
