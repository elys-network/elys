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

func (k Keeper) BorrowRatio(goCtx context.Context, req *types.QueryBorrowRatioRequest) (*types.QueryBorrowRatioResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	pool, found := k.GetPool(ctx, req.PoolId)
	if !found {
		return nil, types.ErrPoolNotFound
	}

	moduleAddr := authtypes.NewModuleAddress(types.ModuleName)

	depositDenom := pool.GetDepositDenom()

	balance := k.bk.GetBalance(ctx, moduleAddr, depositDenom)
	borrowed := pool.NetAmount.Sub(balance.Amount)
	borrowRatio := sdkmath.LegacyZeroDec()
	if pool.NetAmount.GT(sdkmath.ZeroInt()) {
		borrowRatio = borrowed.ToLegacyDec().Quo(pool.NetAmount.ToLegacyDec())
	}

	return &types.QueryBorrowRatioResponse{
		NetAmount:   pool.NetAmount,
		TotalBorrow: borrowed,
		BorrowRatio: borrowRatio,
	}, nil
}
