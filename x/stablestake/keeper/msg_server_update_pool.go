package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/elys-network/elys/v7/x/stablestake/types"
)

func (k msgServer) UpdatePool(goCtx context.Context, msg *types.MsgUpdatePool) (*types.MsgUpdatePoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if k.authority != msg.Authority {
		return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
	}

	pool, found := k.GetPool(ctx, msg.PoolId)
	if !found {
		return nil, errorsmod.Wrap(types.ErrPoolNotFound, "pool does not exist")
	}

	pool.InterestRateDecrease = msg.InterestRateDecrease
	pool.InterestRateIncrease = msg.InterestRateIncrease
	pool.HealthGainFactor = msg.HealthGainFactor
	pool.MaxLeverageRatio = msg.MaxLeverageRatio
	pool.InterestRateMax = msg.InterestRateMax
	pool.InterestRateMin = msg.InterestRateMin

	k.SetPool(ctx, pool)

	return &types.MsgUpdatePoolResponse{}, nil
}
