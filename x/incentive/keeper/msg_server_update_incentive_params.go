package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/elys-network/elys/x/incentive/types"
)

func (k msgServer) UpdateIncentiveParams(goCtx context.Context, msg *types.MsgUpdateIncentiveParams) (*types.MsgUpdateIncentiveParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if k.authority != msg.Authority {
		return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
	}

	params := k.GetParams(ctx)
	params.RewardPortionForLps = msg.RewardPortionForLps
	params.ElysStakeSnapInterval = msg.ElysStakeSnapInterval
	params.MaxEdenRewardAprLps = msg.MaxEdenRewardAprLps
	params.MaxEdenRewardAprStakers = msg.MaxEdenRewardAprStakers
	params.DistributionInterval = msg.DistributionInterval

	k.SetParams(ctx, params)

	return &types.MsgUpdateIncentiveParamsResponse{}, nil
}
