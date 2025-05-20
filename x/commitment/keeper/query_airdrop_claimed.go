package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v4/x/commitment/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) TotalAirdropClaimed(goCtx context.Context, req *types.QueryTotalAirDropClaimedRequest) (*types.QueryTotalAirDropClaimedResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	claimed := k.GetTotalClaimed(ctx)
	return &types.QueryTotalAirDropClaimedResponse{
		TotalElysClaimed: claimed.TotalElysClaimed,
		TotalEdenClaimed: claimed.TotalEdenClaimed,
	}, nil
}
