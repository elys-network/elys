package keeper

import (
	"context"

	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/elys-network/elys/x/stablestake/types"
	"github.com/osmosis-labs/osmosis/osmomath"
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
	borrowed := pool.TotalValue.Sub(balance.Amount)
	borrowRatio := osmomath.ZeroBigDec()
	if pool.TotalValue.GT(sdkmath.ZeroInt()) {
		borrowRatio = osmomath.BigDecFromSDKInt(borrowed).Quo(pool.GetBigDecTotalValue())
	}

	return &types.QueryBorrowRatioResponse{
		TotalDeposit: pool.TotalValue,
		TotalBorrow:  borrowed,
		BorrowRatio:  borrowRatio.Dec(),
	}, nil
}
