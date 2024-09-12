package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/stablestake/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) QueryGetAllInterest(goCtx context.Context, req *types.QueryQueryGetAllInterestRequest) (*types.QueryQueryGetAllInterestResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	res := k.GetAllInterest(ctx)

	return &types.QueryQueryGetAllInterestResponse{
		InterestBlocks: convertToPointerSlice(res),
	}, nil
}

func convertToPointerSlice(blocks []types.InterestBlock) []*types.InterestBlock {
	pointerSlice := make([]*types.InterestBlock, len(blocks))
	for i, block := range blocks {
		pointerSlice[i] = &block
	}
	return pointerSlice
}
