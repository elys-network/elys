package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/leveragelp/types"
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

	updatedLeveragePositions := []*types.QueryPosition{}
	for i ,position := range positions {
		updated_leverage :=  position.LeveragedLpAmount.Quo(position.LeveragedLpAmount.Sub(position.Liabilities))
		updatedLeveragePositions[i] = &types.QueryPosition{
			Position: position,
			UpdatedLeverage: updated_leverage,
		}
	}

	return &types.PositionsByPoolResponse{
		Positions:  updatedLeveragePositions,
		Pagination: pageRes,
	}, nil
}
