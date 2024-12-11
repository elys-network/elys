package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/elys-network/elys/x/commitment/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

const airdropWallet = ""

func (k msgServer) ClaimAirdrop(goCtx context.Context, msg *types.MsgClaimAirdrop) (*types.MsgClaimAirdropResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	sender := sdk.MustAccAddressFromBech32(msg.Creator)
	airdropWallet := sdk.MustAccAddressFromBech32(airdropWallet)
	params := k.GetParams(ctx)

	if k.GetAirdropClaimed(ctx, sender).Claimed {
		return nil, types.ErrAirdropAlreadyClaimed
	}

	currentHeight := uint64(ctx.BlockHeight())
	if currentHeight < params.StartAtomStakersHeight {
		return nil, types.ErrAirdropNotStarted
	}

	if currentHeight > params.EndAtomStakersHeight {
		return nil, types.ErrAirdropEnded
	}

	// Elys to be sent
	total_elys := k.GetNFTHolder(ctx, sender).Amount
	total_elys = total_elys.Add(k.GetCadet(ctx, sender).Amount)
	total_elys = total_elys.Add(k.GetGovernor(ctx, sender).Amount)
	// TODO: Update airdrop wallet address
	err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, airdropWallet, types.ModuleName, sdk.NewCoins(sdk.NewCoin(ptypes.Elys, total_elys)))
	if err != nil {
		return nil, err
	}
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sender, sdk.NewCoins(sdk.NewCoin(ptypes.Elys, total_elys)))
	if err != nil {
		return nil, err
	}

	// Eden to be sent
	total_eden := k.GetAtomStaker(ctx, sender).Amount
	// Add eden to commitment
	err = k.MintCoins(ctx, types.ModuleName, sdk.NewCoins(sdk.NewCoin(ptypes.Eden, total_eden)))
	if err != nil {
		return nil, err
	}
	err = k.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sender, sdk.NewCoins(sdk.NewCoin(ptypes.Eden, total_eden)))
	if err != nil {
		return nil, err
	}

	k.SetAirdropClaimed(ctx, sender)

	// Update tracking of total claimed
	params.TotalEdenClaimed = params.TotalEdenClaimed.Add(total_eden)
	params.TotalElysClaimed = params.TotalElysClaimed.Add(total_elys)
	k.SetParams(ctx, params)

	return &types.MsgClaimAirdropResponse{
		ElysAmount: total_elys,
		EdenAmount: total_eden,
	}, nil
}
