package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v6/x/tier/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) Perpetual(goCtx context.Context, req *types.QueryPerpetualRequest) (*types.QueryPerpetualResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	sender := sdk.MustAccAddressFromBech32(req.User)
	total, borrow, _ := k.RetrievePerpetualTotal(ctx, sender)

	return &types.QueryPerpetualResponse{
		TotalValue:   total.Dec(),
		TotalBorrows: borrow.Dec(),
	}, nil
}
