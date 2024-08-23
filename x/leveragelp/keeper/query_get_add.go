package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/leveragelp/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) GetAdd(goCtx context.Context, req *types.QueryGetAddRequest) (*types.QueryGetAddResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	address := types.GetPositionAddress(uint64(req.Id))
	addressString, err := sdk.Bech32ifyAddressBytes("elys", address)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryGetAddResponse{
		Address: addressString,
	}, nil
}
