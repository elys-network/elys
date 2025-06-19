package keeper

import (
	"context"

	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/elys-network/elys/v6/x/clob/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) AllPerpetualADL(goCtx context.Context, req *types.AllPerpetualADLRequest) (*types.AllPerpetualADLResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	prefixStore := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.PerpetualADLPrefix)

	var adlList []types.PerpetualADL
	pageRes, err := query.Paginate(prefixStore, req.Pagination, func(key []byte, value []byte) error {
		var adl types.PerpetualADL
		if err := k.cdc.Unmarshal(value, &adl); err != nil {
			return err
		}

		adlList = append(adlList, adl)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &types.AllPerpetualADLResponse{AdlList: adlList, Pagination: pageRes}, nil
}
