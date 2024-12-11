package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/elys-network/elys/x/commitment/types"
)

func (k msgServer) ClaimAirdrop(goCtx context.Context, msg *types.MsgClaimAirdrop) (*types.MsgClaimAirdropResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Elys to be sent
	total_elys := k.GetNFTHolder(ctx, sdk.MustAccAddressFromBech32(msg.Creator)).Amount
	total_elys = total_elys.Add(k.GetCadet(ctx, sdk.MustAccAddressFromBech32(msg.Creator)).Amount)
	total_elys = total_elys.Add(k.GetGovernor(ctx, sdk.MustAccAddressFromBech32(msg.Creator)).Amount)

	// Eden to be sent
	total_eden := k.GetAtomStaker(ctx, sdk.MustAccAddressFromBech32(msg.Creator)).Amount

	return &types.MsgClaimAirdropResponse{
		ElysAmount: total_elys,
		EdenAmount: total_eden,
	}, nil
}
