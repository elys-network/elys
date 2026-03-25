package keeper

import (
	"context"
	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/types/query"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v6/x/amm/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) Params(goCtx context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	return &types.QueryParamsResponse{Params: k.GetParams(ctx)}, nil
}

func (k Keeper) WeightAndSlippageFee(goCtx context.Context, req *types.QueryWeightAndSlippageFeeRequest) (*types.QueryWeightAndSlippageFeeResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	return &types.QueryWeightAndSlippageFeeResponse{Value: k.GetWeightAndSlippageFee(ctx, req.PoolId, req.Date).Amount}, nil
}

func (k Keeper) MigrationReceipt(goCtx context.Context, req *types.QueryMigrationReceipt) (*types.QueryMigrationReceiptResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	addr, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &types.QueryMigrationReceiptResponse{Value: k.GetMigrationReceipt(ctx, addr)}, nil
}

func (k Keeper) MigrationReceiptPaginated(goCtx context.Context, req *types.QueryMigrationReceiptPaginated) (*types.QueryMigrationReceiptPaginatedResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var list []types.BalanceMigrationReceipt
	ctx := sdk.UnwrapSDKContext(goCtx)

	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	accountedPoolStore := prefix.NewStore(store, types.MigrationReceiptPrefix)

	pageRes, err := query.Paginate(accountedPoolStore, req.Pagination, func(key []byte, value []byte) error {
		var res types.BalanceMigrationReceipt
		if err := k.cdc.Unmarshal(value, &res); err != nil {
			return err
		}

		list = append(list, res)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryMigrationReceiptPaginatedResponse{List: list, Pagination: pageRes}, nil
}
