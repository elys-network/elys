package keeper

import (
	"context"

	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	ptypes "github.com/elys-network/elys/v6/x/parameter/types"
	"github.com/elys-network/elys/v6/x/perpetual/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) InvariantCustodyLiab(goCtx context.Context, req *types.QueryInvariantCustodyLiabRequest) (*types.QueryInvariantCustodyLiabResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	entry, found := k.assetProfileKeeper.GetEntry(ctx, ptypes.BaseCurrency)
	if !found {
		return nil, status.Error(codes.NotFound, "base currency not found")
	}
	baseCurrency := entry.Denom
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	mtpStore := prefix.NewStore(store, types.MTPPrefix)

	pagination := req.Pagination
	if pagination == nil {
		pagination = &query.PageRequest{
			Limit: 300000,
		}
	}

	totalCustodyWithoutLazy := sdk.Coins{}
	totalLiabilitiesWithoutLazy := sdk.Coins{}

	totalCustody := sdk.Coins{}
	totalLiabilities := sdk.Coins{}

	pageRes, err := query.Paginate(mtpStore, pagination, func(key []byte, value []byte) error {
		var mtp types.MTP
		k.cdc.MustUnmarshal(value, &mtp)

		if mtp.AmmPoolId != req.PoolId {
			return nil
		}

		totalCustodyWithoutLazy = totalCustodyWithoutLazy.Add(sdk.NewCoin(mtp.CustodyAsset, mtp.Custody))
		totalLiabilitiesWithoutLazy = totalLiabilitiesWithoutLazy.Add(sdk.NewCoin(mtp.LiabilitiesAsset, mtp.Liabilities))

		mtpAndPrice, err := k.fillMTPData(ctx, mtp, baseCurrency)
		if err != nil {
			return err
		}

		totalCustody = totalCustody.Add(sdk.NewCoin(mtpAndPrice.Mtp.CustodyAsset, mtpAndPrice.Mtp.Custody))
		totalLiabilities = totalLiabilities.Add(sdk.NewCoin(mtpAndPrice.Mtp.LiabilitiesAsset, mtpAndPrice.Mtp.Liabilities))

		return nil
	})
	if err != nil {
		return nil, err
	}

	pool, _ := k.GetPool(ctx, req.PoolId)

	return &types.QueryInvariantCustodyLiabResponse{
		TotalCustody:         totalCustody,
		TotalLiabilities:     totalLiabilities,
		TotalCustodyWithoutLazy: totalCustodyWithoutLazy,
		TotalLiabilitiesWithoutLazy: totalLiabilitiesWithoutLazy,
		PoolCustody: pool.,
		PoolLiabilities: pool.GetPoolAssetsShort(),
	}, nil
}
