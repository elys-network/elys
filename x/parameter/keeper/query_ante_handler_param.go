package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/elys-network/elys/x/parameter/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) AnteHandlerParamAll(goCtx context.Context, req *types.QueryAllAnteHandlerParamRequest) (*types.QueryAllAnteHandlerParamResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var anteHandlerParams []types.AnteHandlerParam
	ctx := sdk.UnwrapSDKContext(goCtx)

	store := ctx.KVStore(k.storeKey)
	anteHandlerParamStore := prefix.NewStore(store, types.KeyPrefix(types.AnteHandlerParamKeyPrefix))

	pageRes, err := query.Paginate(anteHandlerParamStore, req.Pagination, func(key []byte, value []byte) error {
		var anteHandlerParam types.AnteHandlerParam
		if err := k.cdc.Unmarshal(value, &anteHandlerParam); err != nil {
			return err
		}

		anteHandlerParams = append(anteHandlerParams, anteHandlerParam)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllAnteHandlerParamResponse{AnteHandlerParam: anteHandlerParams, Pagination: pageRes}, nil
}
