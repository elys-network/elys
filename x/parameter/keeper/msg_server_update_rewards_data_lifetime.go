package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/elys-network/elys/x/parameter/types"
)

func (k msgServer) UpdateRewardsDataLifetime(goCtx context.Context, msg *types.MsgUpdateRewardsDataLifetime) (*types.MsgUpdateRewardsDataLifetimeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if k.authority != msg.Creator {
		return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Creator)
	}

	rewardsDataLifetime, ok := sdk.NewIntFromString(msg.RewardsDataLifetime)
	if !ok {
		return nil, errorsmod.Wrapf(types.ErrInvalidRewardsDataLifecycle, "invalid data in rewards_data_lifecycle")
	}

	if !rewardsDataLifetime.IsPositive() {
		return nil, errorsmod.Wrapf(types.ErrInvalidRewardsDataLifecycle, "rewards_data_lifecycle must be positive")
	}

	params := k.GetParams(ctx)
	params.RewardsDataLifetime = rewardsDataLifetime.Int64()
	k.SetParams(ctx, params)

	return &types.MsgUpdateRewardsDataLifetimeResponse{}, nil
}
