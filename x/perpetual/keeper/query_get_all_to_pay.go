package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/perpetual/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) GetAllToPay(goCtx context.Context, req *types.QueryGetAllToPayRequest) (*types.QueryGetAllToPayResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	res := k.GetAllToPayStore(ctx)

	var toPay []*types.ToPay
	for _, item := range res {
		toPay = append(toPay, &item)
	}

	return &types.QueryGetAllToPayResponse{
		ToPay: toPay,
	}, nil
}
