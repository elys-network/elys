package keeper

import (
	"context"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
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

func (k Keeper) GetAllPrices(goCtx context.Context, req *types.QueryGetAllPricesRequest) (*types.QueryGetAllPricesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	var prices []*types.Price

	assetentries := k.assetProfileKeeper.GetAllEntry(ctx)
	for _, assetEntry := range assetentries {
		if assetEntry.Denom == ptypes.Eden {
			assetEntry.Denom = ptypes.Elys
		}
		tokenPriceOracle := k.oracleKeeper.GetAssetPriceFromDenom(ctx, assetEntry.Denom).Mul(sdkmath.LegacyNewDec(int64(assetEntry.Decimals)))
		tokenPriceAmm := k.amm.CalcAmmPrice(ctx, assetEntry.Denom, assetEntry.Decimals).Mul(sdkmath.LegacyNewDec(int64(assetEntry.Decimals)))
		prices = append(prices, &types.Price{
			Denom:       assetEntry.Denom,
			OraclePrice: tokenPriceOracle,
			AmmPrice:    tokenPriceAmm,
		})
	}

	return &types.QueryGetAllPricesResponse{
		Prices: prices,
	}, nil
}
