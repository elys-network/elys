package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/elys-network/elys/x/masterchef/types"
)

func (k msgServer) TogglePoolEdenRewards(goCtx context.Context, msg *types.MsgTogglePoolEdenRewards) (*types.MsgTogglePoolEdenRewardsResponse, error) {
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	if k.authority != msg.Authority {
		return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
	}

	pool, found := k.GetPoolInfo(ctx, msg.PoolId)
	if !found {
		return &types.MsgTogglePoolEdenRewardsResponse{}, types.ErrPoolNotFound
	}

	pool.EnableEdenRewards = msg.Enable
	k.SetPoolInfo(ctx, pool)
	return &types.MsgTogglePoolEdenRewardsResponse{}, nil
}
