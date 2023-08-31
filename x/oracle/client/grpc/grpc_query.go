package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/elys-network/elys/x/oracle/client"
	"github.com/elys-network/elys/x/oracle/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Querier struct {
	Q client.Querier
}

var _ types.QueryServer = Querier{}

func (q Querier) PriceAll(grpcCtx context.Context, req *types.QueryAllPriceRequest) (*types.QueryAllPriceResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	ctx := sdk.UnwrapSDKContext(grpcCtx)
	return q.Q.PriceAll(ctx, *req)
}

// implment those functions
func (q Querier) Params(context.Context, *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	panic("implement me")
}

// BandPriceResult defines a rpc handler method for MsgRequestBandPrice.
func (q Querier) BandPriceResult(context.Context, *types.QueryBandPriceRequest) (*types.QueryBandPriceResponse, error) {
	panic("implement me")
}

// LastBandRequestId query the last BandPrice result id
func (q Querier) LastBandRequestId(context.Context, *types.QueryLastBandRequestIdRequest) (*types.QueryLastBandRequestIdResponse, error) {
	panic("implement me")
}

// Queries a AssetInfo by denom.
func (q Querier) AssetInfo(context.Context, *types.QueryGetAssetInfoRequest) (*types.QueryGetAssetInfoResponse, error) {
	panic("implement me")
}

// Queries a list of AssetInfo items.
func (q Querier) AssetInfoAll(context.Context, *types.QueryAllAssetInfoRequest) (*types.QueryAllAssetInfoResponse, error) {
	panic("implement me")
}

// Queries a Price by asset.
func (q Querier) Price(context.Context, *types.QueryGetPriceRequest) (*types.QueryGetPriceResponse, error) {
	panic("implement me")
}

// Queries a PriceFeeder by feeder.
func (q Querier) PriceFeeder(context.Context, *types.QueryGetPriceFeederRequest) (*types.QueryGetPriceFeederResponse, error) {
	panic("implement me")
}

// Queries a list of PriceFeeder items.
func (q Querier) PriceFeederAll(context.Context, *types.QueryAllPriceFeederRequest) (*types.QueryAllPriceFeederResponse, error) {
	panic("implement me")
}
