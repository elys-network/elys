package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/elys-network/elys/x/incentive/types"
)

func (k msgServer) UpdateIncentiveParams(goCtx context.Context, msg *types.MsgUpdateIncentiveParams) (*types.MsgUpdateIncentiveParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if k.authority != msg.Authority {
		return nil, sdkerrors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
	}

	params := k.GetParams(ctx)
	params.CommunityTax = msg.CommunityTax
	params.WithdrawAddrEnabled = msg.WithdrawAddrEnabled
	params.RewardPortionForLps = msg.RewardPortionForLps
	params.ElysStakeTrackingRate = msg.ElysStakeTrackingRate
	params.MaxEdenRewardAprLps = msg.MaxEdenRewardAprLps
	params.MaxEdenRewardAprStakers = msg.MaxEdenRewardAprStakers
	params.DistributionEpochForLpsInBlocks = msg.DistributionEpochForLps
	params.DistributionEpochForStakersInBlocks = msg.DistributionEpochForStakers

	k.SetParams(ctx, params)

	return &types.MsgUpdateIncentiveParamsResponse{}, nil
}
