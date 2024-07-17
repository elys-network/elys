package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/tier/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) GetConsolidatedPrice(goCtx context.Context, req *types.QueryGetConsolidatedPriceRequest) (*types.QueryGetConsolidatedPriceResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	oracle, amm, oracleDec := k.RetrieveConsolidatedPrice(ctx, req.Denom)

	return &types.QueryGetConsolidatedPriceResponse{
		AmmPrice:       amm,
		OraclePrice:    oracle,
		OraclePriceDec: oracleDec,
	}, nil
}
