package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v6/x/perpetual/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) GetPositionsByPool(goCtx context.Context, req *types.PositionsByPoolRequest) (*types.PositionsByPoolResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	mtps, pageRes, err := k.GetMTPsForPool(ctx, req.AmmPoolId, req.Pagination)
	if err != nil {
		return nil, err
	}

	return &types.PositionsByPoolResponse{
		Mtps:       mtps,
		Pagination: pageRes,
	}, nil
}
