package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/runtime"

	"cosmossdk.io/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/elys-network/elys/x/perpetual/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) Pools(goCtx context.Context, req *types.QueryAllPoolRequest) (*types.QueryAllPoolResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var pools []types.PoolResponse
	ctx := sdk.UnwrapSDKContext(goCtx)

	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	poolStore := prefix.NewStore(store, types.PoolKeyPrefix)

	pageRes, err := query.Paginate(poolStore, req.Pagination, func(key []byte, value []byte) error {
		var pool types.Pool
		if err := k.cdc.Unmarshal(value, &pool); err != nil {
			return err
		}

		ammPool, found := k.amm.GetPool(ctx, pool.AmmPoolId)
		if !found {
			return types.ErrPoolDoesNotExist
		}

		// TODO remove this if condition. All pools in perpetual must use oracle.
		if ammPool.PoolParams.UseOracle {
			totalLiabilities, err := k.GetPoolTotalBaseCurrencyLiabilities(ctx, pool)
			if err != nil {
				return err
			}

			longRate, shortRate := k.GetFundingPaymentRates(ctx, pool)
			pools = append(pools, types.PoolResponse{
				AmmPoolId:                            pool.AmmPoolId,
				Health:                               pool.Health,
				BorrowInterestRate:                   pool.BorrowInterestRate,
				PoolAssetsLong:                       pool.PoolAssetsLong,
				PoolAssetsShort:                      pool.PoolAssetsShort,
				LastHeightBorrowInterestRateComputed: pool.LastHeightBorrowInterestRateComputed,
				FundingRate:                          pool.FundingRate,
				NetOpenInterest:                      pool.GetNetOpenInterest(),
				LongRate:                             longRate,
				ShortRate:                            shortRate,
				TotalLiabilities:                     totalLiabilities,
				TotalLongOpenInterest:                pool.GetTotalLongOpenInterest(),
				TotalShortOpenInterest:               pool.GetTotalShortOpenInterest(),
				LeverageMax:                          pool.LeverageMax,
			})
		}

		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllPoolResponse{Pool: pools, Pagination: pageRes}, nil
}

func (k Keeper) Pool(goCtx context.Context, req *types.QueryGetPoolRequest) (*types.QueryGetPoolResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	val, found := k.GetPool(
		ctx,
		req.Index,
	)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	longRate, shortRate := k.GetFundingPaymentRates(ctx, val)

	totalLiabilities, err := k.GetPoolTotalBaseCurrencyLiabilities(ctx, val)
	if err != nil {
		return nil, err
	}

	pool := types.PoolResponse{
		AmmPoolId:                            val.AmmPoolId,
		Health:                               val.Health,
		BorrowInterestRate:                   val.BorrowInterestRate,
		PoolAssetsLong:                       val.PoolAssetsLong,
		PoolAssetsShort:                      val.PoolAssetsShort,
		LastHeightBorrowInterestRateComputed: val.LastHeightBorrowInterestRateComputed,
		FundingRate:                          val.FundingRate,
		NetOpenInterest:                      val.GetNetOpenInterest(),
		LongRate:                             longRate,
		ShortRate:                            shortRate,
		TotalLiabilities:                     totalLiabilities,
		TotalLongOpenInterest:                val.GetTotalLongOpenInterest(),
		TotalShortOpenInterest:               val.GetTotalShortOpenInterest(),
		LeverageMax:                          val.LeverageMax,
	}

	return &types.QueryGetPoolResponse{Pool: pool}, nil
}
