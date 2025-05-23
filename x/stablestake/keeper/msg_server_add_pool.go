package keeper

import (
	"context"
	"errors"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v5/x/stablestake/types"
)

func (k msgServer) AddPool(goCtx context.Context, msg *types.MsgAddPool) (*types.MsgAddPoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	params := k.ammKeeper.GetParams(ctx)

	if !params.IsCreatorAllowed(msg.Sender) {
		return nil, errors.New("sender is not allowed to create pool")
	}

	found := k.HasPoolByDenom(ctx, msg.DepositDenom)
	if found {
		return nil, errorsmod.Wrapf(types.ErrPoolAlreadyExists, "pool with denom %s already exists", msg.DepositDenom)
	}

	poolId := k.GetNextPoolId(ctx)
	pool := types.Pool{
		Id:                   poolId,
		DepositDenom:         msg.DepositDenom,
		InterestRateDecrease: msg.InterestRateDecrease,
		InterestRateIncrease: msg.InterestRateIncrease,
		HealthGainFactor:     msg.HealthGainFactor,
		MaxLeverageRatio:     msg.MaxLeverageRatio,
		MaxWithdrawRatio:     msg.MaxWithdrawRatio,
		InterestRateMax:      msg.InterestRateMax,
		InterestRateMin:      msg.InterestRateMin,
		InterestRate:         msg.InterestRate,
		NetAmount:            math.ZeroInt(),
	}

	k.SetPool(ctx, pool)

	return &types.MsgAddPoolResponse{
		PoolId: poolId,
	}, nil
}
