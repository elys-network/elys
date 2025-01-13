package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/elys-network/elys/x/stablestake/types"
)

func (k msgServer) AddPool(goCtx context.Context, msg *types.MsgAddPool) (*types.MsgAddPoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if k.authority != msg.Authority {
		return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
	}

	found := k.HasPoolByDenom(ctx, msg.DepositDenom)
	if found {
		return nil, errorsmod.Wrapf(types.ErrPoolAlreadyExists, "pool with denom %s already exists", msg.DepositDenom)
	}

	poolId := k.GetNextPoolId(ctx)
	pool := types.Pool{
		PoolId:               poolId,
		DepositDenom:         msg.DepositDenom,
		RedemptionRate:       math.LegacyZeroDec(),
		InterestRateDecrease: msg.InterestRateDecrease,
		InterestRateIncrease: msg.InterestRateIncrease,
		HealthGainFactor:     msg.HealthGainFactor,
		MaxLeverageRatio:     msg.MaxLeverageRatio,
		InterestRateMax:      msg.InterestRateMax,
		InterestRateMin:      msg.InterestRateMin,
		InterestRate:         msg.InterestRate,
		TotalValue:           math.ZeroInt(),
	}

	k.SetPool(ctx, pool)

	return &types.MsgAddPoolResponse{
		PoolId: poolId,
	}, nil
}
