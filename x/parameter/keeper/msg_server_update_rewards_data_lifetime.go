package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/elys-network/elys/v7/x/parameter/types"
)

func (k msgServer) UpdateRewardsDataLifetime(goCtx context.Context, msg *types.MsgUpdateRewardsDataLifetime) (*types.MsgUpdateRewardsDataLifetimeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if k.authority != msg.Creator {
		return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Creator)
	}

	params := k.GetParams(ctx)
	params.RewardsDataLifetime = msg.RewardsDataLifetime
	k.SetParams(ctx, params)

	return &types.MsgUpdateRewardsDataLifetimeResponse{}, nil
}
