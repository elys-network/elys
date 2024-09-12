package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/stablestake/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) QueryGetDebt(goCtx context.Context, req *types.QueryQueryGetDebtRequest) (*types.QueryQueryGetDebtResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	debt := k.getDebt(ctx, sdk.MustAccAddressFromBech32(req.Owner))
	interest_debt := k.GetDebt(ctx, sdk.MustAccAddressFromBech32(req.Owner))

	return &types.QueryQueryGetDebtResponse{
		Debt:         debt,
		DebtInterest: interest_debt,
	}, nil
}
