package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/tier/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) QueryTotalBorrows(goCtx context.Context, req *types.QueryTotalBorrowsRequest) (*types.QueryTotalBorrowsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	sender := sdk.MustAccAddressFromBech32(req.User)
	total := k.RetrieveTotalBorrows(ctx, sender)

	return &types.QueryTotalBorrowsResponse{
		Total: total,
	}, nil
}


