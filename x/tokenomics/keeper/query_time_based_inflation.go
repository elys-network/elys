package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/runtime"

	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/elys-network/elys/v7/x/tokenomics/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) TimeBasedInflationAll(goCtx context.Context, req *types.QueryAllTimeBasedInflationRequest) (*types.QueryAllTimeBasedInflationResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var timeBasedInflations []types.TimeBasedInflation
	ctx := sdk.UnwrapSDKContext(goCtx)

	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	timeBasedInflationStore := prefix.NewStore(store, types.KeyPrefix(types.TimeBasedInflationKeyPrefix))

	pageRes, err := query.Paginate(timeBasedInflationStore, req.Pagination, func(key []byte, value []byte) error {
		var timeBasedInflation types.TimeBasedInflation
		if err := k.cdc.Unmarshal(value, &timeBasedInflation); err != nil {
			return err
		}

		timeBasedInflations = append(timeBasedInflations, timeBasedInflation)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllTimeBasedInflationResponse{TimeBasedInflation: timeBasedInflations, Pagination: pageRes}, nil
}

func (k Keeper) TimeBasedInflation(goCtx context.Context, req *types.QueryGetTimeBasedInflationRequest) (*types.QueryGetTimeBasedInflationResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	val, found := k.GetTimeBasedInflation(
		ctx,
		req.StartBlockHeight,
		req.EndBlockHeight,
	)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetTimeBasedInflationResponse{TimeBasedInflation: val}, nil
}
