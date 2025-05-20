package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v4/x/leveragelp/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) QueryPositionsByPool(goCtx context.Context, req *types.PositionsByPoolRequest) (*types.PositionsByPoolResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	positions, pageRes, err := k.GetPositionsForPool(ctx, req.AmmPoolId, req.Pagination)
	if err != nil {
		return nil, err
	}

	updatedLeveragePositions, err := k.GetLeverageLpUpdatedLeverage(ctx, positions)
	if err != nil {
		return nil, err
	}

	return &types.PositionsByPoolResponse{
		Positions:  updatedLeveragePositions,
		Pagination: pageRes,
	}, nil
}
