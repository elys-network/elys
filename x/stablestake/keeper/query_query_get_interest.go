package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/stablestake/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) QueryGetInterest(goCtx context.Context, req *types.QueryQueryGetInterestRequest) (*types.QueryQueryGetInterestResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	debt := k.getDebt(ctx, sdk.MustAccAddressFromBech32(req.Owner))
	res := k.GetInterest(ctx, debt.LastInterestCalcBlock, debt.LastInterestCalcTime, debt.Borrowed.ToLegacyDec())

	return &types.QueryQueryGetInterestResponse{
		Interest: res,
	}, nil
}
