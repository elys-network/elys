package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/elys-network/elys/x/incentive/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) ElysDelegatorAll(goCtx context.Context, req *types.QueryAllElysDelegatorRequest) (*types.QueryAllElysDelegatorResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var elysDelegators []types.ElysDelegator
	ctx := sdk.UnwrapSDKContext(goCtx)

	store := ctx.KVStore(k.storeKey)
	elysDelegatorStore := prefix.NewStore(store, types.KeyPrefix(types.ElysDelegatorKeyPrefix))

	pageRes, err := query.Paginate(elysDelegatorStore, req.Pagination, func(key []byte, value []byte) error {
		var elysDelegator types.ElysDelegator
		if err := k.cdc.Unmarshal(value, &elysDelegator); err != nil {
			return err
		}

		elysDelegators = append(elysDelegators, elysDelegator)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllElysDelegatorResponse{ElysDelegator: elysDelegators, Pagination: pageRes}, nil
}

func (k Keeper) ElysDelegator(goCtx context.Context, req *types.QueryGetElysDelegatorRequest) (*types.QueryGetElysDelegatorResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	val, found := k.GetElysDelegator(
		ctx,
		req.Index,
	)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetElysDelegatorResponse{ElysDelegator: val}, nil
}
