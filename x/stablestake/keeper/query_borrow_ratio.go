package keeper

import (
	"context"

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

	params := k.GetParams(ctx)
	moduleAddr := authtypes.NewModuleAddress(types.ModuleName)

	depositDenom := k.GetDepositDenom(ctx)

	balance := k.bk.GetBalance(ctx, moduleAddr, depositDenom)
	borrowed := params.TotalValue.Sub(balance.Amount)
	borrowRatio := sdk.ZeroDec()
	if params.TotalValue.GT(sdk.ZeroInt()) {
		borrowRatio = sdk.NewDecFromInt(borrowed).
			Quo(sdk.NewDecFromInt(params.TotalValue))
	}

	return &types.QueryBorrowRatioResponse{
		TotalDeposit: params.TotalValue,
		TotalBorrow:  borrowed,
		BorrowRatio:  borrowRatio,
	}, nil
}
