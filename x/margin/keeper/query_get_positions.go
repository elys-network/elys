package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/margin/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) GetPositions(goCtx context.Context, req *types.PositionsRequest) (*types.PositionsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.Pagination.Limit > types.MaxPageLimit {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("page size greater than max %d", types.MaxPageLimit))
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	mtps, page, err := k.GetMTPs(ctx, req.Pagination)
	if err != nil {
		return nil, err
	}

	return &types.PositionsResponse{
		Mtps:       mtps,
		Pagination: page,
	}, nil
}
