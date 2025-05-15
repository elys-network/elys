package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/osmosis-labs/osmosis/osmomath"

	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/elys-network/elys/x/amm/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) PoolExtraInfo(ctx sdk.Context, pool types.Pool, days int) types.PoolExtraInfo {
	tvl, _ := pool.TVL(ctx, k.oracleKeeper, k.accountedPoolKeeper)
	lpTokenPrice, _ := pool.LpTokenPriceForShare(ctx, k.oracleKeeper, k.accountedPoolKeeper)
	avg := k.GetWeightBreakingSlippageAvg(ctx, pool.PoolId, days)
	apr := osmomath.ZeroBigDec()
	if tvl.IsPositive() {
		apr = avg.Mul(osmomath.NewBigDec(365)).Quo(tvl)
	}
	return types.PoolExtraInfo{
		Tvl:          tvl.Dec(),
		LpTokenPrice: lpTokenPrice.Dec(),
		LpSavedApr:   apr.Dec(),
	}
}

func (k Keeper) PoolAll(goCtx context.Context, req *types.QueryAllPoolRequest) (*types.QueryAllPoolResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var pools []types.Pool
	var extraInfos []types.PoolExtraInfo
	ctx := sdk.UnwrapSDKContext(goCtx)

	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	poolStore := prefix.NewStore(store, types.KeyPrefix(types.PoolKeyPrefix))

	pageRes, err := query.Paginate(poolStore, req.Pagination, func(key []byte, value []byte) error {
		var pool types.Pool
		if err := k.cdc.Unmarshal(value, &pool); err != nil {
			return err
		}

		pools = append(pools, pool)
		extraInfos = append(extraInfos, k.PoolExtraInfo(ctx, pool, int(req.Days)))
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllPoolResponse{Pool: pools, ExtraInfos: extraInfos, Pagination: pageRes}, nil
}

func (k Keeper) Pool(goCtx context.Context, req *types.QueryGetPoolRequest) (*types.QueryGetPoolResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	pool, found := k.GetPool(ctx, req.PoolId)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetPoolResponse{
		Pool:      pool,
		ExtraInfo: k.PoolExtraInfo(ctx, pool, int(req.Days)),
	}, nil
}
