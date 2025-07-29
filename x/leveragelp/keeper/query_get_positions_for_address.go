package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v7/x/leveragelp/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) QueryPositionsForAddress(goCtx context.Context, req *types.PositionsForAddressRequest) (*types.PositionsForAddressResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	addr, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, err
	}

	res, pageRes, err := k.GetPositionsForAddress(sdk.UnwrapSDKContext(goCtx), addr, req.Pagination)
	if err != nil {
		return nil, err
	}

	query_positions, err := k.GetLeverageLpUpdatedLeverage(sdk.UnwrapSDKContext(goCtx), res)
	if err != nil {
		return nil, err
	}

	positions_and_intrest, err := k.GetInterestRateUsd(sdk.UnwrapSDKContext(goCtx), query_positions)
	if err != nil {
		return nil, err
	}

	return &types.PositionsForAddressResponse{Positions: positions_and_intrest, Pagination: pageRes}, nil
}
