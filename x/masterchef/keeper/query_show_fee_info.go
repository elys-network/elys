package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v4/x/masterchef/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) ShowFeeInfo(goCtx context.Context, req *types.QueryShowFeeInfoRequest) (*types.QueryShowFeeInfoResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	return &types.QueryShowFeeInfoResponse{FeeInfo: k.GetFeeInfo(ctx, req.Date)}, nil
}
