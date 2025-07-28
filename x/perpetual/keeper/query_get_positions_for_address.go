package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v7/x/perpetual/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) GetPositionsForAddress(goCtx context.Context, req *types.PositionsForAddressRequest) (*types.PositionsForAddressResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	addr, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, err
	}

	mtps, pageRes, err := k.GetMTPsForAddressWithPagination(sdk.UnwrapSDKContext(goCtx), addr, req.Pagination)
	if err != nil {
		return nil, err
	}

	return &types.PositionsForAddressResponse{Mtps: mtps, Pagination: pageRes}, nil
}
