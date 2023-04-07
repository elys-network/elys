package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/oracle/types"
)

// BandPriceResult returns the BandPrice result by RequestId
func (k Keeper) BandPriceResult(c context.Context, req *types.QueryBandPriceRequest) (*types.QueryBandPriceResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	result, err := k.GetBandPriceResult(ctx, types.OracleRequestID(req.RequestId))
	if err != nil {
		return nil, err
	}
	return &types.QueryBandPriceResponse{Result: &result}, nil
}

// LastBandRequestId returns the last BandPrice request Id
func (k Keeper) LastBandRequestId(c context.Context, req *types.QueryLastBandRequestIdRequest) (*types.QueryLastBandRequestIdResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	id := k.GetLastBandRequestId(ctx)
	return &types.QueryLastBandRequestIdResponse{RequestId: id}, nil
}
