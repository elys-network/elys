package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v6/x/tier/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) Staked(goCtx context.Context, req *types.QueryStakedRequest) (*types.QueryStakedResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	sender := sdk.MustAccAddressFromBech32(req.User)
	ctx := sdk.UnwrapSDKContext(goCtx)

	com, del, unbon, totalVested := k.RetrieveStaked(ctx, sender)

	return &types.QueryStakedResponse{
		Commitments: com.Dec(),
		Delegations: del.Dec(),
		Unbondings:  unbon.Dec(),
		TotalVested: totalVested.Dec(),
	}, nil
}
