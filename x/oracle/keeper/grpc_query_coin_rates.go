package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/oracle/types"
)

// CoinRatesResult returns the CoinRates result by RequestId
func (k Keeper) CoinRatesResult(c context.Context, req *types.QueryCoinRatesRequest) (*types.QueryCoinRatesResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	result, err := k.GetCoinRatesResult(ctx, types.OracleRequestID(req.RequestId))
	if err != nil {
		return nil, err
	}
	return &types.QueryCoinRatesResponse{Result: &result}, nil
}

// LastBandRequestId returns the last CoinRates request Id
func (k Keeper) LastBandRequestId(c context.Context, req *types.QueryLastBandRequestIdRequest) (*types.QueryLastBandRequestIdResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	id := k.GetLastBandRequestId(ctx)
	return &types.QueryLastBandRequestIdResponse{RequestId: id}, nil
}
