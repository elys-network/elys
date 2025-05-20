package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v4/x/leveragelp/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) QueryPositions(goCtx context.Context, req *types.PositionsRequest) (*types.PositionsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.Pagination != nil && req.Pagination.Limit > types.MaxPageLimit {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("page size greater than max %d", types.MaxPageLimit))
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	positions, page, err := k.GetPositions(ctx, req.Pagination)
	if err != nil {
		return nil, err
	}
	updatedLeveragePositions, err := k.GetLeverageLpUpdatedLeverage(ctx, positions)
	if err != nil {
		return nil, err
	}

	return &types.PositionsResponse{
		Positions:  updatedLeveragePositions,
		Pagination: page,
	}, nil
}
