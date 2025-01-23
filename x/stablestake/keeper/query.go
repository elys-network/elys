package keeper

import (
	"context"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/elys-network/elys/x/stablestake/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) Pool(goCtx context.Context, req *types.QueryGetPoolRequest) (*types.QueryGetPoolResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	pool, found := k.GetPool(ctx, req.PoolId)
	if !found {
		return nil, types.ErrPoolNotFound
	}

	poolInfo := k.GetPoolInfo(ctx, pool)

	return &types.QueryGetPoolResponse{
		Pool: poolInfo,
	}, nil
}

func (k Keeper) Pools(goCtx context.Context, req *types.QueryAllPoolRequest) (*types.QueryAllPoolResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	var res []types.PoolResponse
	pools := k.GetAllPools(ctx)
	for _, pool := range pools {
		poolInfo := k.GetPoolInfo(ctx, pool)
		res = append(res, poolInfo)
	}

	return &types.QueryAllPoolResponse{Pools: res}, nil
}

func (k Keeper) GetPoolInfo(ctx sdk.Context, pool types.Pool) types.PoolResponse {
	moduleAddr := authtypes.NewModuleAddress(types.ModuleName)
	depositDenom := pool.GetDepositDenom()

	balance := k.bk.GetBalance(ctx, moduleAddr, depositDenom)
	borrowed := pool.TotalValue.Sub(balance.Amount)
	borrowRatio := sdkmath.LegacyZeroDec()
	if pool.TotalValue.GT(sdkmath.ZeroInt()) {
		borrowRatio = borrowed.ToLegacyDec().Quo(pool.TotalValue.ToLegacyDec())
	}

	return types.PoolResponse{
		DepositDenom:         pool.DepositDenom,
		RedemptionRate:       pool.RedemptionRate,
		InterestRate:         pool.InterestRate,
		InterestRateMax:      pool.InterestRateMax,
		InterestRateMin:      pool.InterestRateMin,
		InterestRateIncrease: pool.InterestRateIncrease,
		InterestRateDecrease: pool.InterestRateDecrease,
		HealthGainFactor:     pool.HealthGainFactor,
		TotalValue:           pool.TotalValue,
		MaxLeverageRatio:     pool.MaxLeverageRatio,
		PoolId:               pool.PoolId,
		TotalDeposit:         pool.TotalValue,
		TotalBorrow:          borrowed,
		BorrowRatio:          borrowRatio,
	}
}
